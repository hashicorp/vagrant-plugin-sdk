// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package core

type CapabilityPlatform interface {
	Capability(name string, args ...interface{}) (interface{}, error)
	HasCapability(name string) (bool, error)
}
