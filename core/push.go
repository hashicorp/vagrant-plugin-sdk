package core

import "github.com/hashicorp/vagrant-plugin-sdk/component"

type Push interface {
	component.Configurable

	Push() (err error)
}
