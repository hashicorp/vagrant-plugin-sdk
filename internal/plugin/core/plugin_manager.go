package core

import (
	"context"

	"github.com/hashicorp/go-argmapper"
	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"

	"github.com/hashicorp/vagrant-plugin-sdk/core"
	vplugin "github.com/hashicorp/vagrant-plugin-sdk/internal/plugin"
	"github.com/hashicorp/vagrant-plugin-sdk/proto/vagrant_plugin_sdk"
)

type PluginManagerPlugin struct {
	plugin.NetRPCUnsupportedPlugin

	Impl core.PluginManager
	*vplugin.BasePlugin
}

func (p *PluginManagerPlugin) GRPCClient(
	ctx context.Context,
	broker *plugin.GRPCBroker,
	c *grpc.ClientConn,
) (interface{}, error) {
	return &pluginManagerClient{
		client:     vagrant_plugin_sdk.NewPluginManagerServiceClient(c),
		BaseClient: p.NewClient(ctx, broker, nil),
	}, nil
}

func (p *PluginManagerPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	vagrant_plugin_sdk.RegisterPluginManagerServiceServer(s, &pluginManagerServer{
		Impl:       p.Impl,
		BaseServer: p.NewServer(broker, nil),
	})
	return nil
}

type pluginManagerClient struct {
	*vplugin.BaseClient

	client vagrant_plugin_sdk.PluginManagerServiceClient
}

type pluginManagerServer struct {
	*vplugin.BaseServer

	Impl core.PluginManager
	vagrant_plugin_sdk.UnimplementedPluginManagerServiceServer
}

func (p *pluginManagerClient) ListPlugins(types ...string) (h []*core.NamedPlugin, err error) {
	r, err := p.client.ListPlugins(p.Ctx, &vagrant_plugin_sdk.PluginManager_PluginsRequest{
		Types: types,
	})
	if err != nil {
		return
	}

	raw, err := p.Map(r, (*[]*core.NamedPlugin)(nil), argmapper.Typed(p.Ctx))
	if err != nil {
		return nil, err
	}

	return raw.([]*core.NamedPlugin), nil
}

func (p *pluginManagerClient) GetPlugin(name, typ string) (*core.NamedPlugin, error) {
	r, err := p.client.GetPlugin(p.Ctx, &vagrant_plugin_sdk.PluginManager_Plugin{
		Name: name,
		Type: typ,
	})
	if err != nil {
		return nil, err
	}

	raw, err := p.Map(r, (**core.NamedPlugin)(nil), argmapper.Typed(p.Ctx))
	if err != nil {
		return nil, err
	}

	return raw.(*core.NamedPlugin), nil
}

func (s *pluginManagerServer) ListPlugins(
	ctx context.Context,
	in *vagrant_plugin_sdk.PluginManager_PluginsRequest,
) (r *vagrant_plugin_sdk.PluginManager_PluginsResponse, err error) {
	plugins, err := s.Impl.ListPlugins(in.Types...)
	if err != nil {
		s.Logger.Error("failed to get plugin list",
			"types", in.Types,
			"error", err,
		)
		return
	}

	raw, err := s.Map(plugins, (**vagrant_plugin_sdk.PluginManager_PluginsResponse)(nil),
		argmapper.Typed(ctx))

	if err != nil {
		s.Logger.Error("failed to map plugin list",
			"error", err,
		)

		return nil, err
	}

	return raw.(*vagrant_plugin_sdk.PluginManager_PluginsResponse), nil
}

func (s *pluginManagerServer) GetPlugin(
	ctx context.Context,
	in *vagrant_plugin_sdk.PluginManager_Plugin,
) (*vagrant_plugin_sdk.PluginManager_Plugin, error) {
	p, err := s.Impl.GetPlugin(in.Name, in.Type)
	if err != nil {
		s.Logger.Error("failed to get plugin",
			"name", in.Name,
			"type", in.Type,
			"error", err,
		)
		return nil, err
	}

	raw, err := s.Map(p, (**vagrant_plugin_sdk.PluginManager_Plugin)(nil),
		argmapper.Typed(ctx))

	if err != nil {
		s.Logger.Error("failed to map plugin",
			"name", in.Name,
			"type", in.Type,
			"error", err,
		)

		return nil, err
	}

	return raw.(*vagrant_plugin_sdk.PluginManager_Plugin), nil
}

var (
	_ plugin.Plugin                                 = (*PluginManagerPlugin)(nil)
	_ plugin.GRPCPlugin                             = (*PluginManagerPlugin)(nil)
	_ vagrant_plugin_sdk.PluginManagerServiceServer = (*pluginManagerServer)(nil)
	_ core.PluginManager                            = (*pluginManagerClient)(nil)
)
