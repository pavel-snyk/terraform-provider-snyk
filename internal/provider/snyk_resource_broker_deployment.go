package provider

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/pavel-snyk/snyk-sdk-go/v2/snyk"
)

var (
	_ resource.Resource                = (*brokerDeploymentResource)(nil)
	_ resource.ResourceWithConfigure   = (*brokerDeploymentResource)(nil)
	_ resource.ResourceWithImportState = (*brokerDeploymentResource)(nil)
)

// brokerDeploymentResource defines the broker deployment resource implementation.
type brokerDeploymentResource struct {
	client *snyk.Client
}

// brokerDeploymentResourceModel describes the broker deployment resource data model.
type brokerDeploymentResourceModel struct {
	AppInstallID types.String `tfsdk:"app_install_id"`
	ID           types.String `tfsdk:"id"`
	Metadata     types.Map    `tfsdk:"metadata"`
	OrgID        types.String `tfsdk:"organization_id"`
	TenantID     types.String `tfsdk:"tenant_id"`
}

func NewBrokerDeploymentResource() resource.Resource {
	return &brokerDeploymentResource{}
}

func (r *brokerDeploymentResource) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = "snyk_broker_deployment"
}

func (r *brokerDeploymentResource) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `
The broker deployment resource allows to manage Snyk broker deployments.

A Snyk broker deployment is the recommended way to manage Snyk Universal Broker. It allows you
to group broker connections into separate deployments for better organization and management.
For more information, see [Universal Broker documentation](https://docs.snyk.io/implementation-and-setup/enterprise-setup/snyk-broker/universal-broker).
`,
		Attributes: map[string]schema.Attribute{
			"app_install_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the app installation for Universal Broker Snyk App.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "The ID of the broker deployment.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"metadata": schema.MapAttribute{
				MarkdownDescription: "A map of string to string to store custom metadata for the broker deployment. " +
					"This can be useful for tracking ownership, environment, or other identifying information.",
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
				Default:     mapdefault.StaticValue(types.MapValueMust(types.StringType, map[string]attr.Value{})),
			},
			"organization_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the organization where the Universal Broker Snyk App is installed. " +
					"If omitted, the provider will search for the app installation across all accessible organizations. " +
					"It's recommended to set for faster performance.",
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"tenant_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the tenant to which the broker deployment belongs.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (r *brokerDeploymentResource) Configure(_ context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}

	client, ok := request.ProviderData.(*snyk.Client)
	if !ok {
		response.Diagnostics.AddError(
			"Unconfigured Snyk client",
			"Please report this issue to the provider developers.",
		)
		return
	}

	r.client = client
}

func (r *brokerDeploymentResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var data brokerDeploymentResourceModel

	// read plan data into the model
	response.Diagnostics.Append(request.Plan.Get(ctx, &data)...)
	if response.Diagnostics.HasError() {
		return
	}

	appInstallID := data.AppInstallID.ValueString()
	orgID := data.OrgID.ValueString()
	tenantID := data.TenantID.ValueString()
	var metadata map[string]string
	if !data.Metadata.IsNull() && !data.Metadata.IsUnknown() {
		response.Diagnostics.Append(data.Metadata.ElementsAs(ctx, &metadata, false)...)
		if response.Diagnostics.HasError() {
			return
		}
	}

	if orgID == "" {
		tflog.Info(ctx, "Searching in all accessible orgs for installed app", map[string]any{"app_install_id": appInstallID})

		tflog.Trace(ctx, "Getting all accessible organizations")
		orgs, errf := r.client.Orgs.AllAccessibleOrgs(ctx, nil)
		for org := range orgs {
			tflog.Debug(ctx, "Searching all app installs for org", map[string]any{"org_id": org.ID, "org_name": org.Attributes.Name})

			tflog.Trace(ctx, "Getting app installs for org", map[string]any{"org_id": org.ID})
			appInstalls, resp, err := r.client.Apps.ListAppInstallsForOrg(ctx, org.ID, nil)
			if err != nil {
				response.Diagnostics.AddError("Unable to get app installs", err.Error())
				return
			}
			tflog.Trace(ctx, "Got app installs for org", map[string]any{
				"data":            appInstalls,
				"org_id":          org.ID,
				"snyk_request_id": resp.SnykRequestID,
			})

			for _, ai := range appInstalls {
				if ai.ID == appInstallID {
					tflog.Info(ctx, "Found org for configured app install", map[string]any{
						"org_id":   org.ID,
						"org_name": org.Attributes.Name,
					})
					orgID = org.ID
					break
				}
			}
		}
		if err := errf(); err != nil {
			response.Diagnostics.AddError("Unable to get organizations", err.Error())
			return
		}

		if orgID == "" {
			response.Diagnostics.AddError(
				"Unable to find organization",
				fmt.Sprintf("No organization with app installation (%v) was found.", appInstallID),
			)
			return
		}
	}

	createRequest := &snyk.BrokerDeploymentCreateOrUpdateRequest{OrgID: orgID, Metadata: metadata}
	tflog.Trace(ctx, "Creating broker deployment", map[string]any{
		"app_install_id": appInstallID,
		"payload":        createRequest,
		"tenant_id":      tenantID,
	})
	brokerDeployment, resp, err := r.client.Brokers.CreateDeployment(ctx, tenantID, appInstallID, createRequest)
	if err != nil {
		response.Diagnostics.AddError("Unable to create broker deployment", err.Error())
		return
	}
	tflog.Trace(ctx, "Created broker deployment", map[string]any{
		"app_install_id":  appInstallID,
		"data":            brokerDeployment,
		"snyk_request_id": resp.SnykRequestID,
		"tenant_id":       tenantID,
	})

	// map response body to model
	data.AppInstallID = types.StringValue(brokerDeployment.Attributes.AppInstallID)
	data.ID = types.StringValue(brokerDeployment.ID)
	data.OrgID = types.StringValue(orgID)
	data.TenantID = types.StringValue(tenantID)
	if len(brokerDeployment.Attributes.Metadata) > 0 {
		metadataMap, diags := types.MapValueFrom(ctx, types.StringType, brokerDeployment.Attributes.Metadata)
		response.Diagnostics.Append(diags...)
		data.Metadata = metadataMap
	} else {
		data.Metadata = types.MapValueMust(types.StringType, map[string]attr.Value{})
	}

	response.Diagnostics.Append(response.State.Set(ctx, &data)...)
}

func (r *brokerDeploymentResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var data brokerDeploymentResourceModel

	// read state data into the model
	response.Diagnostics.Append(request.State.Get(ctx, &data)...)
	if response.Diagnostics.HasError() {
		return
	}

	appInstallID := data.AppInstallID.ValueString()
	brokerDeploymentID := data.ID.ValueString()
	orgID := data.OrgID.ValueString()
	tenantID := data.TenantID.ValueString()

	tflog.Trace(ctx, "Getting broker deployment", map[string]any{
		"app_install_id":       appInstallID,
		"broker_deployment_id": brokerDeploymentID,
		"tenant_id":            tenantID,
	})
	brokerDeployments, resp, err := r.client.Brokers.ListDeployments(ctx, tenantID, appInstallID)
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			// api returns 404 instead of the empty array for no deployments
			// this case is considered as deployment is somehow already destroyed, mark as successfully gone
			response.State.RemoveResource(ctx)
			return
		}
		response.Diagnostics.AddError("Unable to get broker deployment", err.Error())
		return
	}
	var brokerDeployment *snyk.BrokerDeployment
	for _, bd := range brokerDeployments {
		if bd.ID == brokerDeploymentID {
			brokerDeployment = &bd
			break
		}
	}
	if brokerDeployment == nil {
		response.State.RemoveResource(ctx)
		return
	}
	tflog.Trace(ctx, "Got broker deployment", map[string]any{
		"app_install_id":  appInstallID,
		"data":            brokerDeployment,
		"snyk_request_id": resp.SnykRequestID,
		"tenant_id":       tenantID,
	})

	if orgID == "" {
		tflog.Info(ctx, "Searching in all accessible organizations for app install", map[string]any{"app_install_id": appInstallID})

		tflog.Trace(ctx, "Getting all accessible organizations")
		orgs, errf := r.client.Orgs.AllAccessibleOrgs(ctx, nil)
		for org := range orgs {
			tflog.Debug(ctx, "Searching all app installs for org", map[string]any{
				"organization_id": org.ID, "organization_name": org.Attributes.Name,
			})

			tflog.Trace(ctx, "Getting app installs for organization", map[string]any{
				"organization_id": org.ID, "organization_name": org.Attributes.Name,
			})
			appInstalls, resp, err := r.client.Apps.ListAppInstallsForOrg(ctx, org.ID, nil)
			if err != nil {
				response.Diagnostics.AddError("Unable to get app installs", err.Error())
				return
			}
			tflog.Trace(ctx, "Got app installs for organization", map[string]any{
				"data":            appInstalls,
				"organization_id": org.ID,
				"snyk_request_id": resp.SnykRequestID,
			})

			for _, ai := range appInstalls {
				if ai.ID == appInstallID {
					tflog.Info(ctx, "Found org for configured app install", map[string]any{
						"organization_id": org.ID, "organization_name": org.Attributes.Name,
					})
					orgID = org.ID
					break
				}
			}
		}
		if err := errf(); err != nil {
			response.Diagnostics.AddError("Unable to get organizations", err.Error())
			return
		}

		if orgID == "" {
			response.Diagnostics.AddError(
				"Unable to find organization",
				fmt.Sprintf("No organization with app installation (%v) was found.", appInstallID),
			)
			return
		}
	}

	// map response body to model
	data.AppInstallID = types.StringValue(brokerDeployment.Attributes.AppInstallID)
	data.ID = types.StringValue(brokerDeployment.ID)
	data.OrgID = types.StringValue(orgID)
	data.TenantID = types.StringValue(tenantID)
	if len(brokerDeployment.Attributes.Metadata) > 0 {
		metadataMap, diags := types.MapValueFrom(ctx, types.StringType, brokerDeployment.Attributes.Metadata)
		response.Diagnostics.Append(diags...)
		data.Metadata = metadataMap
	} else {
		data.Metadata = types.MapValueMust(types.StringType, map[string]attr.Value{})
	}

	response.Diagnostics.Append(response.State.Set(ctx, &data)...)
}

func (r *brokerDeploymentResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var data brokerDeploymentResourceModel

	// read plan data into the model
	response.Diagnostics.Append(request.Plan.Get(ctx, &data)...)
	if response.Diagnostics.HasError() {
		return
	}

	appInstallID := data.AppInstallID.ValueString()
	brokerDeploymentID := data.ID.ValueString()
	orgID := data.OrgID.ValueString()
	tenantID := data.TenantID.ValueString()
	var metadata map[string]string
	if !data.Metadata.IsNull() && !data.Metadata.IsUnknown() {
		response.Diagnostics.Append(data.Metadata.ElementsAs(ctx, &metadata, false)...)
		if response.Diagnostics.HasError() {
			return
		}
	}

	updateRequest := &snyk.BrokerDeploymentCreateOrUpdateRequest{OrgID: orgID, Metadata: metadata}
	tflog.Trace(ctx, "Updating broker deployment", map[string]any{
		"app_install_id":       appInstallID,
		"broker_deployment_id": brokerDeploymentID,
		"payload":              updateRequest,
		"tenant_id":            tenantID,
	})
	brokerDeployment, resp, err := r.client.Brokers.UpdateDeployment(ctx, tenantID, appInstallID, brokerDeploymentID, updateRequest)
	if err != nil {
		response.Diagnostics.AddError("Unable to update broker deployment", err.Error())
		return
	}
	tflog.Trace(ctx, "Updated broker deployment", map[string]any{
		"app_install_id":  appInstallID,
		"data":            brokerDeployment,
		"snyk_request_id": resp.SnykRequestID,
		"tenant_id":       tenantID,
	})

	// map response body to model
	data.AppInstallID = types.StringValue(brokerDeployment.Attributes.AppInstallID)
	data.ID = types.StringValue(brokerDeployment.ID)
	data.OrgID = types.StringValue(orgID)
	data.TenantID = types.StringValue(tenantID)
	if len(brokerDeployment.Attributes.Metadata) > 0 {
		metadataMap, diags := types.MapValueFrom(ctx, types.StringType, brokerDeployment.Attributes.Metadata)
		response.Diagnostics.Append(diags...)
		data.Metadata = metadataMap
	} else {
		data.Metadata = types.MapValueMust(types.StringType, map[string]attr.Value{})
	}

	response.Diagnostics.Append(response.State.Set(ctx, &data)...)
}

func (r *brokerDeploymentResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var data brokerDeploymentResourceModel

	// read state data into the model
	response.Diagnostics.Append(request.State.Get(ctx, &data)...)
	if response.Diagnostics.HasError() {
		return
	}

	appInstallID := data.AppInstallID.ValueString()
	tenantID := data.TenantID.ValueString()
	brokerDeploymentID := data.ID.ValueString()
	tflog.Trace(ctx, "Deleting broker deployment", map[string]any{
		"app_install_id":       appInstallID,
		"broker_deployment_id": brokerDeploymentID,
		"tenant_id":            tenantID,
	})
	resp, err := r.client.Brokers.DeleteDeployment(ctx, tenantID, appInstallID, brokerDeploymentID)
	if err != nil {
		response.Diagnostics.AddError("Unable to delete broker deployment", err.Error())
		return
	}
	tflog.Trace(ctx, "Deleted broker deployment", map[string]any{
		"app_install_id":       appInstallID,
		"broker_deployment_id": brokerDeploymentID,
		"snyk_request_id":      resp.SnykRequestID,
		"tenant_id":            tenantID,
	})
}

func (r *brokerDeploymentResource) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	idParts := strings.Split(request.ID, ",")

	if len(idParts) != 3 || idParts[0] == "" || idParts[1] == "" || idParts[2] == "" {
		response.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: id,tenant_id,app_install_id. Got: %q", request.ID),
		)
		return
	}

	response.Diagnostics.Append(response.State.SetAttribute(ctx, path.Root("id"), idParts[0])...)
	response.Diagnostics.Append(response.State.SetAttribute(ctx, path.Root("tenant_id"), idParts[1])...)
	response.Diagnostics.Append(response.State.SetAttribute(ctx, path.Root("app_install_id"), idParts[2])...)
}
