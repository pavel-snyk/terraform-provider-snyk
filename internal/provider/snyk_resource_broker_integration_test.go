package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func init() {}

func TestAccSnykBrokerIntegrationResource(t *testing.T) {
	orgName := acctest.RandomWithPrefix(accTestPrefix)
	groupID := accTestGroupID()
	envVarName := acctest.RandString(5)
	connectionName := acctest.RandomWithPrefix(accTestPrefix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccSnykBrokerIntegrationResourceConfig(orgName, groupID, universalBrokerAppID, envVarName, connectionName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"snyk_broker_integration.test",
						tfjsonpath.New("broker_connection_id"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"snyk_broker_integration.test",
						tfjsonpath.New("id"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"snyk_broker_integration.test",
						tfjsonpath.New("organization_id"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"snyk_broker_integration.test",
						tfjsonpath.New("tenant_id"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"snyk_broker_integration.test",
						tfjsonpath.New("type"),
						knownvalue.StringExact("gitlab"),
					),
				},
			},
		},
	})
}

func TestAccSnykBrokerIntegrationResource_expectError(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccSnykBrokerIntegrationResourceWithoutConnectionID,
				ExpectError: regexp.MustCompile(`The argument "broker_connection_id" is required`),
			},
			{
				Config:      testAccSnykBrokerIntegrationResourceWithoutOrganizationID,
				ExpectError: regexp.MustCompile(`The argument "organization_id" is required`),
			},
			{
				Config:      testAccSnykBrokerIntegrationResourceWithoutTenantID,
				ExpectError: regexp.MustCompile(`The argument "tenant_id" is required`),
			},
			{
				Config:      testAccSnykBrokerIntegrationResourceWithoutType,
				ExpectError: regexp.MustCompile(`The argument "type" is required`),
			},
		},
	})
}

func testAccSnykBrokerIntegrationResourceConfig(orgName, groupID, appID, envVarName, connectionName string) string {
	return fmt.Sprintf(`
resource "snyk_broker_integration" "test" {
  broker_connection_id = snyk_broker_connection.test.id
  organization_id      = snyk_organization.test.id
  tenant_id            = snyk_organization.test.tenant_id
  type                 = snyk_broker_connection.test.type
}

resource "snyk_broker_connection" "test" {
  app_install_id       = snyk_app_install.test.id
  tenant_id            = snyk_organization.test.tenant_id
  broker_deployment_id = snyk_broker_deployment.test.id

  type = "gitlab"
  name = %[5]q
  configuration = {
    broker_client_url          = "https://api.snyk.io"
    gitlab_hostname            = "gitlab.com"
    gitlab_token_credential_id = snyk_broker_deployment_credential.test.id
  }
}

resource "snyk_broker_deployment_credential" "test" {
  app_install_id       = snyk_app_install.test.id
  tenant_id            = snyk_organization.test.tenant_id
  broker_deployment_id = snyk_broker_deployment.test.id

  broker_connection_type    = "gitlab"
  environment_variable_name = %[4]q
}

resource "snyk_broker_deployment" "test" {
  app_install_id  = snyk_app_install.test.id
  organization_id = snyk_organization.test.id
  tenant_id       = snyk_organization.test.tenant_id
}

resource "snyk_app_install" "test" {
  app_id          = %[3]q
  organization_id = snyk_organization.test.id
}

resource "snyk_organization" "test" {
  name     = %[1]q
  group_id = %[2]q
}
`, orgName, groupID, appID, envVarName, connectionName)
}

const testAccSnykBrokerIntegrationResourceWithoutConnectionID = `
resource "snyk_broker_integration" "test" {
  organization_id      = "<organization-id>"
  tenant_id            = "<tenant-id>"
  type                 = "gitlab"
}
`

const testAccSnykBrokerIntegrationResourceWithoutOrganizationID = `
resource "snyk_broker_integration" "test" {
  broker_connection_id = "<broker-connection-id>"
  tenant_id            = "<tenant-id>"
  type                 = "gitlab"
}
`

const testAccSnykBrokerIntegrationResourceWithoutTenantID = `
resource "snyk_broker_integration" "test" {
  broker_connection_id = "<broker-connection-id>"
  organization_id      = "<organization-id>"
  type                 = "gitlab"
}
`

const testAccSnykBrokerIntegrationResourceWithoutType = `
resource "snyk_broker_integration" "test" {
  broker_connection_id = "<broker-connection-id>"
  organization_id      = "<organization-id>"
  tenant_id            = "<tenant-id>"
}
`
