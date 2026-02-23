package helper

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/pavel-snyk/snyk-sdk-go/v2/snyk"
)

var _ validator.Object = (*requiresValidConfiguration)(nil)

type requiresValidConfiguration struct {
	connectionType path.Expression
}

func (av requiresValidConfiguration) Description(ctx context.Context) string {
	return av.MarkdownDescription(ctx)
}

func (av requiresValidConfiguration) MarkdownDescription(_ context.Context) string {
	return fmt.Sprintf("Ensure that if a connection type '%s' is set, all of its required attributes also configured.", av.connectionType)
}

func (av requiresValidConfiguration) ValidateObject(ctx context.Context, request validator.ObjectRequest, response *validator.ObjectResponse) {
	matchedPaths, diags := request.Config.PathMatches(ctx, av.connectionType)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	if len(matchedPaths) != 1 {
		response.Diagnostics.Append(diag.NewAttributeErrorDiagnostic(
			request.Path,
			"Invalid Broker Connection",
			fmt.Sprintf("Expected a single attribute path for the connection type, but got %d paths.", len(matchedPaths)),
		))
		return
	}

	var connectionTypeVal types.String
	response.Diagnostics.Append(request.Config.GetAttribute(ctx, matchedPaths[0], &connectionTypeVal)...)
	if connectionTypeVal.IsNull() || connectionTypeVal.IsUnknown() {
		return
	}
	connectionType := connectionTypeVal.ValueString()

	allowedConnectionTypes := allowedConnectionTypes()
	slices.Sort(allowedConnectionTypes)
	if !slices.Contains(allowedConnectionTypes, connectionType) {
		response.Diagnostics.Append(diag.NewAttributeErrorDiagnostic(
			matchedPaths[0],
			"Invalid Broker Connection",
			fmt.Sprintf("Unexpected connection type: \"%s\".\n"+
				"Allowed connection types: %v", connectionType, strings.Join(allowedConnectionTypes, ", ")),
		))
		return
	}

	switch connectionType {
	case "gitlab":
		av.validateGitLabConfiguration(request.ConfigValue, request.Path, response)
	case "jira":
		av.validateJiraConfiguration(request.ConfigValue, request.Path, response)
	}
}

func (av requiresValidConfiguration) validateGitLabConfiguration(config types.Object, path path.Path, response *validator.ObjectResponse) {
	requiredAttributes := []string{"broker_client_url", "gitlab_hostname", "gitlab_token_credential_id"}
	for _, attributeName := range requiredAttributes {
		value, ok := config.Attributes()[attributeName]
		if !ok || value.IsNull() {
			response.Diagnostics.AddAttributeError(
				path,
				"Missing required attribute for GitLab",
				fmt.Sprintf(`Attribute "%s" is required when connection type is "gitlab".`, attributeName),
			)
		}
	}
}

func (av requiresValidConfiguration) validateJiraConfiguration(config types.Object, path path.Path, response *validator.ObjectResponse) {
	jiraHostname, ok := config.Attributes()["jira_hostname"]
	if !ok || jiraHostname.IsNull() {
		response.Diagnostics.AddAttributeError(
			path,
			"Missing required attribute for Jira",
			`Attribute "jira_hostname" is required when connection type is "jira".`,
		)
	}

	// only PAT or username/password can be configured
	jiraPAT, patOk := config.Attributes()["jira_pat_credential_id"]
	jiraUsername, usernameOk := config.Attributes()["jira_username"]
	jiraPassword, passwordOk := config.Attributes()["jira_password_credential_id"]

	patIsSet := patOk && !jiraPAT.IsNull()
	usernameIsSet := usernameOk && !jiraUsername.IsNull()
	passwordIsSet := passwordOk && !jiraPassword.IsNull()

	if patIsSet && (usernameIsSet || passwordIsSet) {
		response.Diagnostics.AddAttributeError(
			path,
			"Conflicting Jira authentication attributes",
			`"jira_pat_credential_id" cannot be used with "jira_username" or "jira_password_credential_id".
Please provide either the PAT or the username/password combination.`,
		)
	}

	if !patIsSet && (!usernameIsSet || !passwordIsSet) {
		response.Diagnostics.AddAttributeError(
			path,
			"Incomplete Jira authentication attributes",
			`When connection type is "jira", you must provide either "jira_pat_credential_id" or both "jira_username" and "jira_password_credential_id".`,
		)
	}
}

// RequiresValidConfiguration checks that a path.Expression matches a valid
// brokerConnectionResourceModel type and all requires attributes from "configuration"
// are set correctly for the specified type.
func RequiresValidConfiguration(connectionType path.Expression) validator.Object {
	return &requiresValidConfiguration{connectionType}
}

func allowedConnectionTypes() []string {
	return []string{
		string(snyk.BrokerConnectionTypeACR),
		string(snyk.BrokerConnectionTypeArtifactory),
		string(snyk.BrokerConnectionTypeArtifactoryCR),
		string(snyk.BrokerConnectionTypeAzureRepos),
		string(snyk.BrokerConnectionTypeBitbucketServer),
		string(snyk.BrokerConnectionTypeDigitaloceanCR),
		string(snyk.BrokerConnectionTypeDockerHub),
		string(snyk.BrokerConnectionTypeECR),
		string(snyk.BrokerConnectionTypeGCR),
		string(snyk.BrokerConnectionTypeGitHub),
		string(snyk.BrokerConnectionTypeGitHubCloudApp),
		string(snyk.BrokerConnectionTypeGitHubCR),
		string(snyk.BrokerConnectionTypeGitHubEnterprise),
		string(snyk.BrokerConnectionTypeGitHubServerApp),
		string(snyk.BrokerConnectionTypeGitLab),
		string(snyk.BrokerConnectionTypeGitLabCR),
		string(snyk.BrokerConnectionTypeGoogleArtifactCR),
		string(snyk.BrokerConnectionTypeJira),
		string(snyk.BrokerConnectionTypeHarborCR),
		string(snyk.BrokerConnectionTypeNexus),
		string(snyk.BrokerConnectionTypeNexusCR),
		string(snyk.BrokerConnectionTypeQuayCR),
	}
}
