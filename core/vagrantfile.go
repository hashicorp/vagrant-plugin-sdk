package core

import (
	"github.com/hashicorp/vagrant-plugin-sdk/component"
)

type Vagrantfile interface {
	GetConfig(namespace string) (*component.ConfigData, error)
	PrimaryTargetName() (name string, err error)
	Target(name, provider string) (Target, error)
	TargetConfig(name, provider string, validateProvider bool) (Vagrantfile, error)
	TargetNames() (names []string, err error)
	//TargetNamesAndOptions() (names []string, options map[string]interface{}, err error)
}
