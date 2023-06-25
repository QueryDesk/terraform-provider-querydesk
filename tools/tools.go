// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

//go:build tools

package tools

import (
	// Documentation generation
	_ "github.com/Khan/genqlient/generate"
	// graphql generation
	_ "github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs"
	_ "github.com/suessflorian/gqlfetch"
)
