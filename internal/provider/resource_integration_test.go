package provider

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/pavel-snyk/snyk-sdk-go/snyk"
)

func TestAccResourceIntegration_basic(t *testing.T) {
	var integration snyk.Integration
	organizationName := fmt.Sprintf("tf-test-acc_%s", acctest.RandString(10))
	groupID := os.Getenv("SNYK_GROUP_ID")
	token := acctest.RandString(20)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccResourceIntegrationConfig(organizationName, groupID, ""),
				ExpectError: regexp.MustCompile("Wrong credentials for given integration type"),
			},
			{
				Config: testAccResourceIntegrationConfig(organizationName, groupID, token),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckResourceIntegrationExists("snyk_integration.test", organizationName, &integration),
					resource.TestCheckResourceAttrSet("snyk_integration.test", "id"),
					resource.TestCheckResourceAttr("snyk_integration.test", "type", "gitlab"),
					resource.TestCheckNoResourceAttr("snyk_integration.test", "pull_request_sca.enabled"),
				),
			},
		},
	})
}

func TestAccResourceIntegration_pullRequestSCA(t *testing.T) {
	var integration snyk.Integration
	organizationName := fmt.Sprintf("tf-test-acc_%s", acctest.RandString(10))
	groupID := os.Getenv("SNYK_GROUP_ID")
	token := acctest.RandString(20)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceIntegrationConfigWithPullRequestSCA(organizationName, groupID, token),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckResourceIntegrationExists("snyk_integration.test", organizationName, &integration),
					resource.TestCheckResourceAttrSet("snyk_integration.test", "pull_request_sca.enabled"),
					resource.TestCheckResourceAttr("snyk_integration.test", "pull_request_sca.fail_on_any_issue", "true"),
					resource.TestCheckResourceAttr("snyk_integration.test", "pull_request_sca.fail_only_on_issues_with_fix", "false"),
				),
			},
		},
	})
}

func testAccCheckResourceIntegrationExists(resourceName, organizationName string, integration *snyk.Integration) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		// retrieve resource from state
		rs := state.RootModule().Resources[resourceName]

		if rs.Primary.ID == "" {
			return fmt.Errorf("integration ID is not set")
		}

		client := testSnykClient()
		orgs, _, err := client.Orgs.List(context.Background())
		if err != nil {
			return err
		}

		organizationID := ""
		for _, org := range orgs {
			if org.Name == organizationName {
				organizationID = org.ID
				break
			}
		}
		if organizationID == "" {
			return fmt.Errorf("organization (%s) for integration (%s) not found", organizationName, rs.Primary.ID)
		}

		integrations, _, err := client.Integrations.List(context.Background(), organizationID)
		if err != nil {
			return err
		}

		for t, id := range integrations {
			if id == rs.Primary.ID {
				integration = &snyk.Integration{
					ID:   id,
					Type: t,
				}
				return nil
			}
		}

		return fmt.Errorf("integration (%s) not found", rs.Primary.ID)
	}
}

func testAccResourceIntegrationConfig(organizationName, groupID, token string) string {
	return fmt.Sprintf(`
resource "snyk_organization" "test" {
  name = "%s"
  group_id = "%s"
}
resource "snyk_integration" "test" {
  organization_id = snyk_organization.test.id

  type  = "gitlab"
  url   = "https://testing.gitlab.local"
  token = "%s"
}
`, organizationName, groupID, token)
}

func testAccResourceIntegrationConfigWithPullRequestSCA(organizationName, groupID, token string) string {
	return fmt.Sprintf(`
resource "snyk_organization" "test" {
  name = "%s"
  group_id = "%s"
}
resource "snyk_integration" "test" {
  organization_id = snyk_organization.test.id

  type  = "gitlab"
  url   = "https://testing.gitlab.local"
  token = "%s"

  pull_request_sca = {
    enabled = true

    fail_on_any_issue            = true
    fail_only_on_issues_with_fix = false
  }
}
`, organizationName, groupID, token)
}
