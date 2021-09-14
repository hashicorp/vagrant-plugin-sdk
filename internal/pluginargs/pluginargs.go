// Package pluginargs
package pluginargs

import (
	"github.com/hashicorp/go-argmapper"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/hashicorp/vagrant-plugin-sdk/internal-shared/cacher"
)

// Internal is a struct that is available to mappers. This is an internal-only
// type that is not possible for plugins to register for since it is only
// exported in an internal package.
type Internal struct {
	Broker  *plugin.GRPCBroker
	Mappers []*argmapper.Func
	Cleanup *Cleanup
	Cache   cacher.Cache
	Logger  hclog.Logger
}

// Cleanup can be used to register cleanup functions.
type Cleanup struct {
	f func()
}

// Do registers a cleanup function that will be called when the plugin RPC
// call is complete.
func (c *Cleanup) Do(f func()) {
	oldF := c.f
	c.f = func() {
		if oldF != nil {
			defer oldF()
		}
		f()
	}
}

func (c *Cleanup) Close() error {
	if c.f != nil {
		c.f()
	}

	return nil
}
