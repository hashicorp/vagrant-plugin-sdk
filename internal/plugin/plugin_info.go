package plugin

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/hashicorp/go-argmapper"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"

	"github.com/hashicorp/vagrant-plugin-sdk/component"
	"github.com/hashicorp/vagrant-plugin-sdk/proto/vagrant_plugin_sdk"
)

type pluginInfo struct {
	types []component.Type
}

func (p *pluginInfo) ComponentTypes() []component.Type {
	return p.types
}

type PluginInfoPlugin struct {
	plugin.NetRPCUnsupportedPlugin

	Impl    component.PluginInfo // Impl is the concrete implementation
	Mappers []*argmapper.Func    // Mappers
	Logger  hclog.Logger         // Logger
}

func (p *PluginInfoPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	vagrant_plugin_sdk.RegisterPluginInfoServiceServer(s, &pluginInfoServer{
		Impl: p.Impl,
		baseServer: &baseServer{
			base: &base{
				Mappers: p.Mappers,
				Logger:  p.Logger,
				Broker:  broker,
			},
		},
	})
	return nil
}

func (p *PluginInfoPlugin) GRPCClient(
	ctx context.Context,
	broker *plugin.GRPCBroker,
	c *grpc.ClientConn,
) (interface{}, error) {
	return &pluginInfoClient{
		client: vagrant_plugin_sdk.NewPluginInfoServiceClient(c),
		baseClient: &baseClient{
			ctx: context.Background(),
			base: &base{
				Mappers: p.Mappers,
				Logger:  p.Logger,
				Broker:  broker,
			},
		},
	}, nil
}

// pluginInfoClient is an implementation of component.PluginInfo over gRPC.
type pluginInfoClient struct {
	*baseClient

	client vagrant_plugin_sdk.PluginInfoServiceClient
}

type pluginInfoServer struct {
	*baseServer

	Impl component.PluginInfo
	vagrant_plugin_sdk.UnimplementedPluginInfoServiceServer
}

func (c *pluginInfoClient) ComponentTypes() (result []component.Type) {
	result = []component.Type{}
	resp, err := c.client.ComponentTypes(c.ctx, &empty.Empty{})
	if err != nil {
		c.Logger.Error("unexpected error when requesting component types",
			"error", err)
		return
	}
	for _, t := range resp.Component {
		result = append(result, component.Type(t))
	}
	return
}

func (s *pluginInfoServer) ComponentTypes(
	ctx context.Context,
	_ *empty.Empty,
) (*vagrant_plugin_sdk.PluginInfo_ComponentList, error) {
	list := []uint32{}
	for _, t := range s.Impl.ComponentTypes() {
		list = append(list, uint32(t))
	}
	return &vagrant_plugin_sdk.PluginInfo_ComponentList{
		Component: list,
	}, nil
}
