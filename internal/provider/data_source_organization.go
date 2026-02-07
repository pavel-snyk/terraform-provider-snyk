package provider

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/pavel-snyk/snyk-sdk-go/v2/snyk"
)

var (
	_ datasource.DataSource              = &organizationDataSource{}
	_ datasource.DataSourceWithConfigure = &organizationDataSource{}
)

// organizationDataSource is the organization datasource implementation.
type organizationDataSource struct {
	client *snyk.Client
}

// organizationDataSourceModel maps the organization datasource schema data.
type organizationDataSourceModel struct {
	GroupID  types.String `tfsdk:"group_id"`
	ID       types.String `tfsdk:"id"`
	Name     types.String `tfsdk:"name"`
	Slug     types.String `tfsdk:"slug"`
	TenantID types.String `tfsdk:"tenant_id"`
}

func NewOrganizationDataSource() datasource.DataSource {
	return &organizationDataSource{}
}

func (d *organizationDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = "snyk_organization"
}

func (d *organizationDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: `
The organization data source provides information about an existing Snyk organization.

An Organization in Snyk is a way to collect and organize your Projects. Members of Organizations
have access to these Projects. See [Manage Groups and Organizations](https://docs.snyk.io/snyk-platform-administration/groups-and-organizations).
`,
		Attributes: map[string]schema.Attribute{
			"group_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the group to which the organization belongs.",
				Computed:            true,
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "The ID of the organization.",
				Computed:            true,
				Optional:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the organization.",
				Computed:            true,
				Optional:            true,
			},
			"slug": schema.StringAttribute{
				MarkdownDescription: "The canonical (unique and URL-friendly) name of the organization.",
				Computed:            true,
			},
			"tenant_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the tenant to which the organization belongs.",
				Computed:            true,
			},
		},
	}
}

func (d *organizationDataSource) Configure(_ context.Context, request datasource.ConfigureRequest, response *datasource.ConfigureResponse) {
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

func (d *organizationDataSource) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var data organizationDataSourceModel

	diags := request.Config.Get(ctx, &data)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	orgID := data.ID.ValueString()
	orgName := data.Name.ValueString()
	if orgID == "" && orgName == "" {
		response.Diagnostics.AddError(
			"Missing required attributes",
			`The attribute "id" or "name" must be defined.`,
		)
		return
	}

	if orgID != "" {
		tflog.Info(ctx, "Searching for organization by id", map[string]any{"organization_id": orgID})

		tflog.Debug(ctx, "Getting organization by id", map[string]any{"organization_id": orgID})
		organization, resp, err := d.client.Orgs.Get(ctx, orgID, nil)
		if err != nil {
			if resp != nil && resp.StatusCode == http.StatusNotFound {
				response.Diagnostics.AddError("No search results", "Please refine your search.")
				return
			}
			response.Diagnostics.AddError("Unable to get organization", err.Error())
			return
		}
		tflog.Debug(ctx, "Got organization", map[string]any{"data": organization})

		// if name is defined check that it's equal
		if orgName != "" && orgName != organization.Attributes.Name {
			response.Diagnostics.AddError(
				"Ambiguous search results",
				fmt.Sprintf("Specified and actual organization name are different: expected '%s', got '%s'", orgName, organization.Attributes.Name),
			)
			return
		}

		// map response body to attributes
		data.GroupID = types.StringValue(organization.Attributes.GroupID)
		data.ID = types.StringValue(organization.ID)
		data.Name = types.StringValue(organization.Attributes.Name)
		data.Slug = types.StringValue(organization.Attributes.Slug)
		if organization.Relationships != nil && organization.Relationships.Tenant != nil && organization.Relationships.Tenant.Data != nil {
			data.TenantID = types.StringValue(organization.Relationships.Tenant.Data.ID)
		}
	} else {
		tflog.Info(ctx, "Searching for organization by name", map[string]any{"organization_name": orgName})

		tflog.Debug(ctx, "Getting all accessible organizations", map[string]any{"organization_name": orgName})
		var organization *snyk.Organization
		orgs, errf := d.client.Orgs.AllAccessibleOrgs(ctx, nil)
		for org := range orgs {
			if orgName == org.Attributes.Name {
				tflog.Info(ctx, "Found organization by name", map[string]any{"organization_name": orgName, "data": org})
				organization = &org
				break
			}
		}
		if err := errf(); err != nil {
			response.Diagnostics.AddError("Unable to get organizations", err.Error())
			return
		}
		if organization == nil {
			response.Diagnostics.AddError("No search results", "Please refine your search.")
			return
		}

		// enrich data because not all fields are exposed via list API
		tflog.Debug(ctx, "Getting organization by id", map[string]any{"organization_id": organization.ID, "organization_name": orgName})
		organization, _, err := d.client.Orgs.Get(ctx, organization.ID, nil)
		if err != nil {
			response.Diagnostics.AddError("Unable to get organization", err.Error())
			return
		}
		tflog.Debug(ctx, "Got organization by id", map[string]any{"organization_id": organization.ID, "data": organization})

		// map response body to attributes
		data.GroupID = types.StringValue(organization.Attributes.GroupID)
		data.ID = types.StringValue(organization.ID)
		data.Name = types.StringValue(organization.Attributes.Name)
		data.Slug = types.StringValue(organization.Attributes.Slug)
		if organization.Relationships != nil && organization.Relationships.Tenant != nil && organization.Relationships.Tenant.Data != nil {
			data.TenantID = types.StringValue(organization.Relationships.Tenant.Data.ID)
		}
	}

	diags = response.State.Set(ctx, &data)
	response.Diagnostics.Append(diags...)
}
