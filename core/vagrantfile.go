package core

type MachineConfig map[string]interface{}

type Vagrantfile interface {
	Machine(name, provider string, boxes BoxCollection, dataPath string, env Environment) (machine Machine, err error)
	MachineConfig(name, provider string, boxes BoxCollection, dataPath string, validateProvider bool) (config MachineConfig, err error)
	MachineNames() (names []string, err error)
	MachineNamesAndOptions() (names []string, options map[string]interface{}, err error) // TODO(spox): dunno about this one
	PrimaryMachineName() (name string, err error)
}
