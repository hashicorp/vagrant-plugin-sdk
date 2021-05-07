package core

import (
	"context"
	"errors"

	"github.com/hashicorp/go-argmapper"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/hashicorp/vagrant-plugin-sdk/core"
	"github.com/hashicorp/vagrant-plugin-sdk/proto/vagrant_plugin_sdk"
	"google.golang.org/grpc"
)

// Vagrantfile implements core.Vagrantfile interface
type Vagrantfile struct {
	c          *VagrantfileClient
	ServerAddr string
}

// VagrantfilePlugin is just a GRPC client for a vagrantfile
type VagrantfilePlugin struct {
	plugin.NetRPCUnsupportedPlugin
	Mappers []*argmapper.Func // Mappers
	Logger  hclog.Logger      // Logger
	Impl    core.Vagrantfile
}

// Implements plugin.GRPCPlugin
func (p *VagrantfilePlugin) GRPCClient(
	ctx context.Context,
	broker *plugin.GRPCBroker,
	c *grpc.ClientConn,
) (interface{}, error) {
	return &VagrantfileClient{
		client:       vagrant_plugin_sdk.NewVagrantfileServiceClient(c),
		ServerTarget: c.Target(),
		Mappers:      p.Mappers,
		Logger:       p.Logger,
		Broker:       broker,
	}, nil
}

func (p *VagrantfilePlugin) GRPCServer(
	broker *plugin.GRPCBroker,
	s *grpc.Server,
) error {
	return errors.New("Server plugin not provided")
}

func NewVagrantfile(client *VagrantfileClient) *Vagrantfile {
	return &Vagrantfile{
		c:          client,
		ServerAddr: client.ServerTarget,
	}
}

type VagrantfileClient struct {
	Broker       *plugin.GRPCBroker
	Logger       hclog.Logger
	Mappers      []*argmapper.Func
	ServerTarget string
	client       vagrant_plugin_sdk.VagrantfileServiceClient
}

func (v *Vagrantfile) Machine(name, provider string, boxes core.BoxCollection, dataPath string, project core.Project) (machine core.Machine, err error) {
	return
}

func (v *Vagrantfile) MachineConfig(name, provider string, boxes core.BoxCollection, dataPath string, validateProvider bool) (config core.MachineConfig, err error) {
	return
}

func (v *Vagrantfile) MachineNames() (names []string, err error) {
	return
}

func (v *Vagrantfile) MachineNamesAndOptions() (names []string, options map[string]interface{}, err error) {
	return
}

func (v *Vagrantfile) PrimaryMachineName() (name string, err error) {
	return
}

var (
	_ plugin.Plugin     = (*VagrantfilePlugin)(nil)
	_ plugin.GRPCPlugin = (*VagrantfilePlugin)(nil)
	_ core.Vagrantfile  = (*Vagrantfile)(nil)
)
