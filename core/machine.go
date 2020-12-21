package core

import "github.com/hashicorp/vagrant-plugin-sdk/multistep"

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
	// actual workers
	Action(name string, state multistep.StateBag) (err error)
	Communicate() (comm Communicator, err error)
	Guest() (g Guest, err error)
	SetID(value string) (err error)
	IndexUUID() (id string, err error)
	Inspect() (printable string, err error)
	Reload() (err error)
	ConnectionInfo() (info ConnectionInfo, err error)
	State() (state *MachineState, err error)
	UID() (user_id int, err error)
	SyncedFolders() (folders []SyncedFolder, err error)
}

// type MachineClient interface {
// 	GetMachine(id string) (*Machine, error)
// 	ListMachines() ([]*Machine, error)
// 	UpsertMachine(machine *Machine)
// }
