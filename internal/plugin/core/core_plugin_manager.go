package core

import (
	"context"

	"github.com/hashicorp/go-argmapper"
	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/hashicorp/vagrant-plugin-sdk/core"
	vplugin "github.com/hashicorp/vagrant-plugin-sdk/internal/plugin"
	"github.com/hashicorp/vagrant-plugin-sdk/proto/vagrant_plugin_sdk"
)

type CorePluginManagerPlugin struct {
	plugin.NetRPCUnsupportedPlugin

	Impl core.CorePluginManager
	*vplugin.BasePlugin
}

func (p *CorePluginManagerPlugin) GRPCClient(
	ctx context.Context,
	broker *plugin.GRPCBroker,
	c *grpc.ClientConn,
) (interface{}, error) {
	return &corePluginManagerClient{
		client:     vagrant_plugin_sdk.NewCorePluginManagerServiceClient(c),
		BaseClient: p.NewClient(ctx, broker, nil),
	}, nil
}

func (p *CorePluginManagerPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	vagrant_plugin_sdk.RegisterCorePluginManagerServiceServer(s, &corePluginManagerServer{
		Impl:       p.Impl,
		BaseServer: p.NewServer(broker, nil),
	})
	return nil
}

type corePluginManagerClient struct {
	*vplugin.BaseClient

	client vagrant_plugin_sdk.CorePluginManagerServiceClient
}

func (p *corePluginManagerClient) GetPlugin(pluginType core.Type) (plg interface{}, err error) {
	r, err := p.client.GetPlugin(p.Ctx, &vagrant_plugin_sdk.CorePluginManager_GetPluginRequest{
		Type: core.TypeStringMap[pluginType],
	})
	if err != nil {
		return nil, err
	}

	return p.Map(r, core.TypeMap[pluginType], argmapper.Typed(p.Ctx))
}

type corePluginManagerServer struct {
	*vplugin.BaseServer

	Impl core.CorePluginManager
	vagrant_plugin_sdk.UnimplementedCorePluginManagerServiceServer
}

func (p *corePluginManagerServer) GetPlugin(
	ctx context.Context, in *vagrant_plugin_sdk.CorePluginManager_GetPluginRequest,
) (plg *vagrant_plugin_sdk.CorePluginManager_GetPluginResponse, err error) {
	var pluginType core.Type
	for k, v := range core.TypeStringMap {
		if v == in.Type {
			pluginType = k
		}
	}

	plugin, err := p.Impl.GetPlugin(pluginType)
	if err != nil {
		return nil, err
	}

	pluginProto, err := p.Map(plugin, (**anypb.Any)(nil),
		argmapper.Typed(ctx))
	if err != nil {
		return nil, err
	}

	return &vagrant_plugin_sdk.CorePluginManager_GetPluginResponse{Plugin: pluginProto.(*anypb.Any)}, nil
}

var (
	_ plugin.Plugin                                 = (*PluginManagerPlugin)(nil)
	_ plugin.GRPCPlugin                             = (*PluginManagerPlugin)(nil)
	_ vagrant_plugin_sdk.PluginManagerServiceServer = (*pluginManagerServer)(nil)
	_ core.CorePluginManager                        = (*corePluginManagerClient)(nil)
)
