package provider

import (
	"context"
	"fmt"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/pavel-snyk/snyk-sdk-go/v2/snyk"
)

const (
	defaultSnykRegion = "SNYK-US-01"
)

var (
	_ provider.Provider                   = &snykProvider{}
	_ provider.ProviderWithValidateConfig = &snykProvider{}
)

// snykProvider defines the provider implementation.
type snykProvider struct {
	// version is set to
	//  - the provider version on release
	//  - "dev" when the provider is built and ran locally
	//  - "testacc" when running acceptance tests
	version string
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &snykProvider{
			version: version,
		}
	}
}

func (p *snykProvider) Metadata(_ context.Context, _ provider.MetadataRequest, response *provider.MetadataResponse) {
	response.TypeName = "snyk"
	response.Version = p.version
}

func (p *snykProvider) Schema(_ context.Context, _ provider.SchemaRequest, response *provider.SchemaResponse) {
	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"region": schema.SingleNestedAttribute{
				MarkdownDescription: fmt.Sprintf("Configuration for the Snyk Region. If not provided, the provider will "+
					"use the `SNYK_REGION` environment variable, or default to  **%s**.\n"+
					"    - to use a **predefined Snyk region** (e.g., `SNYK-EU-01`, `SNYK-AU-01`), provide only the `name` attribute. "+
					"See the official Snyk documentation for a list of [available region names](https://docs.snyk.io/snyk-data-and-governance/regional-hosting-and-data-residency#available-snyk-regions).\n"+
					"    - to use a **custom or private Snyk region**, provide all attributes: `name`, `app_base_url`, `rest_base_url` and `v1_base_url`.", defaultSnykRegion),
				Attributes: map[string]schema.Attribute{
					"name": schema.StringAttribute{
						MarkdownDescription: "The name of Snyk region. For predefined regions, this is the short name (e.g. `SNYK-EU-01`). For customer regions, this is a user-defined identifier.",
						Optional:            true,
					},
					"app_base_url": schema.StringAttribute{
						MarkdownDescription: "The application base URL for a custom region. Must be provided along with `name`, `rest_base_url` and `v1_base_url` when defining a custom region.",
						Optional:            true,
					},
					"rest_base_url": schema.StringAttribute{
						MarkdownDescription: "The REST API base URL for a custom region. Must be provided along with `name`, `app_base_url` and `v1_base_url` when defining a custom region.",
						Optional:            true,
					},
					"v1_base_url": schema.StringAttribute{
						MarkdownDescription: "The V1 API base URL for a custom region. Must be provided along with `name`, `app_base_url` and `rest_base_url` when defining a custom region.",
						Optional:            true,
					},
				},
				Optional: true,
			},
			"token": schema.StringAttribute{
				MarkdownDescription: "This Snyk API token. It can also be sourced from the `SNYK_TOKEN` environment variable.",
				Optional:            true,
				Sensitive:           true,
			},
		},
	}
}

type snykProviderModel struct {
	Region types.Object `tfsdk:"region"` // snykProviderRegionModel
	Token  types.String `tfsdk:"token"`
}

type snykProviderRegionModel struct {
	AppBaseURL  types.String `tfsdk:"app_base_url"`
	Name        types.String `tfsdk:"name"`
	RESTBaseURL types.String `tfsdk:"rest_base_url"`
	V1BaseURL   types.String `tfsdk:"v1_base_url"`
}

func (p *snykProvider) ValidateConfig(ctx context.Context, request provider.ValidateConfigRequest, response *provider.ValidateConfigResponse) {
	var config snykProviderModel

	response.Diagnostics.Append(request.Config.Get(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}

	// validate region block
	if !config.Region.IsNull() && !config.Region.IsUnknown() {
		var regionConfig snykProviderRegionModel
		diags := request.Config.GetAttribute(ctx, path.Root("region"), &regionConfig)
		response.Diagnostics.Append(diags...)
		if response.Diagnostics.HasError() {
			return
		}

		hasAppBaseURL := !regionConfig.AppBaseURL.IsNull()
		hasRESTBaseURL := !regionConfig.RESTBaseURL.IsNull()
		hasV1BaseURL := !regionConfig.V1BaseURL.IsNull()

		// for custom region all parts must be present
		if hasAppBaseURL || hasRESTBaseURL || hasV1BaseURL {
			if regionConfig.Name.IsNull() || !hasAppBaseURL || !hasRESTBaseURL || !hasV1BaseURL {
				response.Diagnostics.AddAttributeError(
					path.Root("region"),
					"Invalid provider config",
					`For a custom region, the "name", "app_base_url", "rest_base_url" and "v1_base_url" attributes must all be set.`,
				)
			}
		}
	}

	// validate token
	if config.Token.IsNull() && os.Getenv("SNYK_TOKEN") == "" {
		response.Diagnostics.AddAttributeError(
			path.Root("token"),
			"Invalid provider config",
			`The Snyk API token must be provided via the "token" attribute or the "SNYK_TOKEN" environment variable.`,
		)
	}
}

func (p *snykProvider) Configure(ctx context.Context, request provider.ConfigureRequest, response *provider.ConfigureResponse) {
	var config snykProviderModel

	response.Diagnostics.Append(request.Config.Get(ctx, &config)...)
	if response.Diagnostics.HasError() {
		return
	}

	opts := []snyk.ClientOption{snyk.WithUserAgent(p.userAgent())}

	// region logic
	if config.Region.IsNull() || config.Region.IsUnknown() {
		// no region block, fallback to env var, then default
		regionName := os.Getenv("SNYK_REGION")
		if regionName == "" {
			regionName = defaultSnykRegion
		}
		opts = append(opts, snyk.WithRegionAlias(regionName))
	} else {
		var regionConfig snykProviderRegionModel
		diags := request.Config.GetAttribute(ctx, path.Root("region"), &regionConfig)
		response.Diagnostics.Append(diags...)
		if response.Diagnostics.HasError() {
			return
		}

		if !regionConfig.AppBaseURL.IsNull() {
			// all attributes are present because ValidateConfig passed
			opts = append(opts, snyk.WithRegion(snyk.Region{
				Alias:       regionConfig.Name.ValueString(),
				AppBaseURL:  regionConfig.AppBaseURL.ValueString(),
				RESTBaseURL: regionConfig.RESTBaseURL.ValueString(),
				V1BaseURL:   regionConfig.V1BaseURL.ValueString(),
			}))
		} else {
			if !regionConfig.Name.IsNull() {
				// predefined region
				opts = append(opts, snyk.WithRegionAlias(regionConfig.Name.ValueString()))
			} else {
				// fallback to default by empty region block
				opts = append(opts, snyk.WithRegionAlias(defaultSnykRegion))
			}
		}
	}

	// token logic
	token := config.Token.ValueString()
	if token == "" {
		token = os.Getenv("SNYK_TOKEN")
	}

	client, err := snyk.NewClient(token, opts...)
	if err != nil {
		response.Diagnostics.AddError(
			"Unable to create Snyk API client",
			"An unexpected error occurred when creating the Snyk API client. "+
				"Please contact the plugin developers.\n\n"+
				"Snyk Client error:"+err.Error(),
		)
		return
	}

	response.DataSourceData = client
	response.ResourceData = client
}

func (p *snykProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewOrganizationDataSource,
		//NewProjectDataSource,
		//NewUserDataSource,
	}
}

func (p *snykProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		//NewIntegrationResource,
		//NewOrganizationResource,
	}
}

func (p *snykProvider) userAgent() string {
	name := "terraform-provider-snyk"
	comment := "https://registry.terraform.io/providers/pavel-snyk/snyk"
	return fmt.Sprintf("%s/%s (+%s)", name, p.version, comment)
}
