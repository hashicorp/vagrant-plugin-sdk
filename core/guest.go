package core

import (
	"github.com/hashicorp/vagrant-plugin-sdk/docs"
)

type Guest interface {
	Config() interface{}
	Documentation() (*docs.Documentation, error)
	Detect(machine Machine) (detected bool, err error)
}
