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
	"github.com/hashicorp/vagrant-plugin-sdk/internal/pluginargs"
	"github.com/hashicorp/vagrant-plugin-sdk/proto/vagrant_plugin_sdk"
	"github.com/hashicorp/vagrant-plugin-sdk/terminal"
)

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
	return &projectClient{
		client: vagrant_plugin_sdk.NewProjectServiceClient(c),
		base: &base{
			Mappers: p.Mappers,
			Logger:  p.Logger,
			Broker:  broker,
			Cleanup: &pluginargs.Cleanup{},
		},
	}, nil
}

func (p *ProjectPlugin) GRPCServer(
	broker *plugin.GRPCBroker,
	s *grpc.Server,
) error {
	vagrant_plugin_sdk.RegisterProjectServiceServer(s, &projectServer{
		Impl: p.Impl,
		base: &base{
			Mappers: p.Mappers,
			Logger:  p.Logger,
			Broker:  broker,
			Cleanup: &pluginargs.Cleanup{},
		},
	})
	return nil
}

type projectClient struct {
	*base

	client vagrant_plugin_sdk.ProjectServiceClient
}

type projectServer struct {
	*base

	Impl core.Project
	vagrant_plugin_sdk.UnimplementedProjectServiceServer
}

func (p *projectClient) CWD() (path string, err error) {

	return e.Cwd, nil
}

func (p *projectClient) DataDir() (path string, err error) {
	return e.Datadir, nil
}

func (p *projectClient) VagrantfileName() (name string, err error) {
	return e.Vagrantfilename, nil
}

func (p *projectClient) UI() (ui terminal.UI, err error) {
	return
}

func (p *projectClient) Home() (path string, err error) {
	return e.HomePath, nil
}
func (p *projectClient) LocalData() (path string, err error) {
	return e.LocalDataPath, nil
}

func (p *projectClient) Tmp() (path string, err error) {
	return e.TmpPath, nil
}

func (p *projectClient) DefaultPrivateKey() (path string, err error) {
	return e.DefaultPrivateKeyPath, nil
}

func (p *projectClient) MachineNames() (names []string, err error) {
	r, err := e.c.client.MachineNames(context.Background(), &empty.Empty{})
	if err != nil {
		return
	}
	names = r.Names
	return
}

func (p *projectClient) Host() (h core.Host, err error) {
	// TODO
	return nil, nil
}

var (
	_ plugin.Plugin     = (*ProjectPlugin)(nil)
	_ plugin.GRPCPlugin = (*ProjectPlugin)(nil)
	_ core.Project      = (*projectClient)(nil)
)
