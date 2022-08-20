//go:build tools
// +build tools

package tools

import (
	// Documentation generation and validation (tfplugindocs)
	_ "github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs"
)
