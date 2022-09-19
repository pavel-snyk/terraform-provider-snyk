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

	"github.com/pavel-snyk/terraform-provider-snyk/internal/validators"
)

type integrationResourceType struct{}

func (r integrationResourceType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Description: "The integration resource allows you to manage Snyk integration.",
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Description: "The ID of the integration.",
				Computed:    true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					resource.UseStateForUnknown(),
				},
				Type: types.StringType,
			},
			"organization_id": {
				Description:   "The ID of the organization that the integration belongs to.",
				Required:      true,
				PlanModifiers: tfsdk.AttributePlanModifiers{resource.RequiresReplace()},
				Type:          types.StringType,
			},
			"password": {
				Description: "The password used by the integration.",
				Optional:    true,
				Sensitive:   true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					resource.UseStateForUnknown(),
				},
				Type: types.StringType,
			},
			"pull_request_sca": {
				Description: "The pull request testing configuration for SCA (Software Composition Analysis). Snyk will checks " +
					"projects imported through the SCM integration for security and license issues whenever a new PR is opened.",
				Optional: true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					resource.UseStateForUnknown(),
				},
				Attributes: tfsdk.SingleNestedAttributes(map[string]tfsdk.Attribute{
					"enabled": {
						Description: "Denotes the pull request SCA feature should be enabled for this integration.",
						Computed:    true,
						Optional:    true,
						PlanModifiers: tfsdk.AttributePlanModifiers{
							resource.UseStateForUnknown(),
						},
						Type: types.BoolType,
					},
					"fail_on_any_issue": {
						Description: "Fails an opened pull request if any vulnerable dependencies have been detected, otherwise " +
							"the pull request should only fail when a dependency with issues is added.",
						Computed: true,
						Optional: true,
						PlanModifiers: tfsdk.AttributePlanModifiers{
							resource.UseStateForUnknown(),
						},
						Type: types.BoolType,
					},
					"fail_only_for_high_and_critical_severity": {
						Description: "Fails an opened pull request if any dependencies are marked as being of high or critical severity.",
						Computed:    true,
						Optional:    true,
						PlanModifiers: tfsdk.AttributePlanModifiers{
							resource.UseStateForUnknown(),
						},
						Type: types.BoolType,
					},
					"fail_only_on_issues_with_fix": {
						Description: "Fails an opened pull request only when issues found have a fix available.",
						Computed:    true,
						Optional:    true,
						PlanModifiers: tfsdk.AttributePlanModifiers{
							resource.UseStateForUnknown(),
						},
						Type: types.BoolType,
					},
				}),
			},
			"region": {
				Description: "The region used by the integration.",
				Optional:    true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					resource.UseStateForUnknown(),
				},
				Type: types.StringType,
			},
			"registry_url": {
				Description: "The URL for container registries used by the integration (e.g. for ECR).",
				Optional:    true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					resource.UseStateForUnknown(),
				},
				Type: types.StringType,
			},
			"role_arn": {
				Description: "The role ARN used by the integration (ECR only).",
				Optional:    true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					resource.UseStateForUnknown(),
				},
				Type: types.StringType,
			},
			"token": {
				Description: "The token used by the integration.",
				Optional:    true,
				Sensitive:   true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					resource.UseStateForUnknown(),
				},
				Type: types.StringType,
			},
			"type": {
				Description:   "The integration type, e.g. 'github'.",
				Required:      true,
				PlanModifiers: tfsdk.AttributePlanModifiers{resource.RequiresReplace()},
				Validators: []tfsdk.AttributeValidator{
					validators.NotEmptyString(),
				},
				Type: types.StringType,
			},
			"url": {
				Description: "The URL used by the integration.",
				Optional:    true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					resource.UseStateForUnknown(),
				},
				Type: types.StringType,
			},
			"username": {
				Description: "The username used by the integration.",
				Optional:    true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					resource.UseStateForUnknown(),
				},
				Type: types.StringType,
			},
		},
	}, nil
}

func (r integrationResourceType) NewResource(_ context.Context, p provider.Provider) (resource.Resource, diag.Diagnostics) {
	return integrationResource{
		p: p.(*snykProvider),
	}, nil
}

type integrationResource struct {
	p *snykProvider
}

type integrationData struct {
	ID             types.String    `tfsdk:"id"`
	OrganizationID types.String    `tfsdk:"organization_id"`
	Password       types.String    `tfsdk:"password"`
	PullRequestSCA *pullRequestSCA `tfsdk:"pull_request_sca"`
	Region         types.String    `tfsdk:"region"`
	RegistryURL    types.String    `tfsdk:"registry_url"`
	RoleARN        types.String    `tfsdk:"role_arn"`
	Token          types.String    `tfsdk:"token"`
	Type           types.String    `tfsdk:"type"`
	URL            types.String    `tfsdk:"url"`
	Username       types.String    `tfsdk:"username"`
}

type pullRequestSCA struct {
	Enabled                            types.Bool `tfsdk:"enabled"`
	FailOnAnyIssue                     types.Bool `tfsdk:"fail_on_any_issue"`
	FailOnlyForHighAndCriticalSeverity types.Bool `tfsdk:"fail_only_for_high_and_critical_severity"`
	FailOnlyOnIssuesWithFix            types.Bool `tfsdk:"fail_only_on_issues_with_fix"`
}

func (r integrationResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var plan integrationData
	diags := request.Plan.Get(ctx, &plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	orgID := plan.OrganizationID.Value
	integrations, _, err := r.p.client.Integrations.List(ctx, orgID)
	if err != nil {
		response.Diagnostics.AddError("Error reading integrations", err.Error())
		return
	}
	integrationType := snyk.IntegrationType(plan.Type.Value)
	integrationID := ""
	for t, id := range integrations {
		if t == integrationType {
			integrationID = id
			break
		}
	}

	result := &integrationData{
		OrganizationID: plan.OrganizationID,
		Password:       plan.Password,
		Region:         plan.Region,
		RegistryURL:    plan.RegistryURL,
		RoleARN:        plan.RoleARN,
		Token:          plan.Token,
		Type:           plan.Type,
		URL:            plan.URL,
		Username:       plan.Username,
	}

	if integrationID == "" {
		// new integration
		createRequest := &snyk.IntegrationCreateRequest{
			Integration: &snyk.Integration{
				Credentials: &snyk.IntegrationCredentials{
					Password:     plan.Password.Value,
					Region:       plan.Region.Value,
					RegistryBase: plan.RegistryURL.Value,
					RoleARN:      plan.RoleARN.Value,
					Token:        plan.Token.Value,
					URL:          plan.URL.Value,
					Username:     plan.Username.Value,
				},
				Type: integrationType,
			},
		}
		tflog.Trace(ctx, "integrationResource.Create(create)", map[string]interface{}{"payload": createRequest})
		integration, _, err := r.p.client.Integrations.Create(ctx, orgID, createRequest)
		if err != nil {
			response.Diagnostics.AddError("Error creating integration", err.Error())
			return
		}
		result.ID = types.String{Value: integration.ID}
	} else {
		//update existing integration with new credentials
		updateRequest := &snyk.IntegrationUpdateRequest{
			Integration: &snyk.Integration{
				Credentials: &snyk.IntegrationCredentials{
					Password:     plan.Password.Value,
					Region:       plan.Region.Value,
					RegistryBase: plan.RegistryURL.Value,
					RoleARN:      plan.RoleARN.Value,
					Token:        plan.Token.Value,
					URL:          plan.URL.Value,
					Username:     plan.Username.Value,
				},
				Type: integrationType,
			},
		}
		tflog.Trace(ctx, "integrationResource.Create(update)", map[string]interface{}{"payload": updateRequest})
		integration, _, err := r.p.client.Integrations.Update(ctx, orgID, integrationID, updateRequest)
		if err != nil {
			response.Diagnostics.AddError("Error updating integration", err.Error())
			return
		}
		result.ID = types.String{Value: integration.ID}
	}

	if plan.PullRequestSCA != nil {
		updateRequest := &snyk.IntegrationSettingsUpdateRequest{
			IntegrationSettings: &snyk.IntegrationSettings{
				PullRequestTestEnabled:                        toBoolPtr(plan.PullRequestSCA.Enabled),
				PullRequestFailOnAnyVulnerability:             toBoolPtr(plan.PullRequestSCA.FailOnAnyIssue),
				PullRequestFailOnlyForHighAndCriticalSeverity: toBoolPtr(plan.PullRequestSCA.FailOnlyForHighAndCriticalSeverity),
				PullRequestFailOnlyForIssuesWithFix:           toBoolPtr(plan.PullRequestSCA.FailOnlyOnIssuesWithFix),
			},
		}

		settings, _, err := r.p.client.Integrations.UpdateSettings(ctx, orgID, result.ID.Value, updateRequest)
		if err != nil {
			response.Diagnostics.AddError("Error updating pull request settings", err.Error())
			return
		}
		result.PullRequestSCA = &pullRequestSCA{
			Enabled:                            fromBoolPtr(settings.PullRequestTestEnabled),
			FailOnAnyIssue:                     fromBoolPtr(settings.PullRequestFailOnAnyVulnerability),
			FailOnlyForHighAndCriticalSeverity: fromBoolPtr(settings.PullRequestFailOnlyForHighAndCriticalSeverity),
			FailOnlyOnIssuesWithFix:            fromBoolPtr(settings.PullRequestFailOnlyForIssuesWithFix),
		}
	}

	diags = response.State.Set(ctx, result)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (r integrationResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var state integrationData
	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	organizationID := state.OrganizationID.Value

	integration, _, err := r.p.client.Integrations.GetByType(ctx, organizationID, snyk.IntegrationType(state.Type.Value))
	if err != nil {
		response.Diagnostics.AddError("Error reading integration", err.Error())
		return
	}

	if state.PullRequestSCA != nil {
		settings, _, err := r.p.client.Integrations.GetSettings(ctx, organizationID, integration.ID)
		if err != nil {
			response.Diagnostics.AddError("Error reading integration settings", err.Error())
			return
		}

		pullRequestSCA := &pullRequestSCA{
			Enabled:                            fromBoolPtr(settings.PullRequestTestEnabled),
			FailOnAnyIssue:                     fromBoolPtr(settings.PullRequestFailOnAnyVulnerability),
			FailOnlyForHighAndCriticalSeverity: fromBoolPtr(settings.PullRequestFailOnlyForHighAndCriticalSeverity),
			FailOnlyOnIssuesWithFix:            fromBoolPtr(settings.PullRequestFailOnlyForIssuesWithFix),
		}
		state.PullRequestSCA = pullRequestSCA
	}

	state.ID = types.String{Value: integration.ID}

	diags = response.State.Set(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (r integrationResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var plan integrationData
	diags := request.Plan.Get(ctx, &plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	organizationID := plan.OrganizationID.Value
	integrationID := plan.ID.Value
	updateRequest := &snyk.IntegrationUpdateRequest{
		Integration: &snyk.Integration{
			Credentials: &snyk.IntegrationCredentials{
				Password:     plan.Password.Value,
				Region:       plan.Region.Value,
				RegistryBase: plan.RegistryURL.Value,
				RoleARN:      plan.RoleARN.Value,
				Token:        plan.Token.Value,
				URL:          plan.URL.Value,
				Username:     plan.Username.Value,
			},
			Type: snyk.IntegrationType(plan.Type.Value),
		},
	}
	tflog.Info(ctx, "integrationResource.Update", map[string]interface{}{"payload": updateRequest})
	integration, _, err := r.p.client.Integrations.Update(ctx, organizationID, integrationID, updateRequest)
	if err != nil {
		response.Diagnostics.AddError("Error updating integration", err.Error())
		return
	}

	if plan.PullRequestSCA != nil {
		updateRequest := &snyk.IntegrationSettingsUpdateRequest{
			IntegrationSettings: &snyk.IntegrationSettings{
				PullRequestTestEnabled:                        toBoolPtr(plan.PullRequestSCA.Enabled),
				PullRequestFailOnAnyVulnerability:             toBoolPtr(plan.PullRequestSCA.FailOnAnyIssue),
				PullRequestFailOnlyForHighAndCriticalSeverity: toBoolPtr(plan.PullRequestSCA.FailOnlyForHighAndCriticalSeverity),
				PullRequestFailOnlyForIssuesWithFix:           toBoolPtr(plan.PullRequestSCA.FailOnlyOnIssuesWithFix),
			},
		}

		settings, _, err := r.p.client.Integrations.UpdateSettings(ctx, organizationID, integrationID, updateRequest)
		if err != nil {
			response.Diagnostics.AddError("Error updating pull request settings", err.Error())
			return
		}
		plan.PullRequestSCA = &pullRequestSCA{
			Enabled:                            fromBoolPtr(settings.PullRequestTestEnabled),
			FailOnAnyIssue:                     fromBoolPtr(settings.PullRequestFailOnAnyVulnerability),
			FailOnlyForHighAndCriticalSeverity: fromBoolPtr(settings.PullRequestFailOnlyForHighAndCriticalSeverity),
			FailOnlyOnIssuesWithFix:            fromBoolPtr(settings.PullRequestFailOnlyForIssuesWithFix),
		}
	}

	plan.ID = types.String{Value: integration.ID}

	diags = response.State.Set(ctx, &plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (r integrationResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var state integrationData
	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	integrationID := state.ID.Value
	organizationID := state.OrganizationID.Value

	_, err := r.p.client.Integrations.DeleteCredentials(ctx, organizationID, integrationID)
	if err != nil {
		response.Diagnostics.AddError("Error deleting integration", err.Error())
		return
	}

	response.State.RemoveResource(ctx)
}
