package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/pavel-snyk/snyk-sdk-go/snyk"
)

type organizationResourceType struct{}

func (r organizationResourceType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Description: "The organization resource allows you to manage Snyk organization.",
		Attributes: map[string]tfsdk.Attribute{
			"group_id": {
				Description:   "The ID of the group that contains the organization.",
				Optional:      true,
				PlanModifiers: tfsdk.AttributePlanModifiers{resource.RequiresReplace()},
				Type:          types.StringType,
			},
			"id": {
				Description: "The ID of the organization.",
				Computed:    true,
				Type:        types.StringType,
			},
			"name": {
				Description: "The name of the organization.",
				Required:    true,
				Type:        types.StringType,
			},
		},
	}, nil
}

func (r organizationResourceType) NewResource(_ context.Context, p provider.Provider) (resource.Resource, diag.Diagnostics) {
	return organizationResource{
		p: p.(*snykProvider),
	}, nil
}

type organizationResource struct {
	p *snykProvider
}

type organizationData struct {
	GroupID types.String `tfsdk:"group_id"`
	ID      types.String `tfsdk:"id"`
	Name    types.String `tfsdk:"name"`
}

func (r organizationResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	// get values from plan
	var plan organizationData
	if diags := request.Plan.Get(ctx, &plan); diags.HasError() {
		response.Diagnostics.Append(diags...)
		return
	}

	// prepare create call
	createRequest := &snyk.OrganizationCreateRequest{
		Name:    plan.Name.Value,
		GroupID: plan.GroupID.Value,
	}
	tflog.Trace(ctx, "organizationResource.Create", map[string]interface{}{"payload": createRequest})
	org, _, err := r.p.client.Orgs.Create(ctx, createRequest)
	if err != nil {
		response.Diagnostics.AddError("Error creating organization", err.Error())
		return
	}
	tflog.Trace(ctx, "organizationResource.Create", map[string]interface{}{"organization_id": org.ID})

	// map response body to attributes
	result := &organizationData{
		GroupID: types.String{Value: org.Group.ID},
		ID:      types.String{Value: org.ID},
		Name:    types.String{Value: org.Name},
	}

	diags := response.State.Set(ctx, result)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (r organizationResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var state organizationData
	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	organizationID := state.ID.Value

	orgs, _, err := r.p.client.Orgs.List(ctx)
	if err != nil {
		response.Diagnostics.AddError("Error reading organizations", err.Error())
		return
	}

	var organization *snyk.Organization
	for _, org := range orgs {
		if org.ID == organizationID {
			organization = &org
			break
		}
	}
	if organization == nil {
		response.Diagnostics.AddError(
			"Error getting organization",
			"Could not find organization with ID "+organizationID+": "+err.Error(),
		)
		return
	}

	state.GroupID = types.String{Value: organization.Group.ID}
	state.ID = types.String{Value: organization.ID}
	state.Name = types.String{Value: organization.Name}

	diags = response.State.Set(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (r organizationResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
}

func (r organizationResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var state organizationData
	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	organizationID := state.ID.Value

	_, err := r.p.client.Orgs.Delete(ctx, organizationID)
	if err != nil {
		response.Diagnostics.AddError("Error deleting organization", err.Error())
		return
	}

	response.State.RemoveResource(ctx)
}
