package core

import (
	"github.com/hashicorp/vagrant-plugin-sdk/component"
)

type Provisioner interface {
	Provision(machine Machine, config *component.ConfigData) (err error)
	Configure(machine Machine, config *component.ConfigData, rootConfig *component.ConfigData) (err error)
	Cleanup(machine Machine, config *component.ConfigData) (err error)
}
