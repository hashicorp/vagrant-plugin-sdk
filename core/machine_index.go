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

// NOTE: Should include MachineIndexEntry
type MachineEntry interface {
	Valid(homePath string) (valid bool, err error)
	Environment(homePath string) (env Environment, err error)
}

type MachineIndex interface {
	Delete(machine MachineEntry) (err error)
	// Each TODO(spox): enumerators?
	Get(uuid string) (entry MachineEntry, err error)
	Includes(uuid string) (exists bool, err error)
	Release(entry MachineEntry) (err error)
	Set(entry MachineEntry) (updatedEntry MachineEntry, err error)
	Recover(entry MachineEntry) (updatedEntry MachineEntry, err error)
}
