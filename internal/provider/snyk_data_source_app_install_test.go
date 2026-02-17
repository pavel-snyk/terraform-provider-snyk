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

func TestAccSnykAppInstallDataSource(t *testing.T) {
	orgName := acctest.RandomWithPrefix(accTestPrefix)
	groupID := accTestGroupID()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccSnykAppInstallDataSourceConfig(orgName, groupID, universalBrokerAppID),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"data.snyk_app_install.test",
						tfjsonpath.New("app_id"),
						knownvalue.StringExact(universalBrokerAppID),
					),
					statecheck.ExpectKnownValue(
						"data.snyk_app_install.test",
						tfjsonpath.New("app_name"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"data.snyk_app_install.test",
						tfjsonpath.New("client_id"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"data.snyk_app_install.test",
						tfjsonpath.New("id"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"data.snyk_app_install.test",
						tfjsonpath.New("organization_id"),
						knownvalue.NotNull(),
					),
				},
			},
		},
	})
}

func TestAccSnykAppInstallDataSource_expectError(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccSnykAppInstallDataSourceConfigWithoutOrganizationID,
				ExpectError: regexp.MustCompile(`The argument "organization_id" is required`),
			},
			{
				Config:      testAccSnykAppInstallDataSourceConfigWithoutIDAppIDAndAppName,
				ExpectError: regexp.MustCompile(`The attribute "id", "app_id" or "app_name" must be defined`),
			},
		},
	})
}

func testAccSnykAppInstallDataSourceConfig(orgName, groupID, appID string) string {
	return fmt.Sprintf(`
data "snyk_app_install" "test" {
  organization_id = snyk_organization.test.id

  id = snyk_app_install.test.id
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

const testAccSnykAppInstallDataSourceConfigWithoutOrganizationID = `
data "snyk_app_install" "test" {
  id = "<app-install-id>"
}
`

const testAccSnykAppInstallDataSourceConfigWithoutIDAppIDAndAppName = `
data "snyk_app_install" "test" {
  organization_id = "<app-install-id>"
}
`
