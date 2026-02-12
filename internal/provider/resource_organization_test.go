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
	resource.AddTestSweepers("snyk_organization", &resource.Sweeper{
		Name: "snyk_organization",
		F: func(region string) error {
			ctx := context.Background()
			client, err := sharedClient(region)
			if err != nil {
				return fmt.Errorf("could not create shared Snyk SDK client: %w", err)
			}

			orgs, errf := client.Orgs.AllAccessibleOrgs(ctx, nil)
			for org := range orgs {
				if strings.HasPrefix(org.Attributes.Name, accTestPrefix) {
					slog.Info("Deleting snyk_organization", "name", org.Attributes.Name, "id", org.ID)
					resp, err := client.OrgsV1.Delete(ctx, org.ID)
					if err != nil {
						slog.Warn(
							"Error deleting organization during sweep",
							"name", org.Attributes.Name,
							"error", err,
							"snyk_request_id", resp.SnykRequestID,
						)
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

func TestAccSnykOrganizationResource(t *testing.T) {
	name := acctest.RandomWithPrefix(accTestPrefix)
	nameUpdated := acctest.RandomWithPrefix(accTestPrefix)
	groupID := accTestGroupID()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccSnykOrganizationResourceConfig(name, groupID),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"snyk_organization.test",
						tfjsonpath.New("group_id"),
						knownvalue.StringExact(groupID),
					),
					statecheck.ExpectKnownValue(
						"snyk_organization.test",
						tfjsonpath.New("id"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"snyk_organization.test",
						tfjsonpath.New("name"),
						knownvalue.StringExact(name),
					),
					statecheck.ExpectKnownValue(
						"snyk_organization.test",
						tfjsonpath.New("slug"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"snyk_organization.test",
						tfjsonpath.New("tenant_id"),
						knownvalue.NotNull(),
					),
				},
			},
			// ImportState testing
			{
				ResourceName:      "snyk_organization.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccSnykOrganizationResourceConfig(nameUpdated, groupID),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"snyk_organization.test",
						tfjsonpath.New("group_id"),
						knownvalue.StringExact(groupID),
					),
					statecheck.ExpectKnownValue(
						"snyk_organization.test",
						tfjsonpath.New("id"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"snyk_organization.test",
						tfjsonpath.New("name"),
						knownvalue.StringExact(nameUpdated),
					),
					statecheck.ExpectKnownValue(
						"snyk_organization.test",
						tfjsonpath.New("slug"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"snyk_organization.test",
						tfjsonpath.New("tenant_id"),
						knownvalue.NotNull(),
					),
				},
			},
		},
	})
}

func testAccSnykOrganizationResourceConfig(name, groupID string) string {
	return fmt.Sprintf(`
resource "snyk_organization" "test" {
  name     = %[1]q
  group_id = %[2]q
}
`, name, groupID)
}
