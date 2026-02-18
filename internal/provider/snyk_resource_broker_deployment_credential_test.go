package provider

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func init() {
	resource.AddTestSweepers("snyk_broker_deployment_credential", &resource.Sweeper{
		Name: "snyk_broker_deployment_credential",
		F: func(region string) error {
			ctx := context.Background()
			client, err := sharedClient(region)
			if err != nil {
				return fmt.Errorf("could not create shared Snyk SDK client: %w", err)
			}

			// iterate over orgs created by acceptance tests and remove deployments with metadata.env accTestPrefix
			orgs, errf := client.Orgs.AllAccessibleOrgs(ctx, nil)
			for org := range orgs {
				if strings.HasPrefix(org.Attributes.Name, accTestPrefix) {
					appInstalls, resp, err := client.Apps.ListAppInstallsForOrg(ctx, org.ID, nil)
					if err != nil {
						slog.Warn(
							"Error listing app installs for organization",
							"organization_name", org.Attributes.Name,
							"error", err,
							"snyk_request_id", resp.SnykRequestID,
						)
					}

					for _, appInstall := range appInstalls {
						if appInstall.Relationships.App.Data.ID == universalBrokerAppID {
							// get organization with enriched properties (we need tenantID)
							org, resp, err := client.Orgs.Get(ctx, org.ID, nil)
							if err != nil {
								slog.Warn(
									"Error getting organization with enriched properties",
									"organization_id", org.ID,
									"organization_name", org.Attributes.Name,
									"error", err,
									"snyk_request_id", resp.SnykRequestID,
								)
								continue
							}

							tenantID := org.Relationships.Tenant.Data.ID
							brokerDeployments, resp, err := client.Brokers.ListDeployments(ctx, tenantID, appInstall.ID)
							if err != nil {
								slog.Warn(
									"Error listing broker deployments",
									"app_install_id", appInstall.ID,
									"error", err,
									"snyk_request_id", resp.SnykRequestID,
									"tenant_id", tenantID,
								)
								continue
							}

							for _, brokerDeployment := range brokerDeployments {
								metadata := brokerDeployment.Attributes.Metadata
								if len(metadata) == 0 {
									continue
								}

								if strings.HasPrefix(metadata["env"], accTestPrefix) {
									brokerDeploymentCredentials, resp, err := client.Brokers.ListDeploymentCredentials(ctx, tenantID, appInstall.ID, brokerDeployment.ID)
									if err != nil {
										slog.Warn(
											"Error listing broker deployment credentials",
											"app_install_id", appInstall.ID,
											"broker_deployment_id", brokerDeployment.ID,
											"error", err,
											"snyk_request_id", resp.SnykRequestID,
											"tenant_id", tenantID,
										)
										continue
									}

									for _, brokerDeploymentCredential := range brokerDeploymentCredentials {
										slog.Info(
											"Deleting snyk_broker_deployment credentials",
											"app_install_id", appInstall.ID,
											"broker_deployment_credential_id", brokerDeploymentCredential.ID,
											"broker_deployment_id", brokerDeployment.ID,
											"tenant_id", tenantID,
										)
										resp, err := client.Brokers.DeleteDeploymentCredential(ctx, tenantID, appInstall.ID, brokerDeployment.ID, brokerDeploymentCredential.ID)
										if err != nil {
											slog.Warn(
												"Error deleting broker deployment credential during sweep",
												"broker_deployment_credential_id", brokerDeploymentCredential.ID,
												"error", err,
												"snyk_request_id", resp.SnykRequestID,
											)
										}
									}
								}
							}
						}
					}
				}
			}
			if err := errf(); err != nil {
				return fmt.Errorf("unable to iterate over all accessible orgs: %w", err)
			}

			return nil
		},
	})
}

func TestAccSnykBrokerDeploymentCredentialResource(t *testing.T) {
	orgName := acctest.RandomWithPrefix(accTestPrefix)
	groupID := accTestGroupID()
	envVarName := acctest.RandString(5)
	envVarNameUpdated := acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccSnykBrokerDeploymentCredentialResourceConfig(orgName, groupID, universalBrokerAppID, envVarName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"snyk_broker_deployment_credential.test",
						tfjsonpath.New("app_install_id"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"snyk_broker_deployment_credential.test",
						tfjsonpath.New("broker_connection_type"),
						knownvalue.StringExact("gitlab"),
					),
					statecheck.ExpectKnownValue(
						"snyk_broker_deployment_credential.test",
						tfjsonpath.New("broker_deployment_id"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"snyk_broker_deployment_credential.test",
						tfjsonpath.New("environment_variable_name"),
						knownvalue.StringExact(envVarName),
					),
					statecheck.ExpectKnownValue(
						"snyk_broker_deployment_credential.test",
						tfjsonpath.New("id"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"snyk_broker_deployment_credential.test",
						tfjsonpath.New("tenant_id"),
						knownvalue.NotNull(),
					),
				},
			},
			// Update and Read testing
			{
				Config: testAccSnykBrokerDeploymentCredentialResourceConfig(orgName, groupID, universalBrokerAppID, envVarNameUpdated),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"snyk_broker_deployment_credential.test",
						tfjsonpath.New("app_install_id"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"snyk_broker_deployment_credential.test",
						tfjsonpath.New("broker_connection_type"),
						knownvalue.StringExact("gitlab"),
					),
					statecheck.ExpectKnownValue(
						"snyk_broker_deployment_credential.test",
						tfjsonpath.New("broker_deployment_id"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"snyk_broker_deployment_credential.test",
						tfjsonpath.New("environment_variable_name"),
						knownvalue.StringExact(envVarNameUpdated),
					),
					statecheck.ExpectKnownValue(
						"snyk_broker_deployment_credential.test",
						tfjsonpath.New("id"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"snyk_broker_deployment_credential.test",
						tfjsonpath.New("tenant_id"),
						knownvalue.NotNull(),
					),
				},
			},
		},
	})
}

func testAccSnykBrokerDeploymentCredentialResourceConfig(orgName, groupID, appID, envVarName string) string {
	return fmt.Sprintf(`
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
`, orgName, groupID, appID, envVarName)
}
