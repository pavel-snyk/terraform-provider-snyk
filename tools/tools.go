//go:build tools
// +build tools

package tools

import (
	// changelog generation (git-chglog)
	_ "github.com/git-chglog/git-chglog/cmd/git-chglog"
	// documentation generation and validation (tfplugindocs)
	_ "github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs"
)
