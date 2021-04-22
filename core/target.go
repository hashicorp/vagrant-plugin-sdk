package core

import (
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/hashicorp/vagrant-plugin-sdk/datadir"
)

type Target interface {
	Name() string
	ResourceID() string
	Project() Project
	Metadata() map[string]string
	DataDir() (datadir.Target, error)
	State() (State, error)
	Record() anypb.Any
}
