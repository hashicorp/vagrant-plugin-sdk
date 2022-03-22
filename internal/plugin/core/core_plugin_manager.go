package core

import (
	"context"

	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"

	"github.com/hashicorp/vagrant-plugin-sdk/core"
	vplugin "github.com/hashicorp/vagrant-plugin-sdk/internal/plugin"
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
		// client:     vagrant_plugin_sdk.NewPluginManagerServiceClient(c),
		BaseClient: p.NewClient(ctx, broker, nil),
	}, nil
}

func (p *CorePluginManagerPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	// vagrant_plugin_sdk.RegisterPluginManagerServiceServer(s, &pluginManagerServer{
	// 	Impl:       p.Impl,
	// 	BaseServer: p.NewServer(broker, nil),
	// })
	return nil
}

type corePluginManagerClient struct {
	*vplugin.BaseClient

	// client vagrant_plugin_sdk.PluginManagerServiceClient
}

func (m *corePluginManagerClient) GetPlugin(pluginType core.Type) (plg interface{}, err error) {
	return
}

type corePluginManagerServer struct {
	*vplugin.BaseServer

	Impl core.CorePluginManager
	// vagrant_plugin_sdk.UnimplementedPluginManagerServiceServer
}

var (
	_ plugin.Plugin     = (*PluginManagerPlugin)(nil)
	_ plugin.GRPCPlugin = (*PluginManagerPlugin)(nil)
	// _ vagrant_plugin_sdk.PluginManagerServiceServer = (*pluginManagerServer)(nil)
	_ core.CorePluginManager = (*corePluginManagerClient)(nil)
)
