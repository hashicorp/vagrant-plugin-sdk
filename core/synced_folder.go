// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package core

import (
	"github.com/hashicorp/vagrant-plugin-sdk/helper/path"
)

type Folder struct {
	Source      path.Path
	Destination path.Path
	Options     map[interface{}]interface{}
}

type SyncedFolder interface {
	CapabilityPlatform
	Seeder

	Usable(machine Machine) (bool, error)
	Enable(machine Machine, folders []*Folder, opts ...interface{}) error
	Prepare(machine Machine, folders []*Folder, opts ...interface{}) error
	Disable(machine Machine, folders []*Folder, opts ...interface{}) error
	Cleanup(machine Machine, opts ...interface{}) error
}
