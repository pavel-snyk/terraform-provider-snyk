package validator

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/pavel-snyk/snyk-sdk-go/snyk"
)

var _ tfsdk.AttributeValidator = requiredConfiguredCredentialsValidator{}

const (
	attributePassword    = "password"
	attributeRegion      = "region"
	attributeRegistryURL = "registry_url"
	attributeRoleARN     = "role_arn"
	attributeToken       = "token"
	attributeURL         = "url"
	attributeUsername    = "username"
)

type requiredConfiguredCredentialsValidator struct{}

func (validator requiredConfiguredCredentialsValidator) Description(_ context.Context) string {
	return "Ensure that integration is correctly configured"
}

func (validator requiredConfiguredCredentialsValidator) MarkdownDescription(ctx context.Context) string {
	return validator.Description(ctx)
}

func (validator requiredConfiguredCredentialsValidator) Validate(ctx context.Context, request tfsdk.ValidateAttributeRequest, response *tfsdk.ValidateAttributeResponse) {
	t := request.AttributeConfig.(types.String)

	if t.IsUnknown() {
		return
	}

	typePath := path.MatchRoot("type").String()
	configuredPath := request.AttributePathExpression.String()

	if typePath != configuredPath {
		response.Diagnostics.AddAttributeError(
			request.AttributePath,
			validator.Description(ctx),
			"Validator must be applied for integration 'type' only.",
		)
		return
	}

	switch t.Value {
	case snyk.ACRIntegrationType:
		username, diags := getAttributeValue(ctx, attributeUsername, request)
		response.Diagnostics.Append(diags...)
		if response.Diagnostics.HasError() {
			return
		}

		password, diags := getAttributeValue(ctx, attributePassword, request)
		response.Diagnostics.Append(diags...)
		if response.Diagnostics.HasError() {
			return
		}

		registryURL, diags := getAttributeValue(ctx, attributeRegistryURL, request)
		response.Diagnostics.Append(diags...)
		if response.Diagnostics.HasError() {
			return
		}

		if username.IsNull() || username.Value == "" {
			response.Diagnostics.AddAttributeError(
				request.AttributePath,
				validator.Description(ctx),
				fmt.Sprintf("%v must be defined and not empty for '%v' integration",
					attributeUsername,
					snyk.ACRIntegrationType,
				),
			)
		}

		if password.IsNull() || password.Value == "" {
			response.Diagnostics.AddAttributeError(
				request.AttributePath,
				validator.Description(ctx),
				fmt.Sprintf("%v must be defined and not empty for '%v' integration",
					attributePassword,
					snyk.ACRIntegrationType,
				),
			)
		}

		if registryURL.IsNull() || registryURL.Value == "" {
			response.Diagnostics.AddAttributeError(
				request.AttributePath,
				validator.Description(ctx),
				fmt.Sprintf("%v must be defined and not empty for '%v' integration",
					attributeRegistryURL,
					snyk.ACRIntegrationType,
				),
			)
		}

	case snyk.BitBucketCloudIntegrationType:
		username, diags := getAttributeValue(ctx, attributeUsername, request)
		response.Diagnostics.Append(diags...)
		if response.Diagnostics.HasError() {
			return
		}

		password, diags := getAttributeValue(ctx, attributePassword, request)
		response.Diagnostics.Append(diags...)
		if response.Diagnostics.HasError() {
			return
		}

		if username.IsNull() || username.Value == "" {
			response.Diagnostics.AddAttributeError(
				request.AttributePath,
				validator.Description(ctx),
				fmt.Sprintf("%v must be defined and not empty for '%v' integration",
					attributeUsername,
					snyk.BitBucketCloudIntegrationType,
				),
			)
		}

		if password.IsNull() || password.Value == "" {
			response.Diagnostics.AddAttributeError(
				request.AttributePath,
				validator.Description(ctx),
				fmt.Sprintf("%v must be defined and not empty for '%v' integration",
					attributePassword,
					snyk.BitBucketCloudIntegrationType,
				),
			)
		}

	case snyk.BitBucketServerIntegrationType:
		username, diags := getAttributeValue(ctx, attributeUsername, request)
		response.Diagnostics.Append(diags...)
		if response.Diagnostics.HasError() {
			return
		}

		password, diags := getAttributeValue(ctx, attributePassword, request)
		response.Diagnostics.Append(diags...)
		if response.Diagnostics.HasError() {
			return
		}

		url, diags := getAttributeValue(ctx, attributeURL, request)
		response.Diagnostics.Append(diags...)
		if response.Diagnostics.HasError() {
			return
		}

		if username.IsNull() || username.Value == "" {
			response.Diagnostics.AddAttributeError(
				request.AttributePath,
				validator.Description(ctx),
				fmt.Sprintf("%v must be defined and not empty for '%v' integration",
					attributeUsername,
					snyk.BitBucketServerIntegrationType,
				),
			)
		}

		if password.IsNull() || password.Value == "" {
			response.Diagnostics.AddAttributeError(
				request.AttributePath,
				validator.Description(ctx),
				fmt.Sprintf("%v must be defined and not empty for '%v' integration",
					attributePassword,
					snyk.BitBucketServerIntegrationType,
				),
			)
		}

		if url.IsNull() || url.Value == "" {
			response.Diagnostics.AddAttributeError(
				request.AttributePath,
				validator.Description(ctx),
				fmt.Sprintf("%v must be defined and not empty for '%v' integration",
					attributeURL,
					snyk.BitBucketServerIntegrationType,
				),
			)
		}

	case snyk.GitHubIntegrationType:
		token, diags := getAttributeValue(ctx, attributeToken, request)
		response.Diagnostics.Append(diags...)
		if response.Diagnostics.HasError() {
			return
		}

		if token.IsNull() || token.Value == "" {
			response.Diagnostics.AddAttributeError(
				request.AttributePath,
				validator.Description(ctx),
				fmt.Sprintf("%v must be defined and not empty for '%v' integration",
					attributeToken,
					snyk.GitHubIntegrationType,
				),
			)
		}

	case snyk.GitLabIntegrationType:
		token, diags := getAttributeValue(ctx, attributeToken, request)
		response.Diagnostics.Append(diags...)
		if response.Diagnostics.HasError() {
			return
		}

		if token.IsNull() || token.Value == "" {
			response.Diagnostics.AddAttributeError(
				request.AttributePath,
				validator.Description(ctx),
				fmt.Sprintf("%v must be defined and not empty for '%v' integration",
					attributeToken,
					snyk.GitLabIntegrationType,
				),
			)
		}
	}
}

func getAttributeValue(ctx context.Context, attrName string, request tfsdk.ValidateAttributeRequest) (types.String, diag.Diagnostics) {
	attrPaths, diags := request.Config.PathMatches(ctx, path.MatchRoot(attrName))
	if diags.HasError() {
		return types.String{}, diags
	}

	var value types.String
	diags = request.Config.GetAttribute(ctx, attrPaths[0], &value)
	if diags.HasError() {
		return types.String{}, diags
	}

	return value, nil
}

func AllowedIntegrationTypes() []string {
	return []string{
		snyk.ACRIntegrationType,
		snyk.ArtifactoryCRIntegrationType,
		snyk.AzureReposIntegrationType,
		snyk.BitBucketCloudIntegrationType,
		snyk.BitBucketConnectAppIntegrationType,
		snyk.BitBucketServerIntegrationType,
		snyk.DigitalOceanCRIntegrationType,
		snyk.DockerHubIntegrationType,
		snyk.ECRIntegrationType,
		snyk.GCRIntegrationType,
		snyk.GitHubIntegrationType,
		snyk.GitHubCRIntegrationType,
		snyk.GitHubEnterpriseIntegrationType,
		snyk.GitLabIntegrationType,
		snyk.GitLabCRIntegrationType,
		snyk.GoogleArtifactCRIntegrationType,
		snyk.HarborCRIntegrationType,
		snyk.NexusCRIntegrationType,
		snyk.QuayCRIntegrationType,
	}
}

// RequiresConfiguredCredentials checks that a set of path.Expression match
// specific integration type, e.g.
//
//   - type 'github' has token defined
//   - type 'acr' has username, password and registryBase defined
//
// Full matrix can be found under attributes, see https://snyk.docs.apiary.io/#reference/integrations/integrations/add-new-integration
func RequiresConfiguredCredentials() tfsdk.AttributeValidator {
	return &requiredConfiguredCredentialsValidator{}
}
