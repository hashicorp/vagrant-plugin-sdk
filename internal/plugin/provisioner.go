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

// ProvisionerPlugin implements plugin.Plugin (specifically GRPCPlugin) for
// the Provisioner component type.
type ProvisionerPlugin struct {
	plugin.NetRPCUnsupportedPlugin

	Impl    component.Provisioner // Impl is the concrete implementation
	Mappers []*argmapper.Func     // Mappers
	Logger  hclog.Logger          // Logger
}

func (p *ProvisionerPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	proto.RegisterProvisionerServiceServer(s, &provisionerServer{
		Impl:    p.Impl,
		Mappers: p.Mappers,
		Logger:  p.Logger,
		Broker:  broker,
	})
	return nil
}

func (p *ProvisionerPlugin) GRPCClient(
	ctx context.Context,
	broker *plugin.GRPCBroker,
	c *grpc.ClientConn,
) (interface{}, error) {
	return &provisionerClient{
		client: proto.NewProvisionerServiceClient(c),
		logger: p.Logger,
		broker: broker,
	}, nil
}

// provisionerClient is an implementation of component.Provisioner over gRPC.
type provisionerClient struct {
	client  proto.ProvisionerServiceClient
	logger  hclog.Logger
	broker  *plugin.GRPCBroker
	mappers []*argmapper.Func
}

func (c *provisionerClient) Config() (interface{}, error) {
	return configStructCall(context.Background(), c.client)
}

func (c *provisionerClient) ConfigSet(v interface{}) error {
	return configureCall(context.Background(), c.client, v)
}

func (c *provisionerClient) Documentation() (*docs.Documentation, error) {
	return documentationCall(context.Background(), c.client)
}

func (c *provisionerClient) ProvisionerFunc() interface{} {
	//TODO
	return nil
}

// provisionerServer is a gRPC server that the client talks to and calls a
// real implementation of the component.
type provisionerServer struct {
	Impl    component.Provisioner
	Mappers []*argmapper.Func
	Logger  hclog.Logger
	Broker  *plugin.GRPCBroker
}

func (s *provisionerServer) ConfigStruct(
	ctx context.Context,
	empty *empty.Empty,
) (*proto.Config_StructResp, error) {
	return configStruct(s.Impl)
}

func (s *provisionerServer) Configure(
	ctx context.Context,
	req *proto.Config_ConfigureRequest,
) (*empty.Empty, error) {
	return configure(s.Impl, req)
}

func (s *provisionerServer) Documentation(
	ctx context.Context,
	empty *empty.Empty,
) (*proto.Config_Documentation, error) {
	return documentation(s.Impl)
}

var (
	_ plugin.Plugin                  = (*ProvisionerPlugin)(nil)
	_ plugin.GRPCPlugin              = (*ProvisionerPlugin)(nil)
	_ proto.ProvisionerServiceServer = (*provisionerServer)(nil)
	_ component.Provisioner          = (*provisionerClient)(nil)
)
