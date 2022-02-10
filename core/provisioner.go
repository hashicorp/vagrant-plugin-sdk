package core

type ProvisionerConfig struct {
	Config map[string]interface{}
}

type Provisioner interface {
	Provision(machine Machine, config ProvisionerConfig) (err error)
	Configure(machine Machine, config ProvisionerConfig, rootConfig Vagrantfile) (err error)
	Cleanup(machine Machine, config ProvisionerConfig) (err error)
}
