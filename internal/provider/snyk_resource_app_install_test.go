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

const (
	universalBrokerAppID = "921020b6-b167-426e-867b-3e2856a2f56e"
)

func init() {
	resource.AddTestSweepers("snyk_app_install", &resource.Sweeper{
		Name: "snyk_app_install",
		F: func(region string) error {
			ctx := context.Background()
			client, err := sharedClient(region)
			if err != nil {
				return fmt.Errorf("could not create shared Snyk SDK client: %w", err)
			}

			// iterate over orgs created by acceptance tests and remove any UB Snyk Apps
			orgs, errf := client.Orgs.AllAccessibleOrgs(ctx, nil)
			for org := range orgs {
				if strings.HasPrefix(org.Attributes.Name, accTestPrefix) {
					appInstalls, resp, err := client.Apps.ListAppInstallsForOrg(ctx, org.ID, nil)
					if err != nil {
						slog.Warn(
							"Error listing app installs for organization",
							"org_name", org.Attributes.Name,
							"error", err,
							"snyk_request_id", resp.SnykRequestID,
						)
					}

					for _, appInstall := range appInstalls {
						if appInstall.Relationships.App.Data.ID == universalBrokerAppID {
							slog.Info(
								"Deleting snyk_app_install for organization",
								"app_install_id", appInstall.ID,
								"org_name", org.Attributes.Name,
								"org_id", org.ID,
							)
							resp, err := client.Apps.DeleteAppInstallFromOrg(ctx, org.ID, appInstall.ID)
							if err != nil {
								slog.Warn(
									"Error deleting app install during sweep",
									"app_install_id", appInstall.ID,
									"org_name", org.Attributes.Name,
									"org_id", org.ID,
									"error", err,
									"snyk_request_id", resp.SnykRequestID,
								)
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

func TestAccSnykAppInstallResource(t *testing.T) {
	t.Parallel()

	orgName := acctest.RandomWithPrefix(accTestPrefix)
	groupID := accTestGroupID()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// Create and Read testing
				Config: testAccSnykAppInstallResourceConfig(orgName, groupID, universalBrokerAppID),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"snyk_app_install.test",
						tfjsonpath.New("app_id"),
						knownvalue.StringExact(universalBrokerAppID),
					),
					statecheck.ExpectKnownValue(
						"snyk_app_install.test",
						tfjsonpath.New("app_name"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"snyk_app_install.test",
						tfjsonpath.New("client_id"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectSensitiveValue(
						"snyk_app_install.test",
						tfjsonpath.New("client_secret"),
					),
					statecheck.ExpectKnownValue(
						"snyk_app_install.test",
						tfjsonpath.New("id"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"snyk_app_install.test",
						tfjsonpath.New("organization_id"),
						knownvalue.NotNull(),
					),
				},
			},
		},
	})
}

func testAccSnykAppInstallResourceConfig(orgName, groupID, appID string) string {
	return fmt.Sprintf(`
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
