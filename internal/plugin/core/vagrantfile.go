package core

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"

	"github.com/hashicorp/vagrant-plugin-sdk/core"
	"github.com/hashicorp/vagrant-plugin-sdk/helper/path"
	vplugin "github.com/hashicorp/vagrant-plugin-sdk/internal/plugin"
	"github.com/hashicorp/vagrant-plugin-sdk/proto/vagrant_plugin_sdk"
)

// VagrantfilePlugin is just a GRPC client for a vagrantfile
type VagrantfilePlugin struct {
	plugin.NetRPCUnsupportedPlugin

	Impl core.Vagrantfile
	*vplugin.BasePlugin
}

// Implements plugin.GRPCPlugin
func (p *VagrantfilePlugin) GRPCClient(
	ctx context.Context,
	broker *plugin.GRPCBroker,
	c *grpc.ClientConn,
) (interface{}, error) {
	return &vagrantfileClient{
		client:     vagrant_plugin_sdk.NewVagrantfileServiceClient(c),
		BaseClient: p.NewClient(ctx, broker, nil),
	}, nil
}

func (p *VagrantfilePlugin) GRPCServer(
	broker *plugin.GRPCBroker,
	s *grpc.Server,
) error {
	vagrant_plugin_sdk.RegisterVagrantfileServiceServer(s, &vagrantfileServer{
		Impl:       p.Impl,
		BaseServer: p.NewServer(broker, nil),
	})
	return nil
}

type vagrantfileClient struct {
	*vplugin.BaseClient

	client vagrant_plugin_sdk.VagrantfileServiceClient
}

type vagrantfileServer struct {
	*vplugin.BaseServer

	Impl core.Vagrantfile
	vagrant_plugin_sdk.UnimplementedVagrantfileServiceServer
}

func (v *vagrantfileClient) Target(name, provider string, boxes core.BoxCollection, dataPath path.Path, project core.Project) (machine core.Machine, err error) {
	return
}

func (v *vagrantfileClient) TargetConfig(name, provider string, boxes core.BoxCollection, dataPath path.Path, validateProvider bool) (config core.MachineConfig, err error) {
	return
}

func (v *vagrantfileClient) TargetNames() (names []string, err error) {
	return
}

func (v *vagrantfileClient) PrimaryTargetName() (name string, err error) {
	return
}

// Server

func (v *vagrantfileServer) Target(
	ctx context.Context,
	req *vagrant_plugin_sdk.Vagrantfile_TargetRequest,
) (*vagrant_plugin_sdk.Vagrantfile_TargetResponse, error) {
	return nil, nil
}

func (v *vagrantfileServer) TargetConfig(
	ctx context.Context,
	req *vagrant_plugin_sdk.Vagrantfile_TargetConfigRequest,
) (*vagrant_plugin_sdk.Vagrantfile_TargetConfigResponse, error) {
	return nil, nil
}

func (v *vagrantfileServer) TargetNames(
	ctx context.Context,
	_ *empty.Empty,
) (*vagrant_plugin_sdk.Vagrantfile_TargetNamesResponse, error) {
	return nil, nil
}

func (v *vagrantfileServer) PrimaryTargetName(
	ctx context.Context,
	_ *empty.Empty,
) (*vagrant_plugin_sdk.Vagrantfile_PrimaryTargetNameResponse, error) {
	return nil, nil
}

var (
	_ plugin.Plugin     = (*VagrantfilePlugin)(nil)
	_ plugin.GRPCPlugin = (*VagrantfilePlugin)(nil)
	_ core.Vagrantfile  = (*vagrantfileClient)(nil)
)
