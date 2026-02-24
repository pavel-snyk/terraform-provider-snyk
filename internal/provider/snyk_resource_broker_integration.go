package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/pavel-snyk/snyk-sdk-go/v2/snyk"
)

var (
	_ resource.Resource              = (*brokerIntegrationResource)(nil)
	_ resource.ResourceWithConfigure = (*brokerIntegrationResource)(nil)
)

// brokerIntegrationResource defines the broker integration resource implementation.
type brokerIntegrationResource struct {
	client *snyk.Client
}

// brokerIntegrationResourceModel describes the broker integration resource data model.
type brokerIntegrationResourceModel struct {
	BrokerConnectionID types.String `tfsdk:"broker_connection_id"`
	ID                 types.String `tfsdk:"id"`
	OrgID              types.String `tfsdk:"organization_id"`
	TenantID           types.String `tfsdk:"tenant_id"`
	Type               types.String `tfsdk:"type"`
}

func NewBrokerIntegrationResource() resource.Resource {
	return &brokerIntegrationResource{}
}

func (r *brokerIntegrationResource) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = "snyk_broker_integration"
}

func (r *brokerIntegrationResource) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `
The broker integration resource allows to manage Snyk broker integrations.

A Snyk broker integration creates a link between a broker connection and a Snyk organization.
You can integrate broker connection with many organizations as needed. After integration the
broker connection will be available when you run broker client.
For more information, see [Universal Broker documentation](https://docs.snyk.io/implementation-and-setup/enterprise-setup/snyk-broker/universal-broker).
`,
		Attributes: map[string]schema.Attribute{
			"broker_connection_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the associated broker connection.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "The ID of the broker integration.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"organization_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the organization to the broker integration belongs.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"tenant_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the tenant to which the broker integration belongs.",
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
			},
		},
	}
}

func (r *brokerIntegrationResource) Configure(_ context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
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

func (r *brokerIntegrationResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var data brokerIntegrationResourceModel

	// read plan data into the model
	response.Diagnostics.Append(request.Plan.Get(ctx, &data)...)
	if response.Diagnostics.HasError() {
		return
	}

	tenantID := data.TenantID.ValueString()
	brokerConnectionID := data.BrokerConnectionID.ValueString()
	orgID := data.OrgID.ValueString()

	createRequest := &snyk.BrokerIntegrationCreateRequest{
		Type: snyk.BrokerConnectionType(data.Type.ValueString()),
	}
	tflog.Trace(ctx, "Creating broker integration", map[string]any{
		"broker_connection_id": brokerConnectionID,
		"organization_id":      orgID,
		"payload":              createRequest,
		"tenant_id":            tenantID,
	})
	brokerIntegration, resp, err := r.client.Brokers.CreateIntegration(ctx, tenantID, brokerConnectionID, orgID, createRequest)
	if err != nil {
		response.Diagnostics.AddError("Unable to create broker integration", err.Error())
		return
	}
	tflog.Trace(ctx, "Created broker integration", map[string]any{
		"broker_connection_id": brokerConnectionID,
		"data":                 brokerIntegration,
		"organization_id":      orgID,
		"snyk_request_id":      resp.SnykRequestID,
		"tenant_id":            tenantID,
	})

	// map response body to model
	data.BrokerConnectionID = types.StringValue(brokerConnectionID)
	data.ID = types.StringValue(brokerIntegration.ID)
	data.OrgID = types.StringValue(brokerIntegration.OrgID)
	data.TenantID = types.StringValue(tenantID)
	data.Type = types.StringValue(data.Type.ValueString())

	response.Diagnostics.Append(response.State.Set(ctx, &data)...)
}

func (r *brokerIntegrationResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var data brokerIntegrationResourceModel

	// read state data into the model
	response.Diagnostics.Append(request.State.Get(ctx, &data)...)
	if response.Diagnostics.HasError() {
		return
	}

	tenantID := data.TenantID.ValueString()
	brokerConnectionID := data.BrokerConnectionID.ValueString()
	brokerIntegrationID := data.ID.ValueString()

	tflog.Trace(ctx, "Getting broker integration", map[string]any{
		"broker_connection_id":  brokerConnectionID,
		"broker_integration_id": brokerIntegrationID,
		"tenant_id":             tenantID,
	})
	brokerIntegrations, resp, err := r.client.Brokers.ListIntegrations(ctx, tenantID, brokerConnectionID)
	if err != nil {
		response.Diagnostics.AddError("Unable to get broker integration", err.Error())
		return
	}
	var brokerIntegration *snyk.BrokerIntegration
	for _, i := range brokerIntegrations {
		if i.ID == brokerIntegrationID {
			brokerIntegration = &i
			break
		}
	}
	// if the broker integration is somehow already destroyed, mark as successfully gone
	if brokerIntegration == nil {
		response.State.RemoveResource(ctx)
		return
	}
	tflog.Trace(ctx, "Got broker integration", map[string]any{
		"broker_connection_id":  brokerConnectionID,
		"broker_integration_id": brokerIntegrationID,
		"snyk_request_id":       resp.SnykRequestID,
		"tenant_id":             tenantID,
	})

	// map response body to model
	data.BrokerConnectionID = types.StringValue(brokerConnectionID)
	data.ID = types.StringValue(brokerIntegration.ID)
	data.OrgID = types.StringValue(brokerIntegration.OrgID)
	data.TenantID = types.StringValue(tenantID)
	data.Type = types.StringValue(brokerIntegration.IntegrationType)

	response.Diagnostics.Append(response.State.Set(ctx, &data)...)
}

func (r *brokerIntegrationResource) Update(_ context.Context, _ resource.UpdateRequest, _ *resource.UpdateResponse) {
	// all attributes are mandatory and RequiresReplace()
	// no update available
}

func (r *brokerIntegrationResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var data brokerIntegrationResourceModel

	// read state data into the model
	response.Diagnostics.Append(request.State.Get(ctx, &data)...)
	if response.Diagnostics.HasError() {
		return
	}

	tenantID := data.TenantID.ValueString()
	brokerConnectionID := data.BrokerConnectionID.ValueString()
	orgID := data.OrgID.ValueString()
	brokerIntegrationID := data.ID.ValueString()

	tflog.Trace(ctx, "Deleting broker integration", map[string]any{
		"broker_connection_id":  brokerConnectionID,
		"broker_integration_id": brokerIntegrationID,
		"organization_id":       orgID,
		"tenant_id":             tenantID,
	})
	resp, err := r.client.Brokers.DeleteIntegration(ctx, tenantID, brokerConnectionID, orgID, brokerIntegrationID)
	if err != nil {
		response.Diagnostics.AddError("Unable to delete broker integration", err.Error())
		return
	}
	tflog.Trace(ctx, "Deleted broker integration", map[string]any{
		"broker_connection_id":  brokerConnectionID,
		"broker_integration_id": brokerIntegrationID,
		"organization_id":       orgID,
		"snyk_request_id":       resp.SnykRequestID,
		"tenant_id":             tenantID,
	})
}
