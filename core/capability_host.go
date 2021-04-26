package core

import (
	"github.com/hashicorp/go-argmapper"
)

// TODO: chain of parents for "inheritance" of capabilities
// TODO: this map should have a bunch of argmapper functions
type CapabilityHost struct {
	capabilities map[string]func() interface{}
}

func (c *CapabilityHost) HasCapability(name string) bool {
	if _, ok := c.capabilities[name]; ok {
		return true
	}
	return false
}

func (c *CapabilityHost) Capability(name string, args ...argmapper.Arg) (interface{}, error) {
	f := c.capabilities[name]()
	// TODO: append converters and loggers to args
	// callArgs = append(args,
	// 	argmapper.ConverterFunc(b.Mappers...),
	// 	argmapper.Logger(b.Logger),
	// )
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

func (c *CapabilityHost) RegisterCapability(name string, f func() interface{}) error {
	if c.capabilities == nil {
		c.capabilities = make(map[string]func() interface{})
	}
	c.capabilities[name] = f
	return nil
}
