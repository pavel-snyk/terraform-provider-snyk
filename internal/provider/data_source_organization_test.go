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

func TestAccSnykOrganizationDataSource(t *testing.T) {
	t.Parallel()

	name := acctest.RandomWithPrefix(accTestPrefix)
	groupID := accTestGroupID()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccSnykOrganizationDataSourceConfig(name, groupID),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"data.snyk_organization.test",
						tfjsonpath.New("group_id"),
						knownvalue.StringExact(groupID),
					),
					statecheck.ExpectKnownValue(
						"data.snyk_organization.test",
						tfjsonpath.New("id"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"data.snyk_organization.test",
						tfjsonpath.New("name"),
						knownvalue.StringExact(name),
					),
					statecheck.ExpectKnownValue(
						"data.snyk_organization.test",
						tfjsonpath.New("slug"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"data.snyk_organization.test",
						tfjsonpath.New("tenant_id"),
						knownvalue.NotNull(),
					),
				},
			},
		},
	})
}

func TestAccSnykOrganizationDataSource_expectError(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccSnykOrganizationDataSourceConfigWithoutIDAndName,
				ExpectError: regexp.MustCompile(`The attribute "id" or "name" must be defined`),
			},
		},
	})
}

func testAccSnykOrganizationDataSourceConfig(name, groupID string) string {
	return fmt.Sprintf(`
data "snyk_organization" "test" {
  id = snyk_organization.test.id
}

resource "snyk_organization" "test" {
  name     = %[1]q
  group_id = %[2]q
}
`, name, groupID)
}

const testAccSnykOrganizationDataSourceConfigWithoutIDAndName = `
data "snyk_organization" "test" {}
`
