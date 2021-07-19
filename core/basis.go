package core

import (
	"io"

	"github.com/hashicorp/vagrant-plugin-sdk/datadir"
	"github.com/hashicorp/vagrant-plugin-sdk/terminal"
)

//go:generate mockery --all

type Basis interface {
	DataDir() (dir *datadir.Basis, err error)
	UI() (ui terminal.UI, err error)
	Host() (host Host, err error)

	io.Closer
}
