package provider

import (
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/assert"
)

const accTestPrefix = "tf-acc-test"

var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"snyk": providerserver.NewProtocol6WithError(New("testacc")()),
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("SNYK_APP_BASE_URL"); v == "" {
		t.Fatal("SNYK_APP_BASE_URL must be set to run acceptance tests.")
	}
	if v := os.Getenv("SNYK_GROUP_ID"); v == "" {
		t.Fatal("SNYK_GROUP_ID must be set to run acceptance tests.")
	}
	if v := os.Getenv("SNYK_REGION"); v == "" {
		t.Fatal("SNYK_REGION must be set to run acceptance tests.")
	}
	if v := os.Getenv("SNYK_REST_BASE_URL"); v == "" {
		t.Fatal("SNYK_REST_BASE_URL must be set to run acceptance tests.")
	}
	if v := os.Getenv("SNYK_TOKEN"); v == "" {
		t.Fatal("SNYK_TOKEN must be set to run acceptance tests.")
	}
	if v := os.Getenv("SNYK_V1_BASE_URL"); v == "" {
		t.Fatal("SNYK_V1_BASE_URL must be set to run acceptance tests.")
	}
}

func accTestGroupID() string {
	return os.Getenv("SNYK_GROUP_ID")
}

func TestProvider_MissingTokenAttribute(t *testing.T) {
	t.Parallel()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{{
			Config: `
provider "snyk" {
  token = ""
}
data "snyk_organization" "self" {}
`,
			ExpectError: regexp.MustCompile("The Snyk API token must be provided"),
		}},
	})
}

func TestProvider_UserAgent(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		version           string
		expectedUserAgent string
	}{
		"empty-version": {
			version:           "",
			expectedUserAgent: "terraform-provider-snyk/ (+https://registry.terraform.io/providers/pavel-snyk/snyk)",
		},
		"dev-version": {
			version:           "dev",
			expectedUserAgent: "terraform-provider-snyk/dev (+https://registry.terraform.io/providers/pavel-snyk/snyk)",
		},
		"release-version": {
			version:           "1.1.1",
			expectedUserAgent: "terraform-provider-snyk/1.1.1 (+https://registry.terraform.io/providers/pavel-snyk/snyk)",
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			p := &snykProvider{version: test.version}
			actualUserAgent := p.userAgent()

			assert.Equal(t, test.expectedUserAgent, actualUserAgent)
		})
	}
}
