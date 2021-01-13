package datadir

import (
	"github.com/hashicorp/vagrant-plugin-sdk/helper/path"
)

// Machine is an implementation of Dir that encapsulates the directories for a
// single machine.
type Machine struct {
	Dir
}

// Component returns a Dir implementation scoped to a specific component.
func (m *Machine) Component(typ, name string) (*Component, error) {
	dir, err := NewScopedDir(m, path.NewPath("component").Join(typ, name).String())
	if err != nil {
		return nil, err
	}

	return &Component{Dir: dir}, nil
}
