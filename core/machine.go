package core

import (
	"time"

	"github.com/hashicorp/vagrant-plugin-sdk/datadir"
	"github.com/hashicorp/vagrant-plugin-sdk/helper/path"
	"github.com/hashicorp/vagrant-plugin-sdk/terminal"
)

// the ssh info in vagrant core ends up dumping out
// a bunch of options, but they are also ssh specific
// where this would be used for other stuff too (like
// winrm). need to think on this some more.
type ConnectionInfo map[string]interface{}

type MachineState struct {
	ID               string
	ShortDescription string
	LongDescription  string
}

type Machine interface {
	// Attributes
	Name() (name string, err error)
	SetName(value string) (err error)
	ID() (id string, err error)
	SetID(value string) (err error)
	Box() (b Box, err error)
	DataDir() (d *datadir.Target, err error)
	LocalDataPath() (p path.Path, err error) // TODO(spox): difference from data dir?
	Provider() (p Provider, err error)
	VagrantfileName() (name string, err error)
	VagrantfilePath() (p path.Path, err error)
	UpdatedAt() (t *time.Time, err error)
	UI() (ui *terminal.UI, err error)

	// Functions
	Communicate() (comm Communicator, err error)
	Guest() (g Guest, err error)
	IndexUUID() (id string, err error)
	Inspect() (printable string, err error)
	Reload() (err error)
	ConnectionInfo() (info *ConnectionInfo, err error)
	State() (state *MachineState, err error)
	UID() (userId int, err error)
	SyncedFolders() (folders []SyncedFolder, err error)
}
