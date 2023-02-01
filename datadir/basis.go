// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package datadir

import (
	"github.com/hashicorp/vagrant-plugin-sdk/helper/path"
)

type Basis struct {
	Dir
}

func NewBasis(ident string) (*Basis, error) {
	dir, err := newDir(ident)
	if err != nil {
		return nil, err
	}

	return &Basis{Dir: dir}, nil
}

// Project returns the Dir implementation scoped to a specific project.
func (p *Basis) Project(ident string) (*Project, error) {
	dir, err := NewScopedDir(p, path.NewPath("project").Join(ident).String())
	if err != nil {
		return nil, err
	}

	return &Project{Dir: dir}, nil
}

// Assert implementation
var _ Dir = (*Basis)(nil)
