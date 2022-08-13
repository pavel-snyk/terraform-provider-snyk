package main

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"

	"github.com/pavel-snyk/terraform-provider-snyk/internal/provider"
)

func main() {
	ctx := context.Background()
	_ = providerserver.Serve(ctx, provider.New, providerserver.ServeOpts{
		Address: "registry.terraform.io/pavel-snyk/snyk",
	})
}
