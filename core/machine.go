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
	// accessors
	// Box() (box Box, err error)
	// DataDir() (path string, err error)
	// Environment() (env Environment, err error)
	// ID() (id string, err error)
	// Name() (name string, err error)
	// Provider() (p Provider, err error)
	// // ProviderConfig
	// ProviderName() (name string, err error)
	// // Triggers
	// UI() (ui terminal.UI, err error)
	// Vagrantfile() (file Vagrantfile, err error)

	// actual workers
	ActionFunc() interface{}
	CommunicateFunc() interface{}
	GuestFunc() interface{}
	SetIDFunc() interface{}
	IndexUUIDFunc() interface{}
	InspectFunc() interface{}
	ReloadFunc() interface{}
	ConnectionInfoFunc() interface{}
	StateFunc() interface{}
	UIDFunc() interface{}
	SyncedFoldersFunc() interface{}
	// Action(name string, state multistep.StateBag) (err error)
	// Communicate() (comm Communicator, err error)
	// Guest() (g Guest, err error)
	// SetID(value string) (err error)
	// IndexUUID() (id string, err error)
	// Inspect() (printable string, err error)
	// Reload() (err error)
	// ConnectionInfo() (info ConnectionInfo, err error)
	// State() (state *MachineState, err error)
	// UID() (user_id int, err error)
	// SyncedFolders() (folders []SyncedFolder, err error)
}
