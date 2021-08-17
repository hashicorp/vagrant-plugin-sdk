package core

type MachineConfig struct {
	Box            Box
	Config         map[string]interface{}
	ConfigErrors   []string
	ConfigWarnings []string
	// Not to sure about these provider bits here
	Provider        string
	ProviderOptions map[string]interface{}
}

type Vagrantfile interface {
	// Returns the name of the machine that is designated as the
	// "primary."
	//
	// In the case of a single-machine environment, this is just the
	// single machine name. In the case of a multi-machine environment,
	// then this is the machine that is marked as primary, or nil if
	// no primary machine was specified.
	PrimaryTargetName() (name string, err error)

	// Returns a {Target} for the given name and provider that
	// is represented by this Vagrantfile.
	Target(name, provider string, boxes BoxCollection, dataPath string, env Project) (machine Machine, err error)

	// Returns the configuration for a single machine.
	//
	// When loading a box Vagrantfile, it will be prepended to the
	// key order specified when initializing this class. Sub-machine
	// and provider-specific overrides are appended at the end. The
	// actual order is:
	//
	// - box
	// - keys specified for #initialize
	// - sub-machine
	// - provider
	TargetConfig(name, provider string, boxes BoxCollection, dataPath string, validateProvider bool) (config MachineConfig, err error)

	// Returns a list of the machines that are defined within this
	// Vagrantfile.
	TargetNames() (names []string, err error)

	// Returns a list of the machine names as well as the options that
	// were specified for that machine.
	// TODO: Maybe uncomment this?
	// TargetNamesAndOptions() (names []string, options map[string]interface{}, err error)
}
