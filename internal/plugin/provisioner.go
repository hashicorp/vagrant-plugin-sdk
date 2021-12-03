package plugin

import (
	"context"

	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/hashicorp/vagrant-plugin-sdk/component"
	"github.com/hashicorp/vagrant-plugin-sdk/core"
	"github.com/hashicorp/vagrant-plugin-sdk/docs"
	"github.com/hashicorp/vagrant-plugin-sdk/proto/vagrant_plugin_sdk"
)

// ProvisionerPlugin implements plugin.Plugin (specifically GRPCPlugin) for
// the Provisioner component type.
type ProvisionerPlugin struct {
	plugin.NetRPCUnsupportedPlugin

	Impl component.Provisioner // Impl is the concrete implementation
	*BasePlugin
}

func (p *ProvisionerPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	vagrant_plugin_sdk.RegisterProvisionerServiceServer(s, &provisionerServer{
		Impl:       p.Impl,
		BaseServer: p.NewServer(broker, p.Impl),
	})
	return nil
}

func (p *ProvisionerPlugin) GRPCClient(
	ctx context.Context,
	broker *plugin.GRPCBroker,
	c *grpc.ClientConn,
) (interface{}, error) {
	cl := vagrant_plugin_sdk.NewProvisionerServiceClient(c)
	return &provisionerClient{
		client:     cl,
		BaseClient: p.NewClient(ctx, broker, cl.(SeederClient)),
	}, nil
}

// provisionerClient is an implementation of component.Provisioner over gRPC.
type provisionerClient struct {
	*BaseClient

	client vagrant_plugin_sdk.ProvisionerServiceClient
}

func (c *provisionerClient) Config() (interface{}, error) {
	return configStructCall(c.Ctx, c.client)
}

func (c *provisionerClient) ConfigSet(v interface{}) error {
	return configureCall(c.Ctx, c.client, v)
}

func (c *provisionerClient) Documentation() (*docs.Documentation, error) {
	return documentationCall(c.Ctx, c.client)
}

func (c *provisionerClient) ProvisionerFunc() interface{} {
	//TODO
	return nil
}

// provisionerServer is a gRPC server that the client talks to and calls a
// real implementation of the component.
type provisionerServer struct {
	*BaseServer

	Impl component.Provisioner
	vagrant_plugin_sdk.UnsafeProvisionerServiceServer
}

func (s *provisionerServer) ConfigStruct(
	ctx context.Context,
	empty *emptypb.Empty,
) (*vagrant_plugin_sdk.Config_StructResp, error) {
	return configStruct(s.Impl)
}

func (s *provisionerServer) Configure(
	ctx context.Context,
	req *vagrant_plugin_sdk.Config_ConfigureRequest,
) (*emptypb.Empty, error) {
	return configure(s.Impl, req)
}

func (s *provisionerServer) Documentation(
	ctx context.Context,
	empty *emptypb.Empty,
) (*vagrant_plugin_sdk.Config_Documentation, error) {
	return documentation(s.Impl)
}

var (
	_ plugin.Plugin                               = (*ProvisionerPlugin)(nil)
	_ plugin.GRPCPlugin                           = (*ProvisionerPlugin)(nil)
	_ vagrant_plugin_sdk.ProvisionerServiceServer = (*provisionerServer)(nil)
	_ component.Provisioner                       = (*provisionerClient)(nil)
	_ core.Seeder                                 = (*provisionerClient)(nil)
)
