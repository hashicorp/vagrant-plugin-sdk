package core

import (
	"io"
	"time"

	"google.golang.org/protobuf/types/known/anypb"

	"github.com/hashicorp/vagrant-plugin-sdk/datadir"
	"github.com/hashicorp/vagrant-plugin-sdk/terminal"
)

type Target interface {
	Name() (string, error)
	SetName(value string) (err error)
	ResourceId() (string, error)
	Project() (Project, error)
	Metadata() (map[string]string, error)
	DataDir() (*datadir.Target, error)
	State() (State, error)
	UI() (ui terminal.UI, err error)
	UpdatedAt() (t *time.Time, err error)
	GetUUID() (id string, err error)
	SetUUID(id string) (err error)
	Provider() (p Provider, err error)
	Communicate() (comm Communicator, err error)

	Record() (*anypb.Any, error)
	Specialize(kind interface{}) (specialized interface{}, err error)

	Save() error

	io.Closer
}
