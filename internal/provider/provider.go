package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/pavel-snyk/snyk-sdk-go/snyk"
)

var _ provider.Provider = (*snykProvider)(nil)

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

func (p *snykProvider) Schema(_ context.Context, request provider.SchemaRequest, response *provider.SchemaResponse) {
	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"endpoint": schema.StringAttribute{
				MarkdownDescription: "This can be used to override the base URL for Snyk API requests. It can be also sourced from the `SNYK_ENDPOINT` environment variable.",
				Optional:            true,
			},
			"token": schema.StringAttribute{
				MarkdownDescription: "This is the API token from Snyk. It must be provided, but it can also be sourced from the `SNYK_TOKEN` environment variable.",
				Optional:            true,
				Sensitive:           true,
			},
		},
	}
}

type providerData struct {
	Endpoint         types.String `tfsdk:"endpoint"`
	Token            types.String `tfsdk:"token"`
	terraformVersion string
}

func (p *snykProvider) Configure(ctx context.Context, request provider.ConfigureRequest, response *provider.ConfigureResponse) {
	var config providerData

	diags := request.Config.Get(ctx, &config)
	response.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	config.terraformVersion = request.TerraformVersion

	// fallback to env if unset
	//if config.Endpoint.IsNull() {
	//	config.Endpoint.Value = os.Getenv("SNYK_ENDPOINT")
	//}
	//if config.Token.Null {
	//	config.Token.Value = os.Getenv("SNYK_TOKEN")
	//}

	// required if still unset
	//if config.Token.Value == "" {
	//	response.Diagnostics.AddAttributeError(
	//		path.Root("token"),
	//		"Invalid provider config",
	//		"token must be set.",
	//	)
	//	return
	//}

	opts := []snyk.ClientOption{snyk.WithUserAgent(p.userAgent())}
	//if config.Endpoint.Value != "" {
	//	tflog.Info(ctx, "Overriding default endpoint", map[string]interface{}{
	//		"endpoint": config.Endpoint.Value,
	//	})
	//	opts = append(opts, snyk.WithBaseURL(config.Endpoint.Value))
	//}

	client := snyk.NewClient(config.Token.ValueString(), opts...)

	response.DataSourceData = client
	response.ResourceData = client
}

func (p *snykProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		//NewOrganizationDataSource,
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
