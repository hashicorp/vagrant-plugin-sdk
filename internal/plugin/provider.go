package plugin

import (
	"context"
	"encoding/json"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/hashicorp/go-argmapper"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hashicorp/vagrant-plugin-sdk/component"
	"github.com/hashicorp/vagrant-plugin-sdk/internal/funcspec"
	"github.com/hashicorp/vagrant-plugin-sdk/internal/plugincomponent"
	"github.com/hashicorp/vagrant-plugin-sdk/proto/gen"
)

// ProviderPlugin implements plugin.Plugin (specifically GRPCPlugin) for
// the Provider component type.
type ProviderPlugin struct {
	plugin.NetRPCUnsupportedPlugin

	Impl    component.Provider // Impl is the concrete implementation
	Mappers []*argmapper.Func // Mappers
	Logger  hclog.Logger      // Logger
}

func (p *ProviderPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	base := &base{
		Mappers: p.Mappers,
		Logger:  p.Logger,
		Broker:  broker,
	}

	proto.RegisterProviderServer(s, &builderServer{
		base: base,
		Impl: p.Impl,
	})
	return nil
}

func (p *ProviderPlugin) GRPCClient(
	ctx context.Context,
	broker *plugin.GRPCBroker,
	c *grpc.ClientConn,
) (interface{}, error) {
	client := &providerClient{
		client:  proto.NewProviderClient(c),
		logger:  p.Logger,
		broker:  broker,
		mappers: p.Mappers,
	}

	return result, nil
}

// providerClient is an implementation of component.Provider that
// communicates over gRPC.
type providerClient struct {
	client  proto.ProviderClient
	logger  hclog.Logger
	broker  *plugin.GRPCBroker
	mappers []*argmapper.Func
}

func (c *ProviderClient) Config() (interface{}, error) {
	return configStructCall(context.Background(), c.client)
}

func (c *ProviderClient) ProviderFunc() interface{} {
	
}

// providerServer is a gRPC server that the client talks to and calls a
// real implementation of the component.
type providerServer struct {
	Impl component.Provider
}

var (
	_ plugin.Plugin                = (*ProviderPlugin)(nil)
	_ plugin.GRPCPlugin            = (*ProviderPlugin)(nil)
	_ proto.ProviderServer          = (*providerServer)(nil)
	_ component.Provider            = (*providerClient)(nil)
	_ component.Configurable       = (*providerClient)(nil)
	_ component.Documented         = (*providerClient)(nil)
	_ component.ConfigurableNotify = (*providerClient)(nil)
)
