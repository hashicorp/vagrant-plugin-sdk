package core

type CommunicatorOptions struct {
	// Return error when command fails
	ErrorCheck bool
	// Run command without modification
	ForceRaw bool
	// Run command as privileged
	Privileged bool
	// Valid exit code for command
	GoodExit int
	// Shell to use when running command
	Shell string
}

type Communicator interface {
	// Config() interface{}
	// Documentation() (*docs.Documentation, error)
	Download(machine Machine, source, destination string) error
	Execute(machine Machine, command []string, options *CommunicatorOptions) (status int32, err error)
	Init(machine Machine) error
	Match(machine Machine) (isMatch bool, err error)
	PrivilegedExecute(machine Machine, command []string, options *CommunicatorOptions) (status int32, err error)
	Ready(machine Machine) (isReady bool, err error)
	Reset(machine Machine) error
	Test(machine Machine, command []string, options *CommunicatorOptions) (valid bool, err error)
	Upload(machine Machine, source, destination string) error
	WaitForReady(machine Machine, wait int) (isReady bool, err error)
}
