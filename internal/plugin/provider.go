package plugin

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/hashicorp/go-argmapper"
	"google.golang.org/grpc"

	"github.com/hashicorp/vagrant-plugin-sdk/component"
  "github.com/hashicorp/vagrant-plugin-sdk/docs"
	"github.com/hashicorp/vagrant-plugin-sdk/proto/gen"
)

// ProviderPlugin implements plugin.Plugin (specifically GRPCPlugin) for
// the Provider component type.
type ProviderPlugin struct {
	plugin.NetRPCUnsupportedPlugin

	Impl    component.Provider // Impl is the concrete implementation
	Mappers []*argmapper.Func     // Mappers
	Logger  hclog.Logger          // Logger
}

func (p *ProviderPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	proto.RegisterProviderServer(s, &providerServer{
		Impl:    p.Impl,
		Mappers: p.Mappers,
		Logger:  p.Logger,
		Broker:  broker,
	})
	return nil
}

func (p *ProviderPlugin) GRPCClient(
	ctx context.Context,
	broker *plugin.GRPCBroker,
	c *grpc.ClientConn,
) (interface{}, error) {
	return &providerClient{
		client: proto.NewProviderClient(c),
		logger: p.Logger,
		broker: broker,
	}, nil
}

// providerClient is an implementation of component.Provider over gRPC.
type providerClient struct {
	client proto.ProviderClient
	logger hclog.Logger
	broker *plugin.GRPCBroker
	mappers []*argmapper.Func
}

func (c *providerClient) Config() (interface{}, error) {
	return configStructCall(context.Background(), c.client)
}

func (c *providerClient) ConfigSet(v interface{}) error {
	return configureCall(context.Background(), c.client, v)
}

func (c *providerClient) Documentation() (*docs.Documentation, error) {
	return documentationCall(context.Background(), c.client)
}

func (c *providerClient) ProviderFunc() interface{} {
	//TODO
	return nil
}


// logPlatformServer is a gRPC server that the client talks to and calls a
// real implementation of the component.
type providerServer struct {
	Impl    component.Provider
	Mappers []*argmapper.Func
	Logger  hclog.Logger
	Broker  *plugin.GRPCBroker
}

func (s *providerServer) ConfigStruct(
	ctx context.Context,
	empty *empty.Empty,
) (*proto.Config_StructResp, error) {
	return configStruct(s.Impl)
}

func (s *providerServer) Configure(
	ctx context.Context,
	req *proto.Config_ConfigureRequest,
) (*empty.Empty, error) {
	return configure(s.Impl, req)
}

func (s *providerServer) Documentation(
	ctx context.Context,
	empty *empty.Empty,
) (*proto.Config_Documentation, error) {
	return documentation(s.Impl)
}

var (
	_ plugin.Plugin           = (*ProviderPlugin)(nil)
	_ plugin.GRPCPlugin       = (*ProviderPlugin)(nil)
	_ proto.ProviderServer = (*providerServer)(nil)
	_ component.Provider   = (*providerClient)(nil)
)
