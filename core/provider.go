package core

type SshInfo struct {
	Host           string
	Port           string
	Username       string
	PrivateKeyPath string
}

type Provider interface {
	Usable() (bool, error)
	Installed() (bool, error)
	Init(Machine) error
	Action(name string, args ...interface{}) error
	MachineIdChanged() error
	SshInfo() (*SshInfo, error)
	State() (*MachineState, error)

	// Capability functions
	Capability(name string, args ...interface{}) (interface{}, error)
	HasCapability(name string) (bool, error)
}
