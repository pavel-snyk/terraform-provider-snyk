package provider

import (
	"context"

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
	_ resource.Resource                = (*appInstallResource)(nil)
	_ resource.ResourceWithConfigure   = (*appInstallResource)(nil)
	_ resource.ResourceWithImportState = (*appInstallResource)(nil)
)

// appInstallResource defines the app installation resource implementation.
type appInstallResource struct {
	client *snyk.Client
}

// appInstallResourceModel describes the app installation resource data model.
type appInstallResourceModel struct {
	AppID        types.String `tfsdk:"app_id"`
	AppName      types.String `tfsdk:"app_name"`
	ClientID     types.String `tfsdk:"client_id"`
	ClientSecret types.String `tfsdk:"client_secret"`
	ID           types.String `tfsdk:"id"`
	OrgID        types.String `tfsdk:"organization_id"`
}

func NewAppInstallResource() resource.Resource {
	return &appInstallResource{}
}

func (r *appInstallResource) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = "snyk_app_install"
}

func (r *appInstallResource) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `
The app install resource allows to manage Snyk app installations.

Snyk Apps are the modern and preferred way to build integrations with Snyk,
exposing fine-grained scopes for accessing resources over the Snyk APIs,
powered by OAuth 2.0 for a developer-friendly experience. See [Snyk Apps](https://docs.snyk.io/snyk-api/using-specific-snyk-apis/snyk-apps-apis).
`,
		Attributes: map[string]schema.Attribute{
			"app_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the app.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"app_name": schema.StringAttribute{
				MarkdownDescription: "The name of the app.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"client_id": schema.StringAttribute{
				MarkdownDescription: "The OAuth2 client id for the app installation.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"client_secret": schema.StringAttribute{
				MarkdownDescription: "The OAuth2 client secret for the app.",
				Computed:            true,
				Sensitive:           true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "The ID of the app installation.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"organization_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the organization where the app will be installed.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (r *appInstallResource) Configure(_ context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
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

func (r *appInstallResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var data appInstallResourceModel

	// read plan data into the model
	response.Diagnostics.Append(request.Plan.Get(ctx, &data)...)
	if response.Diagnostics.HasError() {
		return
	}

	appID := data.AppID.ValueString()
	orgID := data.OrgID.ValueString()
	tflog.Trace(ctx, "Creating app install for organization", map[string]any{"app_id": appID, "org_id": orgID})
	appInstall, resp, err := r.client.Apps.CreateAppInstallForOrg(ctx, orgID, appID)
	if err != nil {
		response.Diagnostics.AddError("Unable to create app install", err.Error())
		return
	}
	tflog.Trace(ctx, "Created app install for organization", map[string]any{
		"app_id":          appID,
		"data":            appInstall,
		"org_id":          orgID,
		"snyk_request_id": resp.SnykRequestID,
	})

	// client secret is shown only once by creation
	data.ClientSecret = types.StringValue(appInstall.Attributes.ClientSecret)

	tflog.Trace(ctx, "Getting app install with enriched properties", map[string]any{"app_install_id": appInstall.ID, "org_id": orgID})
	// list all app installed for org and filter our new created app install
	appInstalls, resp, err := r.client.Apps.ListAppInstallsForOrg(ctx, orgID, nil)
	if err != nil {
		response.Diagnostics.AddError("Unable to get app install", err.Error())
		return
	}
	for _, ai := range appInstalls {
		if ai.ID == appInstall.ID {
			appInstall = &ai
			break
		}
	}
	tflog.Trace(ctx, "Got app install with enriched properties", map[string]any{
		"data":            appInstall,
		"org_id":          orgID,
		"snyk_request_id": resp.SnykRequestID,
	})

	// map response body to model
	if appInstall.Relationships != nil {
		data.AppID = types.StringValue(appInstall.Relationships.App.Data.ID)
		data.AppName = types.StringValue(appInstall.Relationships.App.Data.Attributes.Name)
	}
	data.ClientID = types.StringValue(appInstall.Attributes.ClientID)
	data.ID = types.StringValue(appInstall.ID)
	data.OrgID = types.StringValue(orgID)

	response.Diagnostics.Append(response.State.Set(ctx, &data)...)
}

func (r *appInstallResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var data appInstallResourceModel

	// read state data into the model
	response.Diagnostics.Append(request.State.Get(ctx, &data)...)
	if response.Diagnostics.HasError() {
		return
	}

	appInstallID := data.ID.ValueString()
	orgID := data.OrgID.ValueString()
	tflog.Trace(ctx, "Getting app install", map[string]any{"app_install_id": appInstallID, "org_id": orgID})
	appInstalls, resp, err := r.client.Apps.ListAppInstallsForOrg(ctx, orgID, nil)
	if err != nil {
		response.Diagnostics.AddError("Unable to get app install", err.Error())
		return
	}
	var appInstall *snyk.AppInstall
	for _, ai := range appInstalls {
		if ai.ID == appInstallID {
			appInstall = &ai
			break
		}
	}
	// if the app install is somehow already destroyed, mark as successfully gone
	if appInstall == nil {
		response.State.RemoveResource(ctx)
		return
	}
	tflog.Trace(ctx, "Got app install", map[string]any{
		"data":            appInstall,
		"org_id":          orgID,
		"snyk_request_id": resp.SnykRequestID,
	})

	// map response body to model
	if appInstall.Relationships != nil {
		data.AppID = types.StringValue(appInstall.Relationships.App.Data.ID)
		data.AppName = types.StringValue(appInstall.Relationships.App.Data.Attributes.Name)
	}
	data.ClientID = types.StringValue(appInstall.Attributes.ClientID)
	data.ID = types.StringValue(appInstall.ID)
	data.OrgID = types.StringValue(orgID)

	response.Diagnostics.Append(response.State.Set(ctx, &data)...)
}

func (r *appInstallResource) Update(context.Context, resource.UpdateRequest, *resource.UpdateResponse) {
	// app_id and organization_id are mandatory and RequiresReplace()
	// no update available
}

func (r *appInstallResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var data appInstallResourceModel

	// read state data into the model
	response.Diagnostics.Append(request.State.Get(ctx, &data)...)
	if response.Diagnostics.HasError() {
		return
	}

	appInstallID := data.ID.ValueString()
	orgID := data.OrgID.ValueString()
	tflog.Trace(ctx, "Deleting app install for org", map[string]any{"app_install_id": appInstallID, "org_id": orgID})
	resp, err := r.client.Apps.DeleteAppInstallFromOrg(ctx, orgID, appInstallID)
	if err != nil {
		response.Diagnostics.AddError("Unable to delete app install", err.Error())
		return
	}
	tflog.Trace(ctx, "Deleted app install for org", map[string]any{
		"app_install_id":  appInstallID,
		"org_id":          orgID,
		"snyk_request_id": resp.SnykRequestID,
	})
}

func (r *appInstallResource) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), request, response)
}
