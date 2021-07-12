package core

import (
	"io"
)

type Host interface {
	// Config() interface{}
	// Documentation() (*docs.Documentation, error)
	Detect() (bool, error)
	HasCapability(name string) (bool, error)
	Capability(name string, args ...interface{}) (interface{}, error)

	io.Closer
}
