package provider

//import (
//	"context"
//	"fmt"
//	"os"
//	"regexp"
//	"testing"
//
//	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
//	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
//	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
//
//	"github.com/pavel-snyk/snyk-sdk-go/snyk"
//)
//
//func TestAccResourceOrganization_basic(t *testing.T) {
//	var organization snyk.Organization
//	organizationName := fmt.Sprintf("tf-test-acc_%s", acctest.RandString(10))
//	groupID := os.Getenv("SNYK_GROUP_ID")
//
//	resource.Test(t, resource.TestCase{
//		PreCheck:                 func() { testAccPreCheck(t) },
//		ProtoV6ProviderFactories: testAccProviderFactories,
//		CheckDestroy:             testAccCheckResourceOrganizationDestroy,
//		Steps: []resource.TestStep{
//			{
//				Config:      testAccResourceOrganizationConfig("", groupID),
//				ExpectError: regexp.MustCompile("string must not be empty"),
//			},
//			{
//				Config: testAccResourceOrganizationConfig(organizationName, groupID),
//				Check: resource.ComposeAggregateTestCheckFunc(
//					testAccCheckResourceOrganizationExists("snyk_organization.test", &organization),
//					resource.TestCheckResourceAttr("snyk_organization.test", "name", organizationName),
//					resource.TestCheckResourceAttr("snyk_organization.test", "group_id", groupID),
//					resource.TestCheckResourceAttrSet("snyk_organization.test", "id"),
//				),
//			},
//		},
//	})
//}
//
//func testAccCheckResourceOrganizationDestroy(state *terraform.State) error {
//	client := testSnykClient()
//
//	for _, rs := range state.RootModule().Resources {
//		if rs.Type != "snyk_organization" {
//			continue
//		}
//
//		orgs, _, err := client.Orgs.List(context.Background())
//		if err != nil {
//			return err
//		}
//
//		for _, org := range orgs {
//			if org.ID == rs.Primary.ID {
//				return fmt.Errorf("organization (%s) still exists", rs.Primary.ID)
//			}
//		}
//	}
//
//	return nil
//}
//
//func testAccCheckResourceOrganizationExists(resourceName string, organization *snyk.Organization) resource.TestCheckFunc {
//	return func(state *terraform.State) error {
//		// retrieve resource from state
//		rs := state.RootModule().Resources[resourceName]
//
//		if rs.Primary.ID == "" {
//			return fmt.Errorf("organization ID is not set")
//		}
//
//		client := testSnykClient()
//		orgs, _, err := client.Orgs.List(context.Background())
//		if err != nil {
//			return err
//		}
//
//		for _, org := range orgs {
//			if org.ID == rs.Primary.ID {
//				organization = &org
//				return nil
//			}
//		}
//
//		return fmt.Errorf("organization (%s) not found", rs.Primary.ID)
//	}
//}
//
//func testAccResourceOrganizationConfig(organizationName, groupID string) string {
//	return fmt.Sprintf(`
//resource "snyk_organization" "test" {
//  name     = "%s"
//  group_id = "%s"
//}
//`, organizationName, groupID)
//}
