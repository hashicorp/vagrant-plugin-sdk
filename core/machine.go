package core

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
	Target

	Box() (b *Box, err error)
	ConnectionInfo() (info *ConnectionInfo, err error)
	Guest() (g Guest, err error)
	ID() (id string, err error)
	Inspect() (printable string, err error)
	MachineState() (state *MachineState, err error)
	Reload() (err error)
	SetID(value string) (err error)
	SetMachineState(state *MachineState) (err error)
	SyncedFolders() (folders []SyncedFolder, err error)
	UID() (userId string, err error)
}
