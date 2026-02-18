package provider

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/pavel-snyk/snyk-sdk-go/v2/snyk"
)

var (
	_ resource.Resource              = (*brokerDeploymentCredentialResource)(nil)
	_ resource.ResourceWithConfigure = (*brokerDeploymentCredentialResource)(nil)
)

// brokerDeploymentCredentialResource defines the broker deployment credential resource implementation.
type brokerDeploymentCredentialResource struct {
	client *snyk.Client
}

// brokerDeploymentCredentialResourceModel describes the broker deployment credential resource data model.
type brokerDeploymentCredentialResourceModel struct {
	AppInstallID         types.String `tfsdk:"app_install_id"`
	BrokerConnectionType types.String `tfsdk:"broker_connection_type"`
	BrokerDeploymentID   types.String `tfsdk:"broker_deployment_id"`
	//Comment              types.String `tfsdk:"comment"` //todo: fix backend edge case
	EnvVarName types.String `tfsdk:"environment_variable_name"`
	ID         types.String `tfsdk:"id"`
	TenantID   types.String `tfsdk:"tenant_id"`
}

func NewBrokerDeploymentCredentialResource() resource.Resource {
	return &brokerDeploymentCredentialResource{}
}

func (r *brokerDeploymentCredentialResource) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = "snyk_broker_deployment_credential"
}

func (r *brokerDeploymentCredentialResource) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `
The broker deployment credential resource allows to manage Snyk broker deployment credentials.

A Snyk broker deployment credential is a local environment variable expected to be found
in a broker deployment. A broker deployment credential can be shared across broker connections
of the same type in the same broker deployment.
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
			"broker_connection_type": schema.StringAttribute{
				MarkdownDescription: "The ID of the associated broker deployment.",
				Required:            true,
			},
			"broker_deployment_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the associated broker deployment.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			//"comment": schema.StringAttribute{
			//	MarkdownDescription: "The comment of the broker deployment credential.",
			//	Optional:            true,
			//	Computed:            true,
			//	Default:             stringdefault.StaticString(""),
			//},
			"environment_variable_name": schema.StringAttribute{
				MarkdownDescription: "The name of the local environment variable expected to be found in broker deployment.",
				Required:            true,
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "The ID of the broker deployment credential.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"tenant_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the tenant to which the broker deployment credential belongs.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (r *brokerDeploymentCredentialResource) Configure(_ context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
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

func (r *brokerDeploymentCredentialResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var data brokerDeploymentCredentialResourceModel

	// read plan data into the model
	response.Diagnostics.Append(request.Plan.Get(ctx, &data)...)
	if response.Diagnostics.HasError() {
		return
	}

	tenantID := data.TenantID.ValueString()
	appInstallID := data.AppInstallID.ValueString()
	brokerDeploymentID := data.BrokerDeploymentID.ValueString()
	createRequest := &snyk.BrokerDeploymentCredentialCreateOrUpdateRequest{
		EnvVarName: data.EnvVarName.ValueString(),
		Type:       data.BrokerConnectionType.ValueString(),
	}
	tflog.Trace(ctx, "Creating broker deployment credential", map[string]any{
		"app_install_id":       appInstallID,
		"broker_deployment_id": brokerDeploymentID,
		"payload":              createRequest,
		"tenant_id":            tenantID,
	})
	brokerDeploymentCredential, resp, err := r.client.Brokers.CreateDeploymentCredential(ctx, tenantID, appInstallID, brokerDeploymentID, createRequest)
	if err != nil {
		response.Diagnostics.AddError("Unable to create broker deployment credential", err.Error())
		return
	}
	tflog.Trace(ctx, "Created broker deployment credential", map[string]any{
		"app_install_id":       appInstallID,
		"broker_deployment_id": brokerDeploymentID,
		"data":                 brokerDeploymentCredential,
		"snyk_request_id":      resp.SnykRequestID,
		"tenant_id":            tenantID,
	})

	// map response body to model
	data.AppInstallID = types.StringValue(appInstallID)
	data.BrokerConnectionType = types.StringValue(brokerDeploymentCredential.Attributes.Type)
	data.BrokerDeploymentID = types.StringValue(brokerDeploymentCredential.Attributes.BrokerDeploymentID)
	data.EnvVarName = types.StringValue(brokerDeploymentCredential.Attributes.EnvVarName)
	data.ID = types.StringValue(brokerDeploymentCredential.ID)
	data.TenantID = types.StringValue(tenantID)

	response.Diagnostics.Append(response.State.Set(ctx, &data)...)
}

func (r *brokerDeploymentCredentialResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var data brokerDeploymentCredentialResourceModel

	// read state data into the model
	response.Diagnostics.Append(request.State.Get(ctx, &data)...)
	if response.Diagnostics.HasError() {
		return
	}

	tenantID := data.TenantID.ValueString()
	appInstallID := data.AppInstallID.ValueString()
	brokerDeploymentID := data.BrokerDeploymentID.ValueString()
	brokerDeploymentCredentialID := data.ID.ValueString()

	tflog.Trace(ctx, "Getting broker deployment credential", map[string]any{
		"app_install_id":                  appInstallID,
		"broker_deployment_credential_id": brokerDeploymentCredentialID,
		"broker_deployment_id":            brokerDeploymentID,
		"tenant_id":                       tenantID,
	})
	brokerDeploymentCredential, resp, err := r.client.Brokers.GetDeploymentCredential(ctx, tenantID, appInstallID, brokerDeploymentID, brokerDeploymentCredentialID)
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			// if the broker deployment credential is somehow already destroyed, mark as successfully gone
			response.State.RemoveResource(ctx)
			return
		}
		response.Diagnostics.AddError("Unable to get broker deployment credential", err.Error())
		return
	}
	tflog.Trace(ctx, "Got broker deployment credential", map[string]any{
		"app_install_id":                  appInstallID,
		"broker_deployment_credential_id": brokerDeploymentCredentialID,
		"broker_deployment_id":            brokerDeploymentID,
		"data":                            brokerDeploymentCredential,
		"tenant_id":                       tenantID,
	})

	// map response body to model
	data.AppInstallID = types.StringValue(appInstallID)
	data.BrokerConnectionType = types.StringValue(brokerDeploymentCredential.Attributes.Type)
	data.BrokerDeploymentID = types.StringValue(brokerDeploymentCredential.Attributes.BrokerDeploymentID)
	data.EnvVarName = types.StringValue(brokerDeploymentCredential.Attributes.EnvVarName)
	data.ID = types.StringValue(brokerDeploymentCredential.ID)
	data.TenantID = types.StringValue(tenantID)

	response.Diagnostics.Append(response.State.Set(ctx, &data)...)
}

func (r *brokerDeploymentCredentialResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var data brokerDeploymentCredentialResourceModel

	// read plan data into the model
	response.Diagnostics.Append(request.Plan.Get(ctx, &data)...)
	if response.Diagnostics.HasError() {
		return
	}

	tenantID := data.TenantID.ValueString()
	appInstallID := data.AppInstallID.ValueString()
	brokerDeploymentID := data.BrokerDeploymentID.ValueString()
	brokerDeploymentCredentialID := data.ID.ValueString()

	updateRequest := &snyk.BrokerDeploymentCredentialCreateOrUpdateRequest{
		EnvVarName: data.EnvVarName.ValueString(),
		Type:       data.BrokerConnectionType.ValueString(),
	}
	tflog.Trace(ctx, "Updating broker deployment credential", map[string]any{
		"app_install_id":                  appInstallID,
		"broker_deployment_credential_id": brokerDeploymentCredentialID,
		"broker_deployment_id":            brokerDeploymentID,
		"payload":                         updateRequest,
		"tenant_id":                       tenantID,
	})
	brokerDeploymentCredential, resp, err := r.client.Brokers.UpdateDeploymentCredential(ctx, tenantID, appInstallID, brokerDeploymentID, brokerDeploymentCredentialID, updateRequest)
	if err != nil {
		response.Diagnostics.AddError("Unable to update broker deployment credential", err.Error())
		return
	}
	tflog.Trace(ctx, "Updated broker deployment credential", map[string]any{
		"app_install_id":                  appInstallID,
		"broker_deployment_credential_id": brokerDeploymentCredentialID,
		"broker_deployment_id":            brokerDeploymentID,
		"data":                            brokerDeploymentCredential,
		"snyk_request_id":                 resp.SnykRequestID,
		"tenant_id":                       tenantID,
	})

	// map response body to model
	data.AppInstallID = types.StringValue(appInstallID)
	data.BrokerConnectionType = types.StringValue(brokerDeploymentCredential.Attributes.Type)
	data.BrokerDeploymentID = types.StringValue(brokerDeploymentCredential.Attributes.BrokerDeploymentID)
	data.EnvVarName = types.StringValue(brokerDeploymentCredential.Attributes.EnvVarName)
	data.ID = types.StringValue(brokerDeploymentCredential.ID)
	data.TenantID = types.StringValue(tenantID)

	response.Diagnostics.Append(response.State.Set(ctx, &data)...)
}

func (r *brokerDeploymentCredentialResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var data brokerDeploymentCredentialResourceModel

	// read state data into the model
	response.Diagnostics.Append(request.State.Get(ctx, &data)...)
	if response.Diagnostics.HasError() {
		return
	}

	tenantID := data.TenantID.ValueString()
	appInstallID := data.AppInstallID.ValueString()
	brokerDeploymentID := data.BrokerDeploymentID.ValueString()
	brokerDeploymentCredentialID := data.ID.ValueString()
	tflog.Trace(ctx, "Deleting broker deployment credential", map[string]any{
		"app_install_id":                  appInstallID,
		"broker_deployment_credential_id": brokerDeploymentCredentialID,
		"broker_deployment_id":            brokerDeploymentID,
		"tenant_id":                       tenantID,
	})
	resp, err := r.client.Brokers.DeleteDeploymentCredential(ctx, tenantID, appInstallID, brokerDeploymentID, brokerDeploymentCredentialID)
	if err != nil {
		response.Diagnostics.AddError("Unable to delete broker deployment credential", err.Error())
		return
	}
	tflog.Trace(ctx, "Deleted broker deployment credential", map[string]any{
		"app_install_id":                  appInstallID,
		"broker_deployment_credential_id": brokerDeploymentCredentialID,
		"broker_deployment_id":            brokerDeploymentID,
		"snyk_request_id":                 resp.SnykRequestID,
		"tenant_id":                       tenantID,
	})
}
