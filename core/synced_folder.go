package core

type Folder struct {
	Source      string
	Destination string
	Options     map[string]interface{} `mapstructure:",remain"`
}

type SyncedFolder interface {
	Capability(name string, args ...interface{}) (interface{}, error)
	HasCapability(name string) (bool, error)
	Usable(machine Machine) (bool, error)
	Enable(machine Machine, folders []*Folder, opts ...interface{}) error
	Prepare(machine Machine, folders []*Folder, opts ...interface{}) error
	Disable(machine Machine, folders []*Folder, opts ...interface{}) error
	Cleanup(machine Machine, opts ...interface{}) error
}
