package provider

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/pavel-snyk/snyk-sdk-go/v2/snyk"
)

var (
	_ resource.Resource                = (*organizationResource)(nil)
	_ resource.ResourceWithConfigure   = (*organizationResource)(nil)
	_ resource.ResourceWithImportState = (*organizationResource)(nil)
)

// organizationResource defines the organization resource implementation.
type organizationResource struct {
	client *snyk.Client
}

// organizationResourceModel describes the organization resource data model.
type organizationResourceModel struct {
	GroupID  types.String `tfsdk:"group_id"`
	ID       types.String `tfsdk:"id"`
	Name     types.String `tfsdk:"name"`
	Slug     types.String `tfsdk:"slug"`
	TenantID types.String `tfsdk:"tenant_id"`
}

func NewOrganizationResource() resource.Resource {
	return &organizationResource{}
}

func (r *organizationResource) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = "snyk_organization"
}

func (r *organizationResource) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `
The organization resource allows to manage Snyk organizations.

An Organization in Snyk is a way to collect and organize your Projects. Members of Organizations
have access to these Projects. See [Manage Groups and Organizations](https://docs.snyk.io/snyk-platform-administration/groups-and-organizations).
`,
		Attributes: map[string]schema.Attribute{
			"group_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the group to which the organization belongs.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "The ID of the organization.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the organization.",
				Required:            true,
			},
			"slug": schema.StringAttribute{
				MarkdownDescription: "The canonical (unique and URL-friendly) name of the organization.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"tenant_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the tenant to which the organization belongs.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *organizationResource) Configure(_ context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
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

func (r *organizationResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var data organizationResourceModel

	// read plan data into the model
	response.Diagnostics.Append(request.Plan.Get(ctx, &data)...)
	if response.Diagnostics.HasError() {
		return
	}

	createRequest := &snyk.OrganizationV1CreateRequest{
		Name: data.Name.ValueString(),
	}
	if data.GroupID.ValueString() != "" {
		createRequest.GroupID = data.GroupID.ValueString()
	}
	tflog.Trace(ctx, "Creating organization", map[string]any{"payload": createRequest})
	orgV1, resp, err := r.client.OrgsV1.Create(ctx, createRequest)
	if err != nil {
		response.Diagnostics.AddError("Unable to create organization", err.Error())
		return
	}
	tflog.Trace(ctx, "Created organization", map[string]any{
		"data":            orgV1,
		"snyk_request_id": resp.SnykRequestID,
	})

	tflog.Trace(ctx, "Getting organization with enriched properties", map[string]any{"org_id": orgV1.ID})
	organization, resp, err := r.client.Orgs.Get(ctx, orgV1.ID, &snyk.GetOrganizationOptions{Expand: "tenant"})
	if err != nil {
		response.Diagnostics.AddError("Unable to get organization", err.Error())
		return
	}
	tflog.Trace(ctx, "Got organization with enriched properties", map[string]any{
		"data":            organization,
		"org_id":          orgV1.ID,
		"snyk_request_id": resp.SnykRequestID,
	})
	// tenant info is not available immediately after org creation with V1, we fetch tenantID from group endpoint
	tflog.Trace(ctx, "Getting group with tenant information", map[string]any{"group_id": orgV1.Group.ID})
	group, resp, err := r.client.Groups.Get(ctx, orgV1.Group.ID)
	if err != nil {
		response.Diagnostics.AddError("Unable to get group", err.Error())
		return
	}
	tflog.Trace(ctx, "Got group with tenant information", map[string]any{
		"data":            group,
		"group_id":        orgV1.Group.ID,
		"snyk_request_id": resp.SnykRequestID,
	})

	// map response body to model
	data.GroupID = types.StringValue(organization.Attributes.GroupID)
	data.ID = types.StringValue(organization.ID)
	data.Name = types.StringValue(organization.Attributes.Name)
	data.Slug = types.StringValue(organization.Attributes.Slug)
	if group.Relationships != nil {
		data.TenantID = types.StringValue(group.Relationships.Tenant.Data.ID)
	}

	response.Diagnostics.Append(response.State.Set(ctx, &data)...)
}

func (r *organizationResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var data organizationResourceModel

	// read state data into the model
	response.Diagnostics.Append(request.State.Get(ctx, &data)...)
	if response.Diagnostics.HasError() {
		return
	}

	orgID := data.ID.ValueString()
	tflog.Trace(ctx, "Getting organization", map[string]any{"org_id": orgID})
	organization, resp, err := r.client.Orgs.Get(ctx, orgID, nil)
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			// if the organization is somehow already destroyed, mark as successfully gone
			response.State.RemoveResource(ctx)
			return
		}
		response.Diagnostics.AddError("Unable to get organization", err.Error())
		return
	}
	tflog.Trace(ctx, "Got organization", map[string]any{
		"data":            organization,
		"org_id":          orgID,
		"snyk_request_id": resp.SnykRequestID,
	})

	// map response body to model
	data.GroupID = types.StringValue(organization.Attributes.GroupID)
	data.ID = types.StringValue(organization.ID)
	data.Name = types.StringValue(organization.Attributes.Name)
	data.Slug = types.StringValue(organization.Attributes.Slug)
	if organization.Relationships != nil {
		data.TenantID = types.StringValue(organization.Relationships.Tenant.Data.ID)
	}

	response.Diagnostics.Append(response.State.Set(ctx, &data)...)
}

func (r *organizationResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var data organizationResourceModel

	// read plan data into the model
	response.Diagnostics.Append(request.Plan.Get(ctx, &data)...)
	if response.Diagnostics.HasError() {
		return
	}

	orgID := data.ID.ValueString()
	updateRequest := &snyk.OrganizationUpdateRequest{
		Name: data.Name.ValueString(),
	}
	tflog.Trace(ctx, "Updating organization", map[string]any{"payload": updateRequest})
	organization, resp, err := r.client.Orgs.Update(ctx, orgID, updateRequest)
	if err != nil {
		response.Diagnostics.AddError("Unable to update organization", err.Error())
		return
	}
	tflog.Trace(ctx, "Updated organization", map[string]any{
		"data":            organization,
		"snyk_request_id": resp.SnykRequestID,
	})

	// map response body to model
	data.GroupID = types.StringValue(organization.Attributes.GroupID)
	data.ID = types.StringValue(organization.ID)
	data.Name = types.StringValue(organization.Attributes.Name)
	data.Slug = types.StringValue(organization.Attributes.Slug)
	if organization.Relationships != nil {
		data.TenantID = types.StringValue(organization.Relationships.Tenant.Data.ID)
	}

	response.Diagnostics.Append(response.State.Set(ctx, &data)...)
}

func (r *organizationResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var data organizationResourceModel

	// read state data into the model
	response.Diagnostics.Append(request.State.Get(ctx, &data)...)
	if response.Diagnostics.HasError() {
		return
	}

	orgID := data.ID.ValueString()
	tflog.Trace(ctx, "Deleting organization", map[string]any{"org_id": orgID})
	resp, err := r.client.OrgsV1.Delete(ctx, orgID)
	if err != nil {
		response.Diagnostics.AddError("Unable to delete organization", err.Error())
		return
	}
	tflog.Trace(ctx, "Deleted organization", map[string]any{
		"org_id":          orgID,
		"snyk_request_id": resp.SnykRequestID,
	})
}

func (r *organizationResource) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), request, response)
}
