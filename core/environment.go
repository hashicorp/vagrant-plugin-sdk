package core

import (
	"github.com/hashicorp/vagrant-plugin-sdk/terminal"
)

type Environment interface {
	// accessors
	CWD() (path string, err error)
	DataDir() (path string, err error)
	VagrantfileName() (name string, err error)
	UI() (ui terminal.UI, err error)
	HomePath() (path string, err error)
	LocalDataPath() (path string, err error)
	TmpPath() (path string, err error)
	DefaultPrivateKeyPath() (path string, err error)

	// actual workers
	// Inspect() (printable string, err error)
	// ActiveMachines() (machines []Machine, err error)
	// DefaultProvider() (name string, err error)
	// CanInstallProvider() (can bool, err error)
	// InstallProvider() (err error)
	// Boxes() (boxes BoxCollection, err error)
	// Environment(v Vagrantfile) (env Environment, err error)
	// Hook(name string) (err error)
	// Host() (h Host, err error)
	// Lock(name string) (err error)
	// Unlock(name string) (err error)
	// Push ?
	// Machine(name, provider string, refresh bool) (machine Machine, err error)
	// MachineIndex() (index MachineIndex, err error)
	MachineNames() (names []string, err error)
	// PrimaryMachineName() (name string, err error)
	// RootPath() (path string, err error)
	// Unload() (err error)
	// Vagrantfile() (v Vagrantfile, err error)
	// SetupHomePath() (homePath string, err error) // TODO(spox): do we need this? probably not
	// SetupLocalDataPath(force bool) (err error)   // TODO(spox): do we need this? - probably not
}
