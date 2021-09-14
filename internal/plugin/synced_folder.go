package plugin

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/hashicorp/go-argmapper"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"

	"github.com/hashicorp/vagrant-plugin-sdk/component"
	"github.com/hashicorp/vagrant-plugin-sdk/docs"
	"github.com/hashicorp/vagrant-plugin-sdk/internal/pluginargs"
	"github.com/hashicorp/vagrant-plugin-sdk/proto/vagrant_plugin_sdk"
)

// SyncedFolderPlugin implements plugin.Plugin (specifically GRPCPlugin) for
// the SyncedFolder component type.
type SyncedFolderPlugin struct {
	plugin.NetRPCUnsupportedPlugin

	Impl    component.SyncedFolder // Impl is the concrete implementation
	Mappers []*argmapper.Func      // Mappers
	Logger  hclog.Logger           // Logger
}

func (p *SyncedFolderPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	vagrant_plugin_sdk.RegisterSyncedFolderServiceServer(s, &syncedFolderServer{
		Impl: p.Impl,
		baseServer: &baseServer{
			base: &base{
				Cleanup: &pluginargs.Cleanup{},
				Mappers: p.Mappers,
				Logger:  p.Logger,
				Broker:  broker,
			},
		},
	})
	return nil
}

func (p *SyncedFolderPlugin) GRPCClient(
	ctx context.Context,
	broker *plugin.GRPCBroker,
	c *grpc.ClientConn,
) (interface{}, error) {
	return &syncedFolderClient{
		client: vagrant_plugin_sdk.NewSyncedFolderServiceClient(c),
		baseClient: &baseClient{
			ctx: context.Background(),
			base: &base{
				Cleanup: &pluginargs.Cleanup{},
				Mappers: p.Mappers,
				Logger:  p.Logger,
				Broker:  broker,
			},
		},
	}, nil
}

// syncedFolderClient is an implementation of component.SyncedFolder over gRPC.
type syncedFolderClient struct {
	*baseClient

	client vagrant_plugin_sdk.SyncedFolderServiceClient
}

func (c *syncedFolderClient) Config() (interface{}, error) {
	return configStructCall(c.ctx, c.client)
}

func (c *syncedFolderClient) ConfigSet(v interface{}) error {
	return configureCall(c.ctx, c.client, v)
}

func (c *syncedFolderClient) Documentation() (*docs.Documentation, error) {
	return documentationCall(c.ctx, c.client)
}

func (c *syncedFolderClient) SyncedFolderFunc() interface{} {
	//TODO
	return nil
}

// syncedFolderServer is a gRPC server that the client talks to and calls a
// real implementation of the component.
type syncedFolderServer struct {
	*baseServer

	Impl component.SyncedFolder
	vagrant_plugin_sdk.UnimplementedSyncedFolderServiceServer
}

func (s *syncedFolderServer) ConfigStruct(
	ctx context.Context,
	empty *empty.Empty,
) (*vagrant_plugin_sdk.Config_StructResp, error) {
	return configStruct(s.Impl)
}

func (s *syncedFolderServer) Configure(
	ctx context.Context,
	req *vagrant_plugin_sdk.Config_ConfigureRequest,
) (*empty.Empty, error) {
	return configure(s.Impl, req)
}

func (s *syncedFolderServer) Documentation(
	ctx context.Context,
	empty *empty.Empty,
) (*vagrant_plugin_sdk.Config_Documentation, error) {
	return documentation(s.Impl)
}

var (
	_ plugin.Plugin                                = (*SyncedFolderPlugin)(nil)
	_ plugin.GRPCPlugin                            = (*SyncedFolderPlugin)(nil)
	_ vagrant_plugin_sdk.SyncedFolderServiceServer = (*syncedFolderServer)(nil)
	_ component.SyncedFolder                       = (*syncedFolderClient)(nil)
)
