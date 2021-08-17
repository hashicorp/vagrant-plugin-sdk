package core

import (
	"io"

	"github.com/hashicorp/vagrant-plugin-sdk/docs"
)

type Guest interface {
	Capability(name string, args ...interface{}) (interface{}, error)
	Config() interface{}
	Detect() (bool, error)
	Documentation() (*docs.Documentation, error)
	HasCapability(name string) (bool, error)
	Parents() ([]string, error)

	io.Closer
}
