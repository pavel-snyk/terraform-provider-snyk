package helper

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/pavel-snyk/snyk-sdk-go/v2/snyk"
)

const (
	organizationTenantIDPopulated    = "tenantIDPopulated"
	organizationTenantIDNotPopulated = "tenantIDNotPopulated"
)

func statusOrganizationTenantIDPopulatedState(ctx context.Context, client *snyk.Client, orgID string) retry.StateRefreshFunc {
	return func() (any, string, error) {
		organization, _, err := client.Orgs.Get(ctx, orgID, nil)
		if err != nil {
			return nil, "", err
		}
		if organization == nil {
			return nil, "", fmt.Errorf("failed to get organization with id: %s", orgID)
		}

		if organization.Relationships != nil &&
			organization.Relationships.Tenant != nil &&
			organization.Relationships.Tenant.Data != nil &&
			organization.Relationships.Tenant.Data.ID != "" {
			return organization, organizationTenantIDPopulated, nil
		}

		return organization, organizationTenantIDNotPopulated, nil
	}
}
