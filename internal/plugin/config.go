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
	"github.com/hashicorp/vagrant-plugin-sdk/proto/vagrant_plugin_sdk"
)

// ConfigPlugin implements plugin.Plugin (specifically GRPCPlugin) for
// the Config component type.
type ConfigPlugin struct {
	plugin.NetRPCUnsupportedPlugin

	Impl    component.Config  // Impl is the concrete implementation
	Mappers []*argmapper.Func // Mappers
	Logger  hclog.Logger      // Logger
}

func (p *ConfigPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	vagrant_plugin_sdk.RegisterConfigServiceServer(s, &configServer{
		Impl:    p.Impl,
		Mappers: p.Mappers,
		Logger:  p.Logger,
		Broker:  broker,
	})
	return nil
}

func (p *ConfigPlugin) GRPCClient(
	ctx context.Context,
	broker *plugin.GRPCBroker,
	c *grpc.ClientConn,
) (interface{}, error) {
	return &configClient{
		client: vagrant_plugin_sdk.NewConfigServiceClient(c),
		logger: p.Logger,
		broker: broker,
	}, nil
}

// configClient is an implementation of component.Config over gRPC.
type configClient struct {
	client  vagrant_plugin_sdk.ConfigServiceClient
	logger  hclog.Logger
	broker  *plugin.GRPCBroker
	mappers []*argmapper.Func
}

func (c *configClient) Config() (interface{}, error) {
	return configStructCall(context.Background(), c.client)
}

func (c *configClient) ConfigSet(v interface{}) error {
	return configureCall(context.Background(), c.client, v)
}

func (c *configClient) Documentation() (*docs.Documentation, error) {
	return documentationCall(context.Background(), c.client)
}

func (c *configClient) ConfigFunc() interface{} {
	//TODO
	return nil
}

// configServer is a gRPC server that the client talks to and calls a
// real implementation of the component.
type configServer struct {
	Impl    component.Config
	Mappers []*argmapper.Func
	Logger  hclog.Logger
	Broker  *plugin.GRPCBroker

	vagrant_plugin_sdk.UnimplementedConfigServiceServer
}

func (s *configServer) ConfigStruct(
	ctx context.Context,
	empty *empty.Empty,
) (*vagrant_plugin_sdk.Config_StructResp, error) {
	return configStruct(s.Impl)
}

func (s *configServer) Configure(
	ctx context.Context,
	req *vagrant_plugin_sdk.Config_ConfigureRequest,
) (*empty.Empty, error) {
	return configure(s.Impl, req)
}

func (s *configServer) Documentation(
	ctx context.Context,
	empty *empty.Empty,
) (*vagrant_plugin_sdk.Config_Documentation, error) {
	return documentation(s.Impl)
}

var (
	_ plugin.Plugin                          = (*ConfigPlugin)(nil)
	_ plugin.GRPCPlugin                      = (*ConfigPlugin)(nil)
	_ vagrant_plugin_sdk.ConfigServiceServer = (*configServer)(nil)
	_ component.Config                       = (*configClient)(nil)
)
