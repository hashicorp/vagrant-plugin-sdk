package core

import (
	"io"

	"github.com/hashicorp/vagrant-plugin-sdk/datadir"
	"github.com/hashicorp/vagrant-plugin-sdk/terminal"
)

type Basis interface {
	Boxes() (boxes BoxCollection, err error)
	CWD() (path string, err error)
	DataDir() (dir *datadir.Basis, err error)
	DefaultPrivateKey() (path string, err error)
	Host() (host Host, err error)
	TargetIndex() (index TargetIndex, err error)
	UI() (ui terminal.UI, err error)

	io.Closer
}
