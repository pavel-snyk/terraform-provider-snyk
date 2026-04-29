package helper

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

func TestRequiresValidConfiguration_validateBitbucketConfiguration(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		request              validator.ObjectRequest
		expectedErrorsCount  int
		expectedErrorDetails []string
	}{
		"basic-auth": {
			request: validator.ObjectRequest{
				ConfigValue: types.ObjectValueMust(
					map[string]attr.Type{
						"broker_client_url":                types.StringType,
						"bitbucket_hostname":               types.StringType,
						"bitbucket_password_credential_id": types.StringType,
						"bitbucket_username":               types.StringType,
					},
					map[string]attr.Value{
						"broker_client_url":                types.StringValue("http://localhost:8080"),
						"bitbucket_hostname":               types.StringValue("bitbucket.test"),
						"bitbucket_password_credential_id": types.StringValue("random-uuid"),
						"bitbucket_username":               types.StringValue("bitbucket-user"),
					},
				),
			},
			expectedErrorsCount:  0,
			expectedErrorDetails: []string{},
		},
		"pat-auth": {
			request: validator.ObjectRequest{
				ConfigValue: types.ObjectValueMust(
					map[string]attr.Type{
						"broker_client_url":           types.StringType,
						"bitbucket_hostname":          types.StringType,
						"bitbucket_pat_credential_id": types.StringType,
					},
					map[string]attr.Value{
						"broker_client_url":           types.StringValue("http://localhost:8080"),
						"bitbucket_hostname":          types.StringValue("bitbucket.test"),
						"bitbucket_pat_credential_id": types.StringValue("random-uuid"),
					},
				),
			},
			expectedErrorsCount:  0,
			expectedErrorDetails: []string{},
		},
		"missing-required-attributes": {
			request: validator.ObjectRequest{
				ConfigValue: types.ObjectValueMust(
					map[string]attr.Type{
						"broker_client_url":                types.StringType,
						"bitbucket_hostname":               types.StringType,
						"bitbucket_password_credential_id": types.StringType,
						"bitbucket_username":               types.StringType,
					},
					map[string]attr.Value{
						"broker_client_url":                types.StringNull(),
						"bitbucket_hostname":               types.StringNull(),
						"bitbucket_password_credential_id": types.StringNull(),
						"bitbucket_username":               types.StringNull(),
					},
				),
			},
			expectedErrorsCount: 3,
			expectedErrorDetails: []string{
				`Attribute "broker_client_url" is required when connection type is "bitbucket-server".`,
				`Attribute "bitbucket_hostname" is required when connection type is "bitbucket-server".`,
				`When connection type is "bitbucket-server", you must provide either "bitbucket_pat_credential_id" or both "bitbucket_username" and "bitbucket_password_credential_id".`,
			},
		},
		"pat-with-basic-together": {
			request: validator.ObjectRequest{
				ConfigValue: types.ObjectValueMust(
					map[string]attr.Type{
						"broker_client_url":                types.StringType,
						"bitbucket_hostname":               types.StringType,
						"bitbucket_password_credential_id": types.StringType,
						"bitbucket_pat_credential_id":      types.StringType,
						"bitbucket_username":               types.StringType,
					},
					map[string]attr.Value{
						"broker_client_url":                types.StringValue("http://localhost:8080"),
						"bitbucket_hostname":               types.StringValue("bitbucket.test"),
						"bitbucket_password_credential_id": types.StringValue("random-uuid"),
						"bitbucket_pat_credential_id":      types.StringValue("random-uuid"),
						"bitbucket_username":               types.StringValue("bitbucket-user"),
					},
				),
			},
			expectedErrorsCount: 1,
			expectedErrorDetails: []string{
				`"bitbucket_pat_credential_id" cannot be used with "bitbucket_username" or "bitbucket_password_credential_id".
Please provide either the PAT or the username/password combination.`,
			},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			response := &validator.ObjectResponse{}
			requiresValidConfiguration{}.validateBitbucketConfiguration(test.request.ConfigValue, path.Empty(), response)
			var diagnosticDetails []string
			for _, diagnostic := range response.Diagnostics {
				diagnosticDetails = append(diagnosticDetails, diagnostic.Detail())
			}

			if test.expectedErrorsCount > 0 {
				assert.True(t, response.Diagnostics.HasError())
				assert.Equal(t, test.expectedErrorsCount, len(diagnosticDetails))
				for _, expectedErrorDetail := range test.expectedErrorDetails {
					assert.Contains(t, diagnosticDetails, expectedErrorDetail)
				}
			} else {
				assert.False(t, response.Diagnostics.HasError())
			}
		})
	}
}

func TestRequiresValidConfiguration_validateGitLabConfiguration(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		request              validator.ObjectRequest
		expectedErrorsCount  int
		expectedErrorDetails []string
	}{
		"base": {
			request: validator.ObjectRequest{
				ConfigValue: types.ObjectValueMust(
					map[string]attr.Type{
						"broker_client_url":          types.StringType,
						"gitlab_hostname":            types.StringType,
						"gitlab_token_credential_id": types.StringType,
					},
					map[string]attr.Value{
						"broker_client_url":          types.StringValue("http://localhost:8080"),
						"gitlab_hostname":            types.StringValue("gitlab.com"),
						"gitlab_token_credential_id": types.StringValue("random-uuid"),
					},
				),
			},
			expectedErrorsCount:  0,
			expectedErrorDetails: []string{},
		},
		"missing-required-attributes": {
			request: validator.ObjectRequest{
				ConfigValue: types.ObjectValueMust(
					map[string]attr.Type{
						"broker_client_url":          types.StringType,
						"gitlab_hostname":            types.StringType,
						"gitlab_token_credential_id": types.StringType,
					},
					map[string]attr.Value{
						"broker_client_url":          types.StringNull(),
						"gitlab_hostname":            types.StringNull(),
						"gitlab_token_credential_id": types.StringNull(),
					},
				),
			},
			expectedErrorsCount: 3,
			expectedErrorDetails: []string{
				`Attribute "broker_client_url" is required when connection type is "gitlab".`,
				`Attribute "gitlab_hostname" is required when connection type is "gitlab".`,
				`Attribute "gitlab_token_credential_id" is required when connection type is "gitlab".`,
			},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			response := &validator.ObjectResponse{}
			requiresValidConfiguration{}.validateGitLabConfiguration(test.request.ConfigValue, path.Empty(), response)
			var diagnosticDetails []string
			for _, diagnostic := range response.Diagnostics {
				diagnosticDetails = append(diagnosticDetails, diagnostic.Detail())
			}

			if test.expectedErrorsCount > 0 {
				assert.True(t, response.Diagnostics.HasError())
				assert.Equal(t, test.expectedErrorsCount, len(diagnosticDetails))
				for _, expectedErrorDetail := range test.expectedErrorDetails {
					assert.Contains(t, diagnosticDetails, expectedErrorDetail)
				}
			} else {
				assert.False(t, response.Diagnostics.HasError())
			}
		})
	}
}

func TestRequiresValidConfiguration_validateJiraConfiguration(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		request              validator.ObjectRequest
		expectedErrorsCount  int
		expectedErrorDetails []string
	}{
		"basic-auth": {
			request: validator.ObjectRequest{
				ConfigValue: types.ObjectValueMust(
					map[string]attr.Type{
						"jira_hostname":               types.StringType,
						"jira_password_credential_id": types.StringType,
						"jira_username":               types.StringType,
					},
					map[string]attr.Value{
						"jira_hostname":               types.StringValue("https://jira.test"),
						"jira_password_credential_id": types.StringValue("random-uuid"),
						"jira_username":               types.StringValue("jira-user"),
					},
				),
			},
			expectedErrorsCount:  0,
			expectedErrorDetails: []string{},
		},
		"pat-auth": {
			request: validator.ObjectRequest{
				ConfigValue: types.ObjectValueMust(
					map[string]attr.Type{
						"jira_hostname":          types.StringType,
						"jira_pat_credential_id": types.StringType,
					},
					map[string]attr.Value{
						"jira_hostname":          types.StringValue("https://jira.test"),
						"jira_pat_credential_id": types.StringValue("random-uuid"),
					},
				),
			},
			expectedErrorsCount:  0,
			expectedErrorDetails: []string{},
		},
		"missing-required-attributes": {
			request: validator.ObjectRequest{
				ConfigValue: types.ObjectValueMust(
					map[string]attr.Type{
						"jira_hostname":          types.StringType,
						"jira_pat_credential_id": types.StringType,
					},
					map[string]attr.Value{
						"jira_hostname":          types.StringNull(),
						"jira_pat_credential_id": types.StringNull(),
					},
				),
			},
			expectedErrorsCount: 2,
			expectedErrorDetails: []string{
				`Attribute "jira_hostname" is required when connection type is "jira".`,
				`When connection type is "jira", you must provide either "jira_pat_credential_id" or both "jira_username" and "jira_password_credential_id".`,
			},
		},
		"pat-with-basic-together": {
			request: validator.ObjectRequest{
				ConfigValue: types.ObjectValueMust(
					map[string]attr.Type{
						"jira_hostname":               types.StringType,
						"jira_password_credential_id": types.StringType,
						"jira_pat_credential_id":      types.StringType,
						"jira_username":               types.StringType,
					},
					map[string]attr.Value{
						"jira_hostname":               types.StringValue("https://jira.test"),
						"jira_password_credential_id": types.StringValue("random-uuid"),
						"jira_pat_credential_id":      types.StringValue("random-uuid"),
						"jira_username":               types.StringValue("jira-user"),
					},
				),
			},
			expectedErrorsCount: 1,
			expectedErrorDetails: []string{
				`"jira_pat_credential_id" cannot be used with "jira_username" or "jira_password_credential_id".
Please provide either the PAT or the username/password combination.`,
			},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			response := &validator.ObjectResponse{}
			requiresValidConfiguration{}.validateJiraConfiguration(test.request.ConfigValue, path.Empty(), response)
			var diagnosticDetails []string
			for _, diagnostic := range response.Diagnostics {
				diagnosticDetails = append(diagnosticDetails, diagnostic.Detail())
			}

			if test.expectedErrorsCount > 0 {
				assert.True(t, response.Diagnostics.HasError())
				assert.Equal(t, test.expectedErrorsCount, len(diagnosticDetails))
				for _, expectedErrorDetail := range test.expectedErrorDetails {
					assert.Contains(t, diagnosticDetails, expectedErrorDetail)
				}
			} else {
				assert.False(t, response.Diagnostics.HasError())
			}
		})
	}
}
