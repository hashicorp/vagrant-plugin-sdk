// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package datadir

import (
	"github.com/hashicorp/vagrant-plugin-sdk/helper/path"
)

// Project is an implementation of Dir that encapsulates the directory
// for an entire project, including multiple apps.
//
// The paths returned by the Dir interface functions will be project-global.
// This means that the data is shared by all applications in the project.
type Project struct {
	Dir
}

// Target returns the Dir implementation scoped to a specific target.
func (p *Project) Target(ident string) (*Target, error) {
	dir, err := NewScopedDir(p, path.NewPath("target").Join(ident).String())
	if err != nil {
		return nil, err
	}

	return &Target{Dir: dir}, nil
}

// Assert implementation
var _ Dir = (*Project)(nil)
