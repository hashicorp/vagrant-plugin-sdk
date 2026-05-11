// Copyright IBM Corp. 2020, 2026
// SPDX-License-Identifier: MPL-2.0

package core

type CapabilityPlatform interface {
	Capability(name string, args ...interface{}) (interface{}, error)
	HasCapability(name string) (bool, error)
}
