package core

import (
	"github.com/hashicorp/go-argmapper"
)

type CapabilityHost struct {
	capabilities map[string]interface{}
}

func (c *CapabilityHost) HasCapabilityFunc() interface{} {
	return c.HasCapability
}

func (c *CapabilityHost) HasCapability(name string) bool {
	if _, ok := c.capabilities[name]; ok {
		return true
	}
	return false
}

func (c *CapabilityHost) CapabilityFunc(capName string) interface{} {
	if c.HasCapability(capName) {
		f, _ := c.capabilities[capName]
		return f
	}
	return nil
}

func (c *CapabilityHost) Capability(capName string, args ...argmapper.Arg) (interface{}, error) {
	f := c.CapabilityFunc(capName)

	mapF, err := argmapper.NewFunc(f)
	if err != nil {
		return nil, err
	}

	callResult := mapF.Call(args...)
	if err := callResult.Err(); err != nil {
		return nil, err
	}

	raw := callResult.Out(0)
	return raw, nil
}

func (c *CapabilityHost) RegisterCapability(name string, f interface{}) error {
	if c.capabilities == nil {
		c.capabilities = make(map[string]interface{})
	}
	c.capabilities[name] = f
	return nil
}
