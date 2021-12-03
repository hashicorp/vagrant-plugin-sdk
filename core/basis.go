package core

import (
	"io"

	"github.com/hashicorp/vagrant-plugin-sdk/datadir"
	"github.com/hashicorp/vagrant-plugin-sdk/terminal"
)

type NamedPlugin struct {
	Plugin interface{}
	Name   string
	Type   string
}

type Basis interface {
	DataDir() (dir *datadir.Basis, err error)
	Host() (host Host, err error)
	UI() (ui terminal.UI, err error)
	Plugins(types ...string) (plugins []*NamedPlugin, err error)

	io.Closer
}
