package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/pavel-snyk/snyk-sdk-go/v2/snyk"
)

var (
	_ datasource.DataSource              = (*appInstallDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*appInstallDataSource)(nil)
)

// appInstallDataSource defines the app installation datasource implementation.
type appInstallDataSource struct {
	client *snyk.Client
}

// appInstallDataSourceModel describes the datasource data model.
type appInstallDataSourceModel struct {
	AppID          types.String `tfsdk:"app_id"`
	AppName        types.String `tfsdk:"app_name"`
	ClientID       types.String `tfsdk:"client_id"`
	ID             types.String `tfsdk:"id"`
	OrganizationID types.String `tfsdk:"organization_id"`
}

func NewAppInstallDataSource() datasource.DataSource {
	return &appInstallDataSource{}
}

func (d *appInstallDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = "snyk_app_install"
}

func (d *appInstallDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `
The app install data source provides information about an existing app installation.

Snyk Apps are the modern and preferred way to build integrations with Snyk,
exposing fine-grained scopes for accessing resources over the Snyk APIs,
powered by OAuth 2.0 for a developer-friendly experience. See (Snyk Apps)[https://docs.snyk.io/snyk-api/using-specific-snyk-apis/snyk-apps-apis].
`,
		Attributes: map[string]schema.Attribute{
			"app_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the app.",
				Computed:            true,
				Optional:            true,
			},
			"app_name": schema.StringAttribute{
				MarkdownDescription: "The name of the app.",
				Computed:            true,
				Optional:            true,
			},
			"client_id": schema.StringAttribute{
				MarkdownDescription: "The OAuth2 client id for the app installation.",
				Computed:            true,
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "The ID of the app installation.",
				Computed:            true,
				Optional:            true,
			},
			"organization_id": schema.StringAttribute{
				MarkdownDescription: "The organization ID of the app installation.",
				Required:            true,
			},
		},
	}
}

func (d *appInstallDataSource) Configure(_ context.Context, request datasource.ConfigureRequest, response *datasource.ConfigureResponse) {
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

	d.client = client
}

func (d *appInstallDataSource) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var data appInstallDataSourceModel

	response.Diagnostics.Append(request.Config.Get(ctx, &data)...)
	if response.Diagnostics.HasError() {
		return
	}

	appInstallID := data.ID.ValueString()
	appID := data.AppID.ValueString()
	appName := data.AppName.ValueString()
	if appInstallID == "" && appID == "" && appName == "" {
		response.Diagnostics.AddError(
			"Missing required attributes",
			`The attribute "id", "app_id" or "app_name" must be defined.`,
		)
		return
	}
	orgID := data.OrganizationID.ValueString()

	tflog.Trace(ctx, "Getting app installs for organization", map[string]any{"organization_id": orgID})
	appInstalls, resp, err := d.client.Apps.ListAppInstallsForOrg(ctx, orgID, nil)
	if err != nil {
		response.Diagnostics.AddError("Unable to get app installs", err.Error())
		return
	}
	tflog.Trace(ctx, "Got app installs for organization", map[string]any{"data": appInstalls, "snyk_request_id": resp.SnykRequestID})

	tflog.Info(ctx, "Searching for app install by criteria", map[string]any{
		"app_install_id": appInstallID,
		"app_id":         appID,
		"app_name":       appName,
	})
	var foundAppInstalls []snyk.AppInstall
	var partialMatchReason string
	for _, ai := range appInstalls {
		fullMatch, reason := checkAppInstallMatch(&ai, appInstallID, appID, appName)
		if fullMatch {
			foundAppInstalls = append(foundAppInstalls, ai)
		} else if partialMatchReason == "" && reason != "" {
			// store the first partial match reason to provide a helpful error diagnostic
			partialMatchReason = reason
		}
	}

	d.handleSearchResults(foundAppInstalls, partialMatchReason, response)
	if response.Diagnostics.HasError() {
		return
	}

	appInstall := foundAppInstalls[0]

	// map response body to attributes
	data.AppID = types.StringValue(appInstall.Relationships.App.Data.ID)
	data.AppName = types.StringValue(appInstall.Relationships.App.Data.Attributes.Name)
	data.ClientID = types.StringValue(appInstall.Attributes.ClientID)
	data.ID = types.StringValue(appInstall.ID)
	data.OrganizationID = types.StringValue(orgID)

	response.Diagnostics.Append(response.State.Set(ctx, &data)...)
}

// checkAppInstallMatch determines if a given app install matches the user-provided criteria.
// It returns true for a full match. For a non-match, it may return a reason for a partial mismatch.
func checkAppInstallMatch(ai *snyk.AppInstall, appInstallID, appID, appName string) (bool, string) {
	if appInstallID != "" && appInstallID != ai.ID {
		return false, ""
	}
	if appID != "" && appID != ai.Relationships.App.Data.ID {
		// if appInstallID matched, this is a partial match
		if appInstallID != "" {
			return false, fmt.Sprintf(
				"An app install with id '%s' was found, but it has a different app_id.\n"+
					"Expected '%s', got '%s'.",
				appInstallID, appID, ai.Relationships.App.Data.ID,
			)
		}
		return false, ""
	}
	if appName != "" && appName != ai.Relationships.App.Data.Attributes.Name {
		// if appInstallID or appID matched, this is a partial match
		if appInstallID != "" {
			return false, fmt.Sprintf(
				"An app install with id '%s' was found, but it has a different app_name.\n"+
					"Expected '%s', got '%s'.",
				appInstallID, appName, ai.Relationships.App.Data.Attributes.Name,
			)
		}
		if appID != "" {
			return false, fmt.Sprintf(
				"An app install with app id '%s' was found, but it has a different app_name.\n"+
					"Expected '%s', got '%s'.",
				appID, appName, ai.Relationships.App.Data.Attributes.Name,
			)
		}
		return false, ""
	}

	// all provided criteria have matched
	return true, ""
}

func (d *appInstallDataSource) handleSearchResults(fullMatches []snyk.AppInstall, partialMatchReason string, response *datasource.ReadResponse) {
	if len(fullMatches) == 0 {
		if partialMatchReason != "" {
			response.Diagnostics.AddError(
				"No search results due to mismatch",
				partialMatchReason,
			)
			return
		}
		response.Diagnostics.AddError(
			"No search results",
			"No app install matched the provided criteria. Please verify the 'organization_id' and other search attributes.",
		)
		return
	}

	if len(fullMatches) > 1 {
		var foundAppInstallIDs []string
		for _, fm := range fullMatches {
			foundAppInstallIDs = append(foundAppInstallIDs, fm.ID)
		}
		response.Diagnostics.AddError(
			"Ambiguous search results",
			fmt.Sprintf("The provided criteria match multiple app installs.\n"+
				"Please provide a more specific combination of search attributes such as 'id', 'app_id' or 'app_name', to uniquely identify one.\n"+
				"Found app install ids: %v", foundAppInstallIDs),
		)
		return
	}
}
