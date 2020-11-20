// Package core has the interfaces for all core types
// used by Vagrant. These are implemented within Vagrant
// and provided to plugins as needed/required. This is
// currently a high level mapping of the interface
// provided by Vagrant in its current state. By doing
// a best effort match (and mapping where things are
// different) it should be easier to integrated with
// Vagrant to keep existing plugins working as well as
// making porting plugins less time consuming.

package core

type MachineIndexEntry struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	Provider        string `json:"provider"`
	DataPath        string `json:"data_path"`
	VagrantfilePath string `json:"vagrantfile_path"`
	State           string `json:"state"`
	UpdatedAt       string `json:"updated_at"`
}

// the ssh info in vagrant core ends up dumping out
// a bunch of options, but they are also ssh specific
// where this would be used for other stuff too (like
// winrm). need to think on this some more.
type ConnectionInfo map[string]interface{}

// for boxes
type Metadata map[string]interface{}

// for vagrantfile
type MachineConfig map[string]interface{}

type MachineState struct {
	ID               string
	ShortDescription string
	LongDescription  string
}

type BoxVersion struct {
	Version   string
	Providers []*BoxProvider
}

type BoxProvider struct {
	Name         string
	URL          string
	Checksum     string
	ChecksumType string
}

type BoxSummary struct {
	Name     string
	Version  string
	Provider string
}

type UI interface{}       // TODO(spox): what is UI in this context?
type StateBag interface{} // TODO(spox): this needs to be a utility (might be in another PR already?)

type MachineIndex interface {
	Delete(machine MachineEntry) (err error)
	// Each TODO(spox): enumerators?
	Get(uuid string) (entry *MachineEntry, err error)
	Includes(uuid string) (exists bool, err error)
	Release(entry *MachineEntry) (err error)
	Set(entry *MachineEntry) (updatedEntry *MachineEntry, err error)
	Recover(entry *MachineEntry) (updatedEntry *MachineEntry, err error)
}

// NOTE: Should include MachineIndexEntry
type MachineEntry interface {
	Valid(homePath string) (valid bool, err error)
	Environment(homePath string) (env Environment, err error)
}

type Box interface {
	// accessors
	Name() (name string, err error)
	Provider() (name string, err error)
	Version() (version string, err error)
	Directory() (path string, err error)
	Metadata() (metadata Metadata, err error)
	MetadataURL() (url string, err error)

	// action functions
	Destroy() (err error)
	InUse(index MachineIndex) (inUse bool, err error)
	LoadMetadata() (metadata BoxMetadata, err error)
	HasUpdate(version string) (updateAvailable bool, err error)
	AutomaticUpdateCheckAllowed() (allowed bool, err error)
	Repackage() (err error)

	// TODO(spox): Needs comparison function for sorting
}

type BoxMetadata interface {
	// accessors
	Name() (name string, err error)
	Description() (description string, err error)

	// action
	Load(pathOrURL string) (err error)
	Version(version string, providers []string) (v BoxVersion, err error)
	Versions(providers []string) (versions []BoxVersion, err error)
}

type BoxCollection interface {
	Add(path, name, version string, metadataURL string, providers []*BoxProvider) (box Box, err error)
	All() (list BoxSummary, err error)
	Find(name string, providers []string, version string) (box Box, err error)
	Clean(name string) (err error)
}

type Vagrantfile interface {
	Machine(name, provider string, boxes BoxCollection, dataPath string, env Environment) (machine Machine, err error)
	MachineConfig(name, provider string, boxes BoxCollection, dataPath string, validateProvider bool) (config MachineConfig, err error)
	MachineNames() (names []string, err error)
	MachineNamesAndOptions() (names []string, options map[string]interface{}, err error) // TODO(spox): dunno about this one
	PrimaryMachineName() (name string, err error)
}

type Machine interface {
	// accessors
	Box() (box Box, err error)
	DataDir() (path string, err error)
	Environment() (env Environment, err error)
	ID() (id string, err error)
	Name() (name string, err error)
	Provider() (p Provider, err error)
	// ProviderConfig
	ProviderName() (name string, err error)
	// Triggers
	UI() (ui UI, err error)
	Vagrantfile() (file Vagrantfile, err error)

	// actual workers
	Action(name string, state StateBag) (err error)
	Communicate() (comm Communicator, err error)
	Guest() (g Guest, err error)
	SetID(value string) (err error)
	IndexUUID() (id string, err error)
	Inspect() (printable string, err error)
	Reload() (err error)
	ConnectionInfo() (info ConnectionInfo, err error)
	State() (state MachineState, err error)
	UID() (user_id int, err error)
	SyncedFolders() (folders []SyncedFolder, err error)
}

type Environment interface {
	// accessors
	CWD() (path string, err error)
	DataDir() (path string, err error)
	VagrantfileName() (name string, err error)
	UI() (ui UI, err error)
	HomePath() (path string, err error)
	LocalDataPath() (path string, err error)
	TmpPath() (path string, err error)
	DefaultPrivateKeyPath() (path string, err error)

	// actual workers
	Inspect() (printable string, err error)
	ActiveMachines() (machines []Machine, err error)
	DefaultProvider() (name string, err error)
	CanInstallProvider() (can bool, err error)
	InstallProvider() (err error)
	Boxes() (boxes BoxCollection, err error)
	Environment(v Vagrantfile) (env Environment, err error)
	Hook(name string) (err error)
	Host() (h Host, err error)
	Lock(name string) (err error)
	Unlock(name string) (err error)
	// Push ?
	Machine(name, provider string, refresh bool) (machine Machine, err error)
	MachineIndex() (index MachineIndex, err error)
	MachineNames() (names []string, err error)
	PrimaryMachineName() (name string, err error)
	RootPath() (path string, err error)
	Unload() (err error)
	Vagrantfile() (v Vagrantfile, err error)
	SetupHomePath() (homePath string, err error) // TODO(spox): do we need this? probably not
	SetupLocalDataPath(force bool) (err error)   // TODO(spox): do we need this? - probably not
}

type Provider interface{}
type Host interface{}
type Guest interface{}
type SyncedFolder interface{}
type Communicator interface{}
