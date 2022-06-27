package core

import (
	"io"
	"time"

	"google.golang.org/protobuf/types/known/anypb"

	"github.com/hashicorp/vagrant-plugin-sdk/datadir"
	"github.com/hashicorp/vagrant-plugin-sdk/terminal"
)

type Target interface {
	Communicate() (comm Communicator, err error)
	//	Config() (config interface{}, err error)
	DataDir() (*datadir.Target, error)
	Destroy() error
	GetUUID() (id string, err error)
	Metadata() (map[string]string, error)
	Name() (string, error)
	Project() (Project, error)
	Provider() (p Provider, err error)
	ProviderName() (name string, err error)
	Record() (*anypb.Any, error)
	ResourceId() (string, error)
	Save() error
	SetName(value string) (err error)
	SetUUID(id string) (err error)
	Specialize(kind interface{}) (specialized interface{}, err error)
	State() (State, error)
	UI() (ui terminal.UI, err error)
	UpdatedAt() (t *time.Time, err error)
	Vagrantfile() (v Vagrantfile, err error)

	io.Closer
}
