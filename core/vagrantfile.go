// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package core

import (
	"github.com/hashicorp/vagrant-plugin-sdk/component"
)

type Vagrantfile interface {
	GetConfig(namespace string) (*component.ConfigData, error)
	GetValue(path ...string) (interface{}, error)
	PrimaryTargetName() (name string, err error)
	Target(name, provider string) (Target, error)
	TargetConfig(name, provider string, validateProvider bool) (Vagrantfile, error)

	// Returns a list of the machines that are defined within this
	// Vagrantfile.
	TargetNames() (names []string, err error)
	//TargetNamesAndOptions() (names []string, options map[string]interface{}, err error)
}
