package core

import (
	"io"

	"google.golang.org/protobuf/types/known/anypb"

	"github.com/hashicorp/vagrant-plugin-sdk/datadir"
)

type Target interface {
	Name() (string, error)
	ResourceId() (string, error)
	Project() (Project, error)
	Metadata() (map[string]string, error)
	DataDir() (*datadir.Target, error)
	State() (State, error)
	Record() (*anypb.Any, error)

	io.Closer
}
