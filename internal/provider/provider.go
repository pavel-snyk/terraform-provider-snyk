package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/pavel-snyk/snyk-sdk-go/snyk"
)

const (
	defaultEndpoint = "https://snyk.io/api/"
)

func New() provider.Provider {
	return &snykProvider{}
}

type snykProvider struct {
	client *snyk.Client
}

func (p *snykProvider) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Version: 1,
		Attributes: map[string]tfsdk.Attribute{
			"endpoint": {
				Description: "This can be used to override the base URL for Snyk API requests.",
				Type:        types.StringType,
				Optional:    true,
			},
			"token": {
				Description: "This is the API token from Snyk. It must be provided, but it can also be sourced from the `SNYK_TOKEN` environment variable",
				Type:        types.StringType,
				Optional:    true,
				Sensitive:   true,
			},
		},
	}, nil
}

type providerData struct {
	Endpoint         types.String `tfsdk:"endpoint"`
	Token            types.String `tfsdk:"token"`
	terraformVersion string
}

func (p *snykProvider) Configure(ctx context.Context, request provider.ConfigureRequest, response *provider.ConfigureResponse) {
	var config providerData

	diags := request.Config.Get(ctx, &config)
	if diags.HasError() {
		response.Diagnostics.Append(diags...)
		return
	}

	config.terraformVersion = request.TerraformVersion

	client := newSnykClient(&config)

	p.client = client
}

func (p *snykProvider) GetResources(_ context.Context) (map[string]provider.ResourceType, diag.Diagnostics) {
	return map[string]provider.ResourceType{
		"snyk_organization": organizationResourceType{},
	}, nil
}

func (p *snykProvider) GetDataSources(_ context.Context) (map[string]provider.DataSourceType, diag.Diagnostics) {
	return map[string]provider.DataSourceType{
		"snyk_user": userDataSourceType{},
	}, nil
}

func newSnykClient(pd *providerData) *snyk.Client {
	if pd.Endpoint.Null {
		pd.Endpoint = types.String{Value: defaultEndpoint}
	}
	return snyk.NewClient(pd.Token.Value,
		snyk.WithBaseURL(pd.Endpoint.Value),
		snyk.WithUserAgent("terraform-provider-snyk/dev (+https://registry.terraform.io/providers/pavel-snyk/snyk)"),
	)
}
