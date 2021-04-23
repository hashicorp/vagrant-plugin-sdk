package core

// TODO: chain of parents for "inheritance" of capabilities
// TODO: this map should have a bunch of argmapper functions
type CapabilityHost struct {
	capabilities map[string]func()
}

func (c *CapabilityHost) HasCapability(name string) bool {
	if _, ok := c.capabilities[name]; ok {
		return true
	}
	return false
}

func (c *CapabilityHost) Capability(name string) {
	c.capabilities[name]()
}

func (c *CapabilityHost) RegisterCapability(name string, f func()) error {
	c.capabilities[name] = f
	return nil
}
