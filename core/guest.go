// Copyright IBM Corp. 2020, 2026
// SPDX-License-Identifier: MPL-2.0

package core

import (
	"io"
)

type Guest interface {
	CapabilityPlatform
	Seeder
	Named

	Detect(Target) (bool, error)
	Parent() (string, error)

	io.Closer
}
