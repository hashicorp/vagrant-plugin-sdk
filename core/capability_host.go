package core

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
		return c.capabilities[capName]
	}
	return nil
}

func (c *CapabilityHost) RegisterCapability(name string, f interface{}) error {
	if c.capabilities == nil {
		c.capabilities = make(map[string]interface{})
	}
	c.capabilities[name] = f
	return nil
}
