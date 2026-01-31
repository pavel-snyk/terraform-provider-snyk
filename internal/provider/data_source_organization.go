package provider

//import (
//	"context"
//
//	"github.com/hashicorp/terraform-plugin-framework/datasource"
//	"github.com/hashicorp/terraform-plugin-framework/diag"
//	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
//	"github.com/hashicorp/terraform-plugin-framework/types"
//
//	"github.com/pavel-snyk/snyk-sdk-go/snyk"
//)
//
//var _ datasource.DataSource = (*organizationDataSource)(nil)
//
//type organizationDataSource struct {
//	client *snyk.Client
//}
//
//func NewOrganizationDataSource() datasource.DataSource {
//	return &organizationDataSource{}
//}
//
//func (d *organizationDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, response *datasource.MetadataResponse) {
//	response.TypeName = "snyk_organization"
//}
//
//func (d *organizationDataSource) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
//	return tfsdk.Schema{
//		Description: "The organization data source provides information about an existing Snyk organization.",
//		Attributes: map[string]tfsdk.Attribute{
//			"group_id": {
//				Description: "The ID of the group that contains the organization.",
//				Optional:    true,
//				Type:        types.StringType,
//			},
//			"id": {
//				Description: "The ID of the organization.",
//				Type:        types.StringType,
//				Required:    true,
//			},
//			"name": {
//				Description: "The name of the organization.",
//				Type:        types.StringType,
//				Computed:    true,
//			},
//		},
//	}, nil
//}
//
//func (d *organizationDataSource) Configure(_ context.Context, request datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
//	if request.ProviderData == nil {
//		return
//	}
//
//	client := request.ProviderData.(*snyk.Client)
//	d.client = client
//}
//
//func (d *organizationDataSource) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
//	var oData organizationData
//	diags := request.Config.Get(ctx, &oData)
//	response.Diagnostics.Append(diags...)
//	if response.Diagnostics.HasError() {
//		return
//	}
//
//	orgs, _, err := d.client.Orgs.List(ctx)
//	if err != nil {
//		response.Diagnostics.AddError("Error retrieving organizations", err.Error())
//		return
//	}
//
//	var organization *snyk.Organization
//	for _, org := range orgs {
//		if org.ID == oData.ID.Value {
//			organization = &org
//			break
//		}
//	}
//
//	oData.Name = types.String{Value: organization.Name}
//	oData.GroupID = types.String{Value: organization.Group.ID}
//
//	diags = response.State.Set(ctx, &oData)
//	response.Diagnostics.Append(diags...)
//	if response.Diagnostics.HasError() {
//		return
//	}
//}
