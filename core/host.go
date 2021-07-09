package core

// import (
// 	"github.com/hashicorp/vagrant-plugin-sdk/docs"
// )

type Host interface {
	// Config() interface{}
	// Documentation() (*docs.Documentation, error)
	Detect() (bool, error)
	HasCapability(name string) (bool, error)
	Capability(name string, args ...interface{}) (interface{}, error)
}
