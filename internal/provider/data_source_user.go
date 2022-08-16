package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type userDataSourceType struct{}

func (d userDataSourceType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Description: "The user data source provides information about an existing Snyk user.",
		Attributes: map[string]tfsdk.Attribute{
			"email": {
				Description: "The email of the user.",
				Type:        types.StringType,
				Computed:    true,
			},
			"id": {
				Description: "The ID of the user.",
				Type:        types.StringType,
				Computed:    true,
				Optional:    true,
			},
			"name": {
				Description: "The name of the user.",
				Type:        types.StringType,
				Computed:    true,
			},
			"username": {
				Description: "The username of the user.",
				Type:        types.StringType,
				Computed:    true,
			},
		},
	}, nil
}

func (d userDataSourceType) NewDataSource(_ context.Context, p provider.Provider) (datasource.DataSource, diag.Diagnostics) {
	return userDataSource{
		p: p.(*snykProvider),
	}, nil
}

type userDataSource struct {
	p *snykProvider
}

type userData struct {
	Email    types.String `tfsdk:"email"`
	ID       types.String `tfsdk:"id"`
	Name     types.String `tfsdk:"name"`
	Username types.String `tfsdk:"username"`
}

func (d userDataSource) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var uData userData
	if diags := request.Config.Get(ctx, &uData); diags.HasError() {
		response.Diagnostics = diags
		return
	}

	if uData.ID.Null {
		currentUser, _, err := d.p.client.Users.GetCurrent(ctx)
		if err != nil {
			response.Diagnostics.AddError("Error retrieving user", err.Error())
			return
		}

		uData.ID = types.String{Value: currentUser.ID}
	}

	// workaround to get the name (not username!)
	user, _, err := d.p.client.Users.Get(ctx, uData.ID.Value)
	if err != nil {
		response.Diagnostics.AddError("Error retrieving user", err.Error())
		return
	}

	uData.Email = types.String{Value: user.Email}
	uData.Name = types.String{Value: user.Name}
	uData.Username = types.String{Value: user.Username}

	diags := response.State.Set(ctx, uData)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
}
