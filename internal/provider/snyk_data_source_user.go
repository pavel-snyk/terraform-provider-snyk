package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/pavel-snyk/snyk-sdk-go/v2/snyk"
)

var (
	_ datasource.DataSource              = (*userDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*userDataSource)(nil)
)

// userDataSource is the user datasource implementation.
type userDataSource struct {
	client *snyk.Client
}

// userDataSourceModel maps the user datasource schema data.
type userDataSourceModel struct {
	Email    types.String `tfsdk:"email"`
	ID       types.String `tfsdk:"id"`
	Name     types.String `tfsdk:"name"`
	Username types.String `tfsdk:"username"`
}

func NewUserDataSource() datasource.DataSource {
	return &userDataSource{}
}

func (d *userDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = "snyk_user"
}

func (d *userDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `
The user data source provides information about an existing Snyk user.

A user in Snyk is a member of an Organization, that have access to Projects.
See [Manage users in Organizations](https://docs.snyk.io/snyk-platform-administration/groups-and-organizations/organizations/manage-users-in-organizations).
`,
		Attributes: map[string]schema.Attribute{
			"email": schema.StringAttribute{
				MarkdownDescription: "The email of the user.",
				Computed:            true,
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "The ID of the user.",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the user.",
				Computed:            true,
			},
			"username": schema.StringAttribute{
				MarkdownDescription: "The username of the user.",
				Computed:            true,
			},
		},
	}
}

func (d *userDataSource) Configure(_ context.Context, request datasource.ConfigureRequest, response *datasource.ConfigureResponse) {
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

func (d *userDataSource) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var data userDataSourceModel

	response.Diagnostics.Append(request.Config.Get(ctx, &data)...)
	if response.Diagnostics.HasError() {
		return
	}

	tflog.Trace(ctx, "Getting self user")
	user, resp, err := d.client.Users.GetSelf(ctx)
	if err != nil {
		response.Diagnostics.AddError("Unable to get user", err.Error())
		return
	}
	tflog.Trace(ctx, "Got self user", map[string]any{
		"data":            user,
		"snyk_request_id": resp.SnykRequestID,
	})

	// map response body to model
	data.Email = types.StringValue(user.Attributes.Email)
	data.ID = types.StringValue(user.ID)
	data.Name = types.StringValue(user.Attributes.Name)
	data.Username = types.StringValue(user.Attributes.Username)

	response.Diagnostics.Append(response.State.Set(ctx, &data)...)
}
