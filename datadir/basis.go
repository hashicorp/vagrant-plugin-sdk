package datadir

import (
	"github.com/hashicorp/vagrant-plugin-sdk/helper/path"
)

type Basis struct {
	Dir
}

func NewBasis(path string) (*Basis, error) {
	dir, err := newRootDir(path)
	if err != nil {
		return nil, err
	}

	return &Basis{Dir: dir}, nil
}

// Project returns the Dir implementation scoped to a specific project.
func (p *Basis) Project(name string) (*Project, error) {
	dir, err := NewScopedDir(p, path.NewPath("target").Join(name).String())
	if err != nil {
		return nil, err
	}

	return &Project{Dir: dir}, nil
}

// Assert implementation
var _ Dir = (*Basis)(nil)
