package core

// import (
// 	"github.com/hashicorp/vagrant-plugin-sdk/docs"
// )

type Host interface {
	// Config() interface{}
	// Documentation() (*docs.Documentation, error)
	Detect() (detected bool)
	HasCapability(name string) bool
	CapabilityFunc(name string) interface{}
}
