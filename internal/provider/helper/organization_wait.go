package helper

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/pavel-snyk/snyk-sdk-go/v2/snyk"
)

func WaitOrganizationTenantIDPopulated(ctx context.Context, client *snyk.Client, orgID string) error {
	stateConf := &retry.StateChangeConf{
		Pending:    []string{organizationTenantIDNotPopulated},
		Target:     []string{organizationTenantIDPopulated},
		Timeout:    10 * time.Minute,
		MinTimeout: 5 * time.Second,
		Delay:      3 * time.Second,
		Refresh:    statusOrganizationTenantIDPopulatedState(ctx, client, orgID),
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("failed to wait for organization (%s) to become tenant_id set: %w", orgID, err)
	}
	return nil
}
