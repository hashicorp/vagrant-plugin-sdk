// Copyright IBM Corp. 2020, 2026
// SPDX-License-Identifier: MPL-2.0

package core

import (
	"io"

	"github.com/hashicorp/vagrant-plugin-sdk/component"
)

type Command interface {
	CommandInfo() (*component.CommandInfo, error)
	Execute([]string) (int32, error)

	io.Closer
}
