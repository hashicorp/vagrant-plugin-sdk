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
	"github.com/hashicorp/vagrant-plugin-sdk/proto/gen"
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
	proto.RegisterSyncedFolderServiceServer(s, &syncedFolderServer{
		Impl:    p.Impl,
		Mappers: p.Mappers,
		Logger:  p.Logger,
		Broker:  broker,
	})
	return nil
}

func (p *SyncedFolderPlugin) GRPCClient(
	ctx context.Context,
	broker *plugin.GRPCBroker,
	c *grpc.ClientConn,
) (interface{}, error) {
	return &syncedFolderClient{
		client: proto.NewSyncedFolderServiceClient(c),
		logger: p.Logger,
		broker: broker,
	}, nil
}

// syncedFolderClient is an implementation of component.SyncedFolder over gRPC.
type syncedFolderClient struct {
	client  proto.SyncedFolderServiceClient
	logger  hclog.Logger
	broker  *plugin.GRPCBroker
	mappers []*argmapper.Func
}

func (c *syncedFolderClient) Config() (interface{}, error) {
	return configStructCall(context.Background(), c.client)
}

func (c *syncedFolderClient) ConfigSet(v interface{}) error {
	return configureCall(context.Background(), c.client, v)
}

func (c *syncedFolderClient) Documentation() (*docs.Documentation, error) {
	return documentationCall(context.Background(), c.client)
}

func (c *syncedFolderClient) SyncedFolderFunc() interface{} {
	//TODO
	return nil
}

// syncedFolderServer is a gRPC server that the client talks to and calls a
// real implementation of the component.
type syncedFolderServer struct {
	Impl    component.SyncedFolder
	Mappers []*argmapper.Func
	Logger  hclog.Logger
	Broker  *plugin.GRPCBroker
}

func (s *syncedFolderServer) ConfigStruct(
	ctx context.Context,
	empty *empty.Empty,
) (*proto.Config_StructResp, error) {
	return configStruct(s.Impl)
}

func (s *syncedFolderServer) Configure(
	ctx context.Context,
	req *proto.Config_ConfigureRequest,
) (*empty.Empty, error) {
	return configure(s.Impl, req)
}

func (s *syncedFolderServer) Documentation(
	ctx context.Context,
	empty *empty.Empty,
) (*proto.Config_Documentation, error) {
	return documentation(s.Impl)
}

var (
	_ plugin.Plugin                   = (*SyncedFolderPlugin)(nil)
	_ plugin.GRPCPlugin               = (*SyncedFolderPlugin)(nil)
	_ proto.SyncedFolderServiceServer = (*syncedFolderServer)(nil)
	_ component.SyncedFolder          = (*syncedFolderClient)(nil)
)
