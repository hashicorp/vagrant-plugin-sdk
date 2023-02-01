// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package core

import (
	"io"
)

type Host interface {
	CapabilityPlatform
	Seeder
	Named

	Detect(state StateBag) (bool, error)
	Parent() (string, error)

	io.Closer
}
