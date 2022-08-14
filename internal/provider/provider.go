package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/pavel-snyk/snyk-sdk-go/snyk"
)

func New() provider.Provider {
	return &snykProvider{}
}

type snykProvider struct {
	client *snyk.Client
}

func (p *snykProvider) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		//
	}, nil
}

type providerData struct {
	Token types.String `tfsdk:"token"`
}

func (p *snykProvider) Configure(ctx context.Context, request provider.ConfigureRequest, response *provider.ConfigureResponse) {
	var config providerData

	diags := request.Config.Get(ctx, &config)
	if diags.HasError() {
		response.Diagnostics.Append(diags...)
		return
	}

	client := snyk.NewClient(config.Token.Value)

	p.client = client
}

func (p *snykProvider) GetResources(_ context.Context) (map[string]provider.ResourceType, diag.Diagnostics) {
	return nil, nil
}

func (p *snykProvider) GetDataSources(_ context.Context) (map[string]provider.DataSourceType, diag.Diagnostics) {
	return nil, nil
}
