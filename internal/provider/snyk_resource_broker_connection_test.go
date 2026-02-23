package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func init() {}

func TestAccSnykBrokerConnectionResource(t *testing.T) {
	orgName := acctest.RandomWithPrefix(accTestPrefix)
	groupID := accTestGroupID()
	envVarName := acctest.RandString(5)
	connectionName := acctest.RandomWithPrefix(accTestPrefix)
	connectionNameUpdated := acctest.RandomWithPrefix(accTestPrefix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccSnykBrokerConnectionResourceConfig(orgName, groupID, universalBrokerAppID, envVarName, connectionName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"snyk_broker_connection.test",
						tfjsonpath.New("app_install_id"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"snyk_broker_connection.test",
						tfjsonpath.New("broker_deployment_id"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"snyk_broker_connection.test",
						tfjsonpath.New("id"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"snyk_broker_connection.test",
						tfjsonpath.New("configuration"),
						knownvalue.ObjectPartial(map[string]knownvalue.Check{
							"broker_client_url":          knownvalue.StringExact("https://api.snyk.io"),
							"gitlab_hostname":            knownvalue.StringExact("gitlab.com"),
							"gitlab_token_credential_id": knownvalue.NotNull(),
						}),
					),
					statecheck.ExpectKnownValue(
						"snyk_broker_connection.test",
						tfjsonpath.New("name"),
						knownvalue.StringExact(connectionName),
					),
					statecheck.ExpectKnownValue(
						"snyk_broker_connection.test",
						tfjsonpath.New("tenant_id"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"snyk_broker_connection.test",
						tfjsonpath.New("type"),
						knownvalue.StringExact("gitlab"),
					),
				},
			},
			// Update and Read testing with updated "name"
			{
				Config: testAccSnykBrokerConnectionResourceConfig(orgName, groupID, universalBrokerAppID, envVarName, connectionNameUpdated),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"snyk_broker_connection.test",
						tfjsonpath.New("app_install_id"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"snyk_broker_connection.test",
						tfjsonpath.New("broker_deployment_id"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"snyk_broker_connection.test",
						tfjsonpath.New("id"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"snyk_broker_connection.test",
						tfjsonpath.New("configuration"),
						knownvalue.ObjectPartial(map[string]knownvalue.Check{
							"broker_client_url":          knownvalue.StringExact("https://api.snyk.io"),
							"gitlab_hostname":            knownvalue.StringExact("gitlab.com"),
							"gitlab_token_credential_id": knownvalue.NotNull(),
						}),
					),
					statecheck.ExpectKnownValue(
						"snyk_broker_connection.test",
						tfjsonpath.New("name"),
						knownvalue.StringExact(connectionNameUpdated),
					),
					statecheck.ExpectKnownValue(
						"snyk_broker_connection.test",
						tfjsonpath.New("tenant_id"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"snyk_broker_connection.test",
						tfjsonpath.New("type"),
						knownvalue.StringExact("gitlab"),
					),
				},
			},
		},
	})
}

func testAccSnykBrokerConnectionResourceConfig(orgName, groupID, appID, envVarName, connectionName string) string {
	return fmt.Sprintf(`
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
