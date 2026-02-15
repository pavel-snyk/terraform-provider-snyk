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
	resource.AddTestSweepers("snyk_broker_deployment", &resource.Sweeper{
		Name: "snyk_broker_deployment",
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
							}

							for _, brokerDeployment := range brokerDeployments {
								metadata := brokerDeployment.Attributes.Metadata
								if len(metadata) == 0 {
									continue
								}

								if strings.HasPrefix(metadata["env"], accTestPrefix) {
									slog.Info(
										"Deleting snyk_broker_deployment",
										"app_install_id", appInstall.ID,
										"broker_deployment_id", brokerDeployment.ID,
										"tenant_id", tenantID,
									)
									resp, err := client.Brokers.DeleteDeployment(ctx, tenantID, appInstall.ID, brokerDeployment.ID)
									if err != nil {
										slog.Warn(
											"Error deleting broker deployment during sweep",
											"broker_deployment_id", brokerDeployment.ID,
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
			if err := errf(); err != nil {
				return fmt.Errorf("unable to iterate over all accessible orgs: %w", err)
			}

			return nil
		},
	})
}

func TestAccSnykBrokerDeploymentResource(t *testing.T) {
	t.Parallel()

	orgName := acctest.RandomWithPrefix(accTestPrefix)
	groupID := accTestGroupID()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccSnykBrokerDeploymentResourceConfig(orgName, groupID, universalBrokerAppID),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"snyk_broker_deployment.test",
						tfjsonpath.New("app_install_id"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"snyk_broker_deployment.test",
						tfjsonpath.New("id"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"snyk_broker_deployment.test",
						tfjsonpath.New("metadata"),
						knownvalue.MapExact(map[string]knownvalue.Check{}),
					),
					statecheck.ExpectKnownValue(
						"snyk_broker_deployment.test",
						tfjsonpath.New("organization_id"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"snyk_broker_deployment.test",
						tfjsonpath.New("tenant_id"),
						knownvalue.NotNull(),
					),
				},
			},
			// Update and Read testing with adding "metadata"
			{
				Config: testAccSnykBrokerDeploymentResourceConfigWithMetadata(orgName, groupID, universalBrokerAppID),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"snyk_broker_deployment.test",
						tfjsonpath.New("app_install_id"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"snyk_broker_deployment.test",
						tfjsonpath.New("id"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"snyk_broker_deployment.test",
						tfjsonpath.New("metadata"),
						knownvalue.MapSizeExact(2),
					),
					statecheck.ExpectKnownValue(
						"snyk_broker_deployment.test",
						tfjsonpath.New("metadata").AtMapKey("env"),
						knownvalue.StringExact(orgName),
					),
					statecheck.ExpectKnownValue(
						"snyk_broker_deployment.test",
						tfjsonpath.New("metadata").AtMapKey("comment"),
						knownvalue.StringExact("created by terraform acceptance tests"),
					),
					statecheck.ExpectKnownValue(
						"snyk_broker_deployment.test",
						tfjsonpath.New("organization_id"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"snyk_broker_deployment.test",
						tfjsonpath.New("tenant_id"),
						knownvalue.NotNull(),
					),
				},
			},
			// Update and Read testing with removing "metadata"
			{
				Config: testAccSnykBrokerDeploymentResourceConfig(orgName, groupID, universalBrokerAppID),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"snyk_broker_deployment.test",
						tfjsonpath.New("app_install_id"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"snyk_broker_deployment.test",
						tfjsonpath.New("id"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"snyk_broker_deployment.test",
						tfjsonpath.New("metadata"),
						knownvalue.MapExact(map[string]knownvalue.Check{}),
					),
					statecheck.ExpectKnownValue(
						"snyk_broker_deployment.test",
						tfjsonpath.New("organization_id"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"snyk_broker_deployment.test",
						tfjsonpath.New("tenant_id"),
						knownvalue.NotNull(),
					),
				},
			},
		},
	})
}

func testAccSnykBrokerDeploymentResourceConfig(orgName, groupID, appID string) string {
	return fmt.Sprintf(`
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
`, orgName, groupID, appID)
}

func testAccSnykBrokerDeploymentResourceConfigWithMetadata(orgName, groupID, appID string) string {
	return fmt.Sprintf(`
resource "snyk_broker_deployment" "test" {
  app_install_id  = snyk_app_install.test.id
  organization_id = snyk_organization.test.id
  tenant_id       = snyk_organization.test.tenant_id

  metadata = {
    env     = %[1]q
    comment = "created by terraform acceptance tests"
  }
}

resource "snyk_app_install" "test" {
  app_id          = %[3]q
  organization_id = snyk_organization.test.id
}

resource "snyk_organization" "test" {
  name     = %[1]q
  group_id = %[2]q
}
`, orgName, groupID, appID)
}
