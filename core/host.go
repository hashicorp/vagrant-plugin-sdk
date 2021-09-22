package core

import (
	"io"
)

type Host interface {
	// Config() interface{}
	// Documentation() (*docs.Documentation, error)
	Capability(name string, args ...interface{}) (interface{}, error)
	Detect(state StateBag) (bool, error)
	HasCapability(name string) (bool, error)
	Parents() ([]string, error)

	io.Closer
}
