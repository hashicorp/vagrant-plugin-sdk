package core

import (
	"context"
	"errors"

	"google.golang.org/grpc"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/hashicorp/go-argmapper"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"

	"github.com/hashicorp/vagrant-plugin-sdk/core"
	"github.com/hashicorp/vagrant-plugin-sdk/proto/vagrant_plugin_sdk"
	"github.com/hashicorp/vagrant-plugin-sdk/terminal"
)

// Environment implements core.Environment interface
type Environment struct {
	c          *EnvironmentClient
	ServerAddr string

	Cwd                   string
	Datadir               string
	Vagrantfilename       string
	HomePath              string
	LocalDataPath         string
	TmpPath               string
	AliasesPath           string
	BoxesPath             string
	GemsPath              string
	DefaultPrivateKeyPath string
}

// EnvironmentPlugin is just a GRPC client for a environment
type EnvironmentPlugin struct {
	plugin.NetRPCUnsupportedPlugin
	Mappers []*argmapper.Func // Mappers
	Logger  hclog.Logger      // Logger
	Impl    core.Environment
}

// Implements plugin.GRPCPlugin
func (p *EnvironmentPlugin) GRPCClient(
	ctx context.Context,
	broker *plugin.GRPCBroker,
	c *grpc.ClientConn,
) (interface{}, error) {
	return &EnvironmentClient{
		client:       vagrant_plugin_sdk.NewEnvironmentServiceClient(c),
		ServerTarget: c.Target(),
		Mappers:      p.Mappers,
		Logger:       p.Logger,
		Broker:       broker,
	}, nil
}

func (p *EnvironmentPlugin) GRPCServer(
	broker *plugin.GRPCBroker,
	s *grpc.Server,
) error {
	return errors.New("Server plugin not provided")
}

func NewEnvironment(client *EnvironmentClient) *Environment {
	return &Environment{
		c:          client,
		ServerAddr: client.ServerTarget,
	}
}

type EnvironmentClient struct {
	Broker       *plugin.GRPCBroker
	Logger       hclog.Logger
	Mappers      []*argmapper.Func
	ServerTarget string
	client       vagrant_plugin_sdk.EnvironmentServiceClient
}

func (e *Environment) CWD() (path string, err error) {
	return e.Cwd, nil
}

func (e *Environment) DataDir() (path string, err error) {
	return e.Datadir, nil
}

func (e *Environment) VagrantfileName() (name string, err error) {
	return e.Vagrantfilename, nil
}

func (e *Environment) UI() (ui terminal.UI, err error) {
	return
}

func (e *Environment) Home() (path string, err error) {
	return e.HomePath, nil
}
func (e *Environment) LocalData() (path string, err error) {
	return e.LocalDataPath, nil
}

func (e *Environment) Tmp() (path string, err error) {
	return e.TmpPath, nil
}

func (e *Environment) DefaultPrivateKey() (path string, err error) {
	return e.DefaultPrivateKeyPath, nil
}

func (e *Environment) MachineNames() (names []string, err error) {
	r, err := e.c.client.MachineNames(context.Background(), &empty.Empty{})
	if err != nil {
		return
	}
	names = r.Names
	return
}

var (
	_ plugin.Plugin     = (*EnvironmentPlugin)(nil)
	_ plugin.GRPCPlugin = (*EnvironmentPlugin)(nil)
	_ core.Environment  = (*Environment)(nil)
)
