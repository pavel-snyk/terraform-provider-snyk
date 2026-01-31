package provider

//import (
//	"context"
//
//	"github.com/hashicorp/terraform-plugin-framework/datasource"
//	"github.com/hashicorp/terraform-plugin-framework/diag"
//	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
//	"github.com/hashicorp/terraform-plugin-framework/types"
//	"github.com/pavel-snyk/snyk-sdk-go/snyk"
//
//	"github.com/pavel-snyk/terraform-provider-snyk/internal/validator"
//)
//
//var _ datasource.DataSource = (*projectDataSource)(nil)
//
//type projectDataSource struct {
//	client *snyk.Client
//}
//
//func NewProjectDataSource() datasource.DataSource {
//	return &projectDataSource{}
//}
//
//func (d *projectDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, response *datasource.MetadataResponse) {
//	response.TypeName = "snyk_project"
//}
//
//func (d *projectDataSource) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
//	return tfsdk.Schema{
//		Description: "The project data source provides information about an existing Snyk project.",
//		Attributes: map[string]tfsdk.Attribute{
//			"id": {
//				Description: "The ID of the project.",
//				Computed:    true,
//				Type:        types.StringType,
//			},
//			"name": {
//				Description: "The name of the project.",
//				Required:    true,
//				Validators: []tfsdk.AttributeValidator{
//					validator.NotEmptyString(),
//				},
//				Type: types.StringType,
//			},
//			"organization_id": {
//				Description: "The ID of the organization that the project belongs to.",
//				Required:    true,
//				Validators: []tfsdk.AttributeValidator{
//					validator.NotEmptyString(),
//				},
//				Type: types.StringType,
//			},
//		},
//	}, nil
//}
//
//func (d *projectDataSource) Configure(_ context.Context, request datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
//	if request.ProviderData == nil {
//		return
//	}
//
//	client := request.ProviderData.(*snyk.Client)
//	d.client = client
//}
//
//type projectData struct {
//	ID             types.String `tfsdk:"id"`
//	Name           types.String `tfsdk:"name"`
//	OrganizationID types.String `tfsdk:"organization_id"`
//}
//
//func (d *projectDataSource) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
//	var data projectData
//	diags := request.Config.Get(ctx, &data)
//	response.Diagnostics.Append(diags...)
//	if response.Diagnostics.HasError() {
//		return
//	}
//
//	projects, _, err := d.client.Projects.List(ctx, data.OrganizationID.Value)
//	if err != nil {
//		response.Diagnostics.AddError("Error getting projects", err.Error())
//		return
//	}
//	for _, project := range projects {
//		if project.Name == data.Name.Value {
//			data.ID = types.String{Value: project.ID}
//			data.Name = types.String{Value: project.Name}
//			break
//		}
//	}
//
//	diags = response.State.Set(ctx, &data)
//	response.Diagnostics.Append(diags...)
//	if response.Diagnostics.HasError() {
//		return
//	}
//}
