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

// NewProject creates the directory structure for a project. This will
// create the physical directories on disk if they do not already exist.
func NewProject(path string) (*Project, error) {
	dir, err := newRootDir(path)
	if err != nil {
		return nil, err
	}

	return &Project{Dir: dir}, nil
}

// App returns the Dir implementation scoped to a specific app.
func (p *Project) Machine(name string) (*Machine, error) {
	dir, err := NewScopedDir(p, path.NewPath("machine").Join(name).String())
	if err != nil {
		return nil, err
	}

	return &Machine{Dir: dir}, nil
}

// Assert implementation
var _ Dir = (*Project)(nil)
