package provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/pavel-snyk/snyk-sdk-go/v2/snyk"
)

func TestMain(m *testing.M) {
	resource.TestMain(m)
}

func sharedClient(_ string) (*snyk.Client, error) {
	if v := os.Getenv("SNYK_APP_BASE_URL"); v == "" {
		return nil, fmt.Errorf("empty SNYK_APP_BASE_URL environment variable")
	}
	if v := os.Getenv("SNYK_REGION"); v == "" {
		return nil, fmt.Errorf("empty SNYK_REGION environment variable")
	}
	if v := os.Getenv("SNYK_REST_BASE_URL"); v == "" {
		return nil, fmt.Errorf("empty SNYK_REST_BASE_URL environment variable")
	}
	if v := os.Getenv("SNYK_TOKEN"); v == "" {
		return nil, fmt.Errorf("empty SNYK_TOKEN environment variable")
	}
	if v := os.Getenv("SNYK_V1_BASE_URL"); v == "" {
		return nil, fmt.Errorf("empty SNYK_V1_BASE_URL environment variable")
	}

	regionName := os.Getenv("SNYK_REGION")
	appBaseURL := os.Getenv("SNYK_APP_BASE_URL")
	restBaseURL := os.Getenv("SNYK_REST_BASE_URL")
	v1BaseURL := os.Getenv("SNYK_V1_BASE_URL")
	token := os.Getenv("SNYK_TOKEN")

	opts := []snyk.ClientOption{
		snyk.WithUserAgent("terraform-provider-snyk/sweeper"),
		snyk.WithRegion(snyk.Region{
			Alias:       regionName,
			AppBaseURL:  appBaseURL,
			RESTBaseURL: restBaseURL,
			V1BaseURL:   v1BaseURL,
		}),
	}

	return snyk.NewClient(token, opts...)
}
