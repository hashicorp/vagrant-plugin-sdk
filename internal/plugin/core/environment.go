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

// Project implements core.Project interface
type Project struct {
	c          *ProjectClient
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

// ProjectPlugin is just a GRPC client for a project
type ProjectPlugin struct {
	plugin.NetRPCUnsupportedPlugin
	Mappers []*argmapper.Func // Mappers
	Logger  hclog.Logger      // Logger
	Impl    core.Project
}

// Implements plugin.GRPCPlugin
func (p *ProjectPlugin) GRPCClient(
	ctx context.Context,
	broker *plugin.GRPCBroker,
	c *grpc.ClientConn,
) (interface{}, error) {
	return &ProjectClient{
		client:       vagrant_plugin_sdk.NewProjectServiceClient(c),
		ServerTarget: c.Target(),
		Mappers:      p.Mappers,
		Logger:       p.Logger,
		Broker:       broker,
	}, nil
}

func (p *ProjectPlugin) GRPCServer(
	broker *plugin.GRPCBroker,
	s *grpc.Server,
) error {
	return errors.New("Server plugin not provided")
}

func NewProject(client *ProjectClient) *Project {
	return &Project{
		c:          client,
		ServerAddr: client.ServerTarget,
	}
}

type ProjectClient struct {
	Broker       *plugin.GRPCBroker
	Logger       hclog.Logger
	Mappers      []*argmapper.Func
	ServerTarget string
	client       vagrant_plugin_sdk.ProjectServiceClient
}

func (e *Project) CWD() (path string, err error) {
	return e.Cwd, nil
}

func (e *Project) DataDir() (path string, err error) {
	return e.Datadir, nil
}

func (e *Project) VagrantfileName() (name string, err error) {
	return e.Vagrantfilename, nil
}

func (e *Project) UI() (ui terminal.UI, err error) {
	return
}

func (e *Project) Home() (path string, err error) {
	return e.HomePath, nil
}
func (e *Project) LocalData() (path string, err error) {
	return e.LocalDataPath, nil
}

func (e *Project) Tmp() (path string, err error) {
	return e.TmpPath, nil
}

func (e *Project) DefaultPrivateKey() (path string, err error) {
	return e.DefaultPrivateKeyPath, nil
}

func (e *Project) MachineNames() (names []string, err error) {
	r, err := e.c.client.MachineNames(context.Background(), &empty.Empty{})
	if err != nil {
		return
	}
	names = r.Names
	return
}

func (e *Project) Host() (h core.Host, err error) {
	// TODO
	return nil, nil
}

var (
	_ plugin.Plugin     = (*ProjectPlugin)(nil)
	_ plugin.GRPCPlugin = (*ProjectPlugin)(nil)
	_ core.Project      = (*Project)(nil)
)
