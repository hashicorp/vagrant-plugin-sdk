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

	// Attributes
	ID() (id string, err error)
	SetID(value string) (err error)
	Box() (b Box, err error)

	// Functions
	Guest() (g Guest, err error)
	IndexUUID() (id string, err error)
	SetUUID(id string) (err error)
	Inspect() (printable string, err error)
	Reload() (err error)
	ConnectionInfo() (info *ConnectionInfo, err error)
	MachineState() (state *MachineState, err error)
	SetMachineState(state *MachineState) (err error)
	UID() (userId int, err error)
	SyncedFolders() (folders []SyncedFolder, err error)
}
