// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package datadir

import (
	"github.com/hashicorp/vagrant-plugin-sdk/helper/path"
)

// Target is an implementation of Dir that encapsulates the directories for a
// single target.
type Target struct {
	Dir
}

// Target returns the Dir implementation scoped to a specific subtarget.
func (p *Target) Target(name string) (*Target, error) {
	dir, err := NewScopedDir(p, path.NewPath("subtarget").Join(name).String())
	if err != nil {
		return nil, err
	}

	return &Target{Dir: dir}, nil
}

// Component returns a Dir implementation scoped to a specific component.
func (m *Target) Component(typ, name string) (*Component, error) {
	dir, err := NewScopedDir(m, path.NewPath("component").Join(typ, name).String())
	if err != nil {
		return nil, err
	}

	return &Component{Dir: dir}, nil
}
