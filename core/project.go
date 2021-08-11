package core

import (
	"io"
	//	"github.com/hashicorp/vagrant-plugin-sdk/proto/vagrant_plugin_sdk"
	"github.com/hashicorp/vagrant-plugin-sdk/datadir"
	"github.com/hashicorp/vagrant-plugin-sdk/terminal"
)

type Project interface {
	// accessors
	CWD() (path string, err error)
	DataDir() (dir *datadir.Project, err error)
	VagrantfileName() (name string, err error)
	UI() (ui terminal.UI, err error)
	Home() (path string, err error)
	LocalData() (path string, err error)
	Tmp() (path string, err error)
	DefaultPrivateKey() (path string, err error)

	// actual workers
	// Inspect() (printable string, err error)
	// ActiveMachines() (machines []Machine, err error)
	// DefaultProvider() (name string, err error)
	// CanInstallProvider() (can bool, err error)
	// InstallProvider() (err error)
	// Boxes() (boxes BoxCollection, err error)
	// Project(v Vagrantfile) (env Project, err error)
	// Hook(name string) (err error)
	Host() (h Host, err error)
	// Lock(name string) (err error)
	// Unlock(name string) (err error)
	// Push ?
	// Machine(name, provider string, refresh bool) (machine Machine, err error)
	TargetIndex() (index TargetIndex, err error)
	MachineNames() (names []string, err error)
	// PrimaryMachineName() (name string, err error)
	// RootPath() (path string, err error)
	// Unload() (err error)
	// Vagrantfile() (v Vagrantfile, err error)
	// SetupHomePath() (homePath string, err error) // TODO(spox): do we need this? probably not
	// SetupLocalDataPath(force bool) (err error)   // TODO(spox): do we need this? - probably not

	Target(name string) (t Target, err error)
	TargetNames() (names []string, err error)
	TargetIds() (ids []string, err error)

	io.Closer
}
