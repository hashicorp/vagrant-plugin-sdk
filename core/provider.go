// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package core

type SshInfo struct {
	Host           *string
	Port           *string
	Username       *string
	PrivateKeyPath *string
}

type Provider interface {
	CapabilityPlatform

	Usable() (bool, error)
	Installed() (bool, error)
	Action(name string, args ...interface{}) error
	MachineIdChanged() error
	SshInfo() (*SshInfo, error)
	State() (*MachineState, error)
}
