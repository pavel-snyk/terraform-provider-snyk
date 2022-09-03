package provider

import (
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var testAccProvider = New()
var testAccProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"snyk": func() (tfprotov6.ProviderServer, error) {
		return providerserver.NewProtocol6WithError(testAccProvider)()
	},
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("SNYK_ENDPOINT"); v == "" {
		t.Fatal("SNYK_ENDPOINT must be set to run acceptance tests.")
	}

	if v := os.Getenv("SNYK_TOKEN"); v == "" {
		t.Fatal("SNYK_TOKEN must be set to run acceptance tests.")
	}

	if v := os.Getenv("SNYK_GROUP_ID"); v == "" {
		t.Fatal("SNYK_GROUP_ID must be set to run acceptance tests.")
	}
}

func TestProvider_MissingTokenAttribute(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProviderFactories,

		Steps: []resource.TestStep{
			{
				Config: `
provider "snyk" {
  token = ""
}
data "snyk_user" "self" {}
				`,
				ExpectError: regexp.MustCompile("token must be set"),
			},
		},
	})
}
