package provider

//import (
//	"os"
//	"regexp"
//	"testing"
//
//	"github.com/hashicorp/terraform-plugin-framework/providerserver"
//	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
//	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
//	"github.com/pavel-snyk/snyk-sdk-go/snyk"
//	"github.com/stretchr/testify/assert"
//)
//
//var testAccProvider = New("testacc")()
//var testAccProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
//	"snyk": func() (tfprotov6.ProviderServer, error) {
//		return providerserver.NewProtocol6WithError(testAccProvider)()
//	},
//}
//var snykClient *snyk.Client
//
//func testSnykClient() *snyk.Client {
//	if snykClient == nil {
//		endpoint := os.Getenv("SNYK_ENDPOINT")
//		token := os.Getenv("SNYK_TOKEN")
//
//		snykClient = snyk.NewClient(token,
//			snyk.WithBaseURL(endpoint),
//			snyk.WithUserAgent(testAccProvider.(*snykProvider).userAgent()),
//		)
//	}
//	return snykClient
//}
//
//func testAccPreCheck(t *testing.T) {
//	if v := os.Getenv("SNYK_ENDPOINT"); v == "" {
//		t.Fatal("SNYK_ENDPOINT must be set to run acceptance tests.")
//	}
//
//	if v := os.Getenv("SNYK_TOKEN"); v == "" {
//		t.Fatal("SNYK_TOKEN must be set to run acceptance tests.")
//	}
//
//	if v := os.Getenv("SNYK_GROUP_ID"); v == "" {
//		t.Fatal("SNYK_GROUP_ID must be set to run acceptance tests.")
//	}
//}
//
//func TestProvider_MissingTokenAttribute(t *testing.T) {
//	resource.UnitTest(t, resource.TestCase{
//		ProtoV6ProviderFactories: testAccProviderFactories,
//
//		Steps: []resource.TestStep{
//			{
//				Config: `
//provider "snyk" {
//  token = ""
//}
//data "snyk_user" "self" {}
//				`,
//				ExpectError: regexp.MustCompile("token must be set"),
//			},
//		},
//	})
//}
//
//func TestProvider_UserAgent(t *testing.T) {
//	t.Parallel()
//
//	type testCase struct {
//		version  string
//		expected string
//	}
//	tests := map[string]testCase{
//		"empty_version": {
//			version:  "",
//			expected: "terraform-provider-snyk/ (+https://registry.terraform.io/providers/pavel-snyk/snyk)",
//		},
//		"dev_version": {
//			version:  "dev",
//			expected: "terraform-provider-snyk/dev (+https://registry.terraform.io/providers/pavel-snyk/snyk)",
//		},
//		"release_version": {
//			version:  "1.1.1",
//			expected: "terraform-provider-snyk/1.1.1 (+https://registry.terraform.io/providers/pavel-snyk/snyk)",
//		},
//	}
//
//	for name, test := range tests {
//		t.Run(name, func(t *testing.T) {
//			p := &snykProvider{version: test.version}
//
//			assert.Equal(t, test.expected, p.userAgent())
//		})
//	}
//}
