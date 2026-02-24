package provider

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/pavel-snyk/snyk-sdk-go/v2/snyk"

	"github.com/pavel-snyk/terraform-provider-snyk/internal/provider/helper"
)

var (
	_ resource.Resource              = (*brokerConnectionResource)(nil)
	_ resource.ResourceWithConfigure = (*brokerConnectionResource)(nil)
)

// brokerConnectionResource defines the broker connection resource implementation.
type brokerConnectionResource struct {
	client *snyk.Client
}

// brokerConnectionResourceModel describes the broker connection resource data model.
type brokerConnectionResourceModel struct {
	AppInstallID       types.String `tfsdk:"app_install_id"`
	BrokerDeploymentID types.String `tfsdk:"broker_deployment_id"`
	Configuration      types.Object `tfsdk:"configuration"` // brokerConnectionConfigurationModel
	ID                 types.String `tfsdk:"id"`
	Name               types.String `tfsdk:"name"`
	TenantID           types.String `tfsdk:"tenant_id"`
	Type               types.String `tfsdk:"type"`
}

type brokerConnectionResourceConfigurationModel struct {
	BrokerClientURL          types.String `tfsdk:"broker_client_url"`
	GitLabHostname           types.String `tfsdk:"gitlab_hostname"`
	GitLabTokenCredentialID  types.String `tfsdk:"gitlab_token_credential_id"`
	JiraHostname             types.String `tfsdk:"jira_hostname"`
	JiraPasswordCredentialID types.String `tfsdk:"jira_password_credential_id"`
	JiraPATCredentialID      types.String `tfsdk:"jira_pat_credential_id"`
	JiraUsername             types.String `tfsdk:"jira_username"`
}

func NewBrokerConnectionResource() resource.Resource {
	return &brokerConnectionResource{}
}

func (r *brokerConnectionResource) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = "snyk_broker_connection"
}

func (r *brokerConnectionResource) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `
The broker connection resource allows to manage Snyk broker connections.

A Snyk broker connection lives in Snyk broker deployment and is configured to communicate
with specific private resources: SCMs, JIRA, and others. For more information,
see [Universal Broker documentation](https://docs.snyk.io/implementation-and-setup/enterprise-setup/snyk-broker/universal-broker).
`,
		Attributes: map[string]schema.Attribute{
			"app_install_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the app installation for Universal Broker Snyk App.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"broker_deployment_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the associated broker deployment.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"configuration": schema.SingleNestedAttribute{
				MarkdownDescription: "The configuration parameters depending on broker connection type.",
				Required:            true,
				Validators: []validator.Object{
					helper.RequiresValidConfiguration(path.MatchRoot("type")),
				},
				Attributes: map[string]schema.Attribute{
					"broker_client_url": schema.StringAttribute{
						MarkdownDescription: "The URL of the broker client used by webhooks. It's recommended to use regional Snyk API URL.",
						Optional:            true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"gitlab_hostname": schema.StringAttribute{
						MarkdownDescription: "The GitLab hostname.",
						Optional:            true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"gitlab_token_credential_id": schema.StringAttribute{
						MarkdownDescription: "The ID of the broker deployment credential for GitLab token.",
						Optional:            true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"jira_hostname": schema.StringAttribute{
						MarkdownDescription: "The Jira hostname.",
						Optional:            true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"jira_password_credential_id": schema.StringAttribute{
						MarkdownDescription: "The ID of the broker deployment credential for Jira password.",
						Optional:            true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"jira_pat_credential_id": schema.StringAttribute{
						MarkdownDescription: "The ID of the broker deployment credential for Jira PAT token.",
						Optional:            true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"jira_username": schema.StringAttribute{
						MarkdownDescription: "The Jira username.",
						Optional:            true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
				},
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "The ID of the broker connection.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the broker connection.",
				Required:            true,
			},
			"tenant_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the tenant to which the broker connection belongs.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "The type of the broker connection.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				//TODO: remove after RequiresValidConfiguration completion
				Validators: []validator.String{
					stringvalidator.OneOf("gitlab", "jira"),
				},
			},
		},
	}
}

func (r *brokerConnectionResource) Configure(_ context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
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

func (r *brokerConnectionResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var data brokerConnectionResourceModel

	// read plan data into the model
	response.Diagnostics.Append(request.Plan.Get(ctx, &data)...)
	if response.Diagnostics.HasError() {
		return
	}
	var dataConfiguration brokerConnectionResourceConfigurationModel
	response.Diagnostics.Append(request.Plan.GetAttribute(ctx, path.Root("configuration"), &dataConfiguration)...)
	if response.Diagnostics.HasError() {
		return
	}

	tenantID := data.TenantID.ValueString()
	appInstallID := data.AppInstallID.ValueString()
	brokerDeploymentID := data.BrokerDeploymentID.ValueString()

	createRequest := &snyk.BrokerConnectionCreateOrUpdateRequest{
		BrokerClientURL: dataConfiguration.BrokerClientURL.ValueString(),
		GitLabHostname:  dataConfiguration.GitLabHostname.ValueString(),
		GitLabToken:     dataConfiguration.GitLabTokenCredentialID.ValueString(),
		Name:            data.Name.ValueString(),
		Type:            snyk.BrokerConnectionType(data.Type.ValueString()),
	}
	tflog.Trace(ctx, "Creating broker connection", map[string]any{
		"app_install_id":       appInstallID,
		"broker_deployment_id": brokerDeploymentID,
		"payload":              createRequest,
		"tenant_id":            tenantID,
	})
	brokerConnection, resp, err := r.client.Brokers.CreateConnection(ctx, tenantID, appInstallID, brokerDeploymentID, createRequest)
	if err != nil {
		response.Diagnostics.AddError("Unable to create broker connection", err.Error())
		return
	}
	tflog.Trace(ctx, "Created broker connection", map[string]any{
		"app_install_id":       appInstallID,
		"broker_deployment_id": brokerDeploymentID,
		"data":                 brokerConnection,
		"snyk_request_id":      resp.SnykRequestID,
		"tenant_id":            tenantID,
	})

	// map response body to model
	data.AppInstallID = types.StringValue(appInstallID)
	data.BrokerDeploymentID = types.StringValue(brokerConnection.Attributes.BrokerDeploymentID)
	// configuration todo
	data.ID = types.StringValue(brokerConnection.ID)
	data.Name = types.StringValue(brokerConnection.Attributes.Name)
	data.TenantID = types.StringValue(tenantID)
	data.Type = types.StringValue(string(brokerConnection.Attributes.Configuration.Type))

	response.Diagnostics.Append(response.State.Set(ctx, &data)...)
}

func (r *brokerConnectionResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var data brokerConnectionResourceModel

	// read state data into the model
	response.Diagnostics.Append(request.State.Get(ctx, &data)...)
	if response.Diagnostics.HasError() {
		return
	}
	var dataConfiguration brokerConnectionResourceConfigurationModel
	response.Diagnostics.Append(request.State.GetAttribute(ctx, path.Root("configuration"), &dataConfiguration)...)
	if response.Diagnostics.HasError() {
		return
	}

	tenantID := data.TenantID.ValueString()
	appInstallID := data.AppInstallID.ValueString()
	brokerDeploymentID := data.BrokerDeploymentID.ValueString()
	brokerConnectionID := data.ID.ValueString()

	tflog.Trace(ctx, "Getting broker connection", map[string]any{
		"app_install_id":       appInstallID,
		"broker_deployment_id": brokerDeploymentID,
		"broker_connection_id": brokerConnectionID,
		"tenant_id":            tenantID,
	})
	brokerConnection, resp, err := r.client.Brokers.GetConnection(ctx, tenantID, appInstallID, brokerDeploymentID, brokerConnectionID)
	if err != nil {
		// if connection is somehow already destroyed, mark as successfully gone
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			response.State.RemoveResource(ctx)
			return
		}
		response.Diagnostics.AddError("Unable to get broker connection", err.Error())
		return
	}
	tflog.Trace(ctx, "Got broker connection", map[string]any{
		"app_install_id":       appInstallID,
		"broker_deployment_id": brokerDeploymentID,
		"data":                 brokerConnection,
		"snyk_request_id":      resp.SnykRequestID,
		"tenant_id":            tenantID,
	})

	// map response body to model
	data.AppInstallID = types.StringValue(appInstallID)
	data.BrokerDeploymentID = types.StringValue(brokerConnection.Attributes.BrokerDeploymentID)
	// configuration todo
	data.ID = types.StringValue(brokerConnection.ID)
	data.Name = types.StringValue(brokerConnection.Attributes.Name)
	data.TenantID = types.StringValue(tenantID)
	data.Type = types.StringValue(string(brokerConnection.Attributes.Configuration.Type))

	response.Diagnostics.Append(response.State.Set(ctx, &data)...)
}

func (r *brokerConnectionResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var data brokerConnectionResourceModel

	// read plan data into the model
	response.Diagnostics.Append(request.Plan.Get(ctx, &data)...)
	if response.Diagnostics.HasError() {
		return
	}
	var dataConfiguration brokerConnectionResourceConfigurationModel
	response.Diagnostics.Append(request.Plan.GetAttribute(ctx, path.Root("configuration"), &dataConfiguration)...)
	if response.Diagnostics.HasError() {
		return
	}

	tenantID := data.TenantID.ValueString()
	appInstallID := data.AppInstallID.ValueString()
	brokerDeploymentID := data.BrokerDeploymentID.ValueString()
	brokerConnectionID := data.ID.ValueString()

	updateRequest := &snyk.BrokerConnectionCreateOrUpdateRequest{
		BrokerClientURL: dataConfiguration.BrokerClientURL.ValueString(),
		GitLabHostname:  dataConfiguration.GitLabHostname.ValueString(),
		GitLabToken:     dataConfiguration.GitLabTokenCredentialID.ValueString(),
		Name:            data.Name.ValueString(),
		Type:            snyk.BrokerConnectionType(data.Type.ValueString()),
	}
	tflog.Trace(ctx, "Updating broker connection", map[string]any{
		"app_install_id":       appInstallID,
		"broker_deployment_id": brokerDeploymentID,
		"payload":              updateRequest,
		"tenant_id":            tenantID,
	})
	brokerConnection, resp, err := r.client.Brokers.UpdateConnection(ctx, tenantID, appInstallID, brokerDeploymentID, brokerConnectionID, updateRequest)
	if err != nil {
		response.Diagnostics.AddError("Unable to create broker connection", err.Error())
		return
	}
	tflog.Trace(ctx, "Updated broker connection", map[string]any{
		"app_install_id":       appInstallID,
		"broker_deployment_id": brokerDeploymentID,
		"data":                 brokerConnection,
		"snyk_request_id":      resp.SnykRequestID,
		"tenant_id":            tenantID,
	})

	// map response body to model
	data.AppInstallID = types.StringValue(appInstallID)
	data.BrokerDeploymentID = types.StringValue(brokerConnection.Attributes.BrokerDeploymentID)
	// configuration todo
	data.ID = types.StringValue(brokerConnection.ID)
	data.Name = types.StringValue(brokerConnection.Attributes.Name)
	data.TenantID = types.StringValue(tenantID)
	data.Type = types.StringValue(string(brokerConnection.Attributes.Configuration.Type))

	response.Diagnostics.Append(response.State.Set(ctx, &data)...)
}

func (r *brokerConnectionResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var data brokerConnectionResourceModel

	// read state data into the model
	response.Diagnostics.Append(request.State.Get(ctx, &data)...)
	if response.Diagnostics.HasError() {
		return
	}

	tenantID := data.TenantID.ValueString()
	appInstallID := data.AppInstallID.ValueString()
	brokerDeploymentID := data.BrokerDeploymentID.ValueString()
	brokerConnectionID := data.ID.ValueString()

	tflog.Trace(ctx, "Deleting broker connection", map[string]any{
		"app_install_id":       appInstallID,
		"broker_deployment_id": brokerDeploymentID,
		"broker_connection_id": brokerConnectionID,
		"tenant_id":            tenantID,
	})
	resp, err := r.client.Brokers.DeleteConnection(ctx, tenantID, appInstallID, brokerDeploymentID, brokerConnectionID)
	if err != nil {
		response.Diagnostics.AddError("Unable to delete broker connection", err.Error())
		return
	}
	tflog.Trace(ctx, "Deleted broker connection", map[string]any{
		"app_install_id":       appInstallID,
		"broker_deployment_id": brokerDeploymentID,
		"broker_connection_id": brokerConnectionID,
		"snyk_request_id":      resp.SnykRequestID,
		"tenant_id":            tenantID,
	})
}
