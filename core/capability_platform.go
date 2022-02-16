package core

type CapabilityPlatform interface {
	Capability(name string, args ...interface{}) (interface{}, error)
	HasCapability(name string) (bool, error)
}
