package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/pavel-snyk/terraform-provider-snyk/internal/validators"
)

type projectDataSourceType struct{}

func (d projectDataSourceType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Description: "The project data source provides information about an existing Snyk project.",
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Description: "The ID of the project.",
				Computed:    true,
				Type:        types.StringType,
			},
			"name": {
				Description: "The name of the project.",
				Required:    true,
				Validators: []tfsdk.AttributeValidator{
					validators.NotEmptyString(),
				},
				Type: types.StringType,
			},
			"organization_id": {
				Description: "The ID of the organization that the project belongs to.",
				Required:    true,
				Validators: []tfsdk.AttributeValidator{
					validators.NotEmptyString(),
				},
				Type: types.StringType,
			},
		},
	}, nil
}

func (d projectDataSourceType) NewDataSource(_ context.Context, p provider.Provider) (datasource.DataSource, diag.Diagnostics) {
	return projectDataSource{
		p: p.(*snykProvider),
	}, nil
}

type projectDataSource struct {
	p *snykProvider
}

type projectData struct {
	ID             types.String `tfsdk:"id"`
	Name           types.String `tfsdk:"name"`
	OrganizationID types.String `tfsdk:"organization_id"`
}

func (d projectDataSource) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var data projectData
	diags := request.Config.Get(ctx, &data)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	projects, _, err := d.p.client.Projects.List(ctx, data.OrganizationID.Value)
	if err != nil {
		response.Diagnostics.AddError("Error getting projects", err.Error())
		return
	}
	for _, project := range projects {
		if project.Name == data.Name.Value {
			data.ID = types.String{Value: project.ID}
			data.Name = types.String{Value: project.Name}
			break
		}
	}

	diags = response.State.Set(ctx, &data)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
}
