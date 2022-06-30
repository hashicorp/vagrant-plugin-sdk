package core

import (
	"io"

	"github.com/hashicorp/vagrant-plugin-sdk/datadir"
	"github.com/hashicorp/vagrant-plugin-sdk/helper/path"
	"github.com/hashicorp/vagrant-plugin-sdk/terminal"
)

type Basis interface {
	Boxes() (boxes BoxCollection, err error)
	CWD() (path path.Path, err error)
	DataDir() (dir *datadir.Basis, err error)
	DefaultPrivateKey() (path path.Path, err error)
	Host() (host Host, err error)
	ResourceId() (string, error)
	TargetIndex() (index TargetIndex, err error)
	Vagrantfile() (Vagrantfile, error)
	UI() (ui terminal.UI, err error)

	io.Closer
}
