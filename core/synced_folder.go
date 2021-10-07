package core

type Folder struct {
	Name    string
	Options map[string]string
}

type SyncedFolder interface {
	Capability(name string, args ...interface{}) (interface{}, error)
	HasCapability(name string) (bool, error)
	Usable(machine Machine) (bool, error)
	Enable(machine Machine, folders []*Folder, opts map[string]string) error
	Disable(machine Machine, folders []*Folder, opts map[string]string) error
	Cleanup(machine Machine, opts map[string]string) error
}
