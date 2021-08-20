package core

import (
	"io"

	"github.com/hashicorp/vagrant-plugin-sdk/datadir"
	"github.com/hashicorp/vagrant-plugin-sdk/terminal"
)

type Basis interface {
	DataDir() (dir *datadir.Basis, err error)
	Host() (host Host, err error)
	UI() (ui terminal.UI, err error)

	io.Closer
}
