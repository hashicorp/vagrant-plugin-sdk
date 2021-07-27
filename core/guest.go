package core

import (
	"io"

	"github.com/hashicorp/vagrant-plugin-sdk/docs"
)

type Guest interface {
	Config() interface{}
	Documentation() (*docs.Documentation, error)
	Parents() ([]string, error)
	Detect() (bool, error)
	HasCapability(name string) (bool, error)
	Capability(name string, args ...interface{}) (interface{}, error)

	io.Closer
}
