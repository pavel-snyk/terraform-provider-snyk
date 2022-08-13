package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"

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

func (p *snykProvider) Configure(_ context.Context, _ provider.ConfigureRequest, _ *provider.ConfigureResponse) {
	//
}

func (p *snykProvider) GetResources(_ context.Context) (map[string]provider.ResourceType, diag.Diagnostics) {
	return nil, nil
}

func (p *snykProvider) GetDataSources(_ context.Context) (map[string]provider.DataSourceType, diag.Diagnostics) {
	return nil, nil
}
