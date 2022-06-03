package plugin

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/hashicorp/vagrant-plugin-sdk/component"
	"github.com/hashicorp/vagrant-plugin-sdk/proto/vagrant_plugin_sdk"
)

type pluginInfo struct {
	options map[component.Type]interface{}
	types   []component.Type
	name    string
}

func (p *pluginInfo) ComponentOptions() map[component.Type]interface{} {
	return p.options
}

func (p *pluginInfo) ComponentTypes() []component.Type {
	return p.types
}

func (p *pluginInfo) Name() string {
	return p.name
}

type PluginInfoPlugin struct {
	plugin.NetRPCUnsupportedPlugin

	Impl component.PluginInfo // Impl is the concrete implementation
	*BasePlugin
}

func (p *PluginInfoPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	vagrant_plugin_sdk.RegisterPluginInfoServiceServer(s, &pluginInfoServer{
		Impl:       p.Impl,
		BaseServer: p.NewServer(broker, nil),
	})
	return nil
}

func (p *PluginInfoPlugin) GRPCClient(
	ctx context.Context,
	broker *plugin.GRPCBroker,
	c *grpc.ClientConn,
) (interface{}, error) {
	return &pluginInfoClient{
		client:     vagrant_plugin_sdk.NewPluginInfoServiceClient(c),
		BaseClient: p.NewClient(ctx, broker, nil),
	}, nil
}

// pluginInfoClient is an implementation of component.PluginInfo over gRPC.
type pluginInfoClient struct {
	*BaseClient

	client vagrant_plugin_sdk.PluginInfoServiceClient
}

type pluginInfoServer struct {
	*BaseServer

	Impl component.PluginInfo
	vagrant_plugin_sdk.UnimplementedPluginInfoServiceServer
}

func (c *pluginInfoClient) ComponentOptions() (result map[component.Type]interface{}) {
	result = map[component.Type]interface{}{}
	resp, err := c.client.ComponentOptions(c.Ctx, &empty.Empty{})
	if err != nil {
		c.Logger.Error("unexpected error when requesting component options",
			"error", err)
		return
	}
	for t, optsProto := range resp.Options {
		typ := component.Type(t)
		opts, err := component.UnmarshalOptionsProto(typ, optsProto)
		if err != nil {
			c.Logger.Error("cannot unmarshal options for type",
				"type", typ, "options", opts)
			continue
		}
		result[typ] = opts
	}
	return
}

func (c *pluginInfoClient) ComponentTypes() (result []component.Type) {
	result = []component.Type{}
	resp, err := c.client.ComponentTypes(c.Ctx, &empty.Empty{})
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

func (c *pluginInfoClient) Name() string {
	resp, err := c.client.Name(c.Ctx, &empty.Empty{})
	if err != nil {
		c.Logger.Error("unexpected error when requesting component name",
			"error", err)

		return ""
	}
	return resp.Name
}

func (s *pluginInfoServer) ComponentOptions(
	ctx context.Context,
	_ *empty.Empty,
) (*vagrant_plugin_sdk.PluginInfo_ComponentOptionsMap, error) {
	if err := isImplemented(s, "plugin info"); err != nil {
		return nil, err
	}

	optsMap := s.Impl.ComponentOptions()
	result := &vagrant_plugin_sdk.PluginInfo_ComponentOptionsMap{
		Options: map[uint32]*anypb.Any{},
	}
	for t, opts := range optsMap {
		s.Logger.Info("trying to map these opts to a proto",
			"type", component.Type(t).String(), "opts", opts)
		any, err := component.ProtoAny(opts)
		if err != nil {
			s.Logger.Error("unexpected error while encoding component options into any",
				"error", err)
			return nil, err
		}
		if any != nil {
			result.Options[uint32(t)] = any
		}
	}

	return result, nil
}

func (s *pluginInfoServer) ComponentTypes(
	ctx context.Context,
	_ *empty.Empty,
) (*vagrant_plugin_sdk.PluginInfo_ComponentList, error) {
	if err := isImplemented(s, "plugin info"); err != nil {
		return nil, err
	}

	list := []uint32{}
	for _, t := range s.Impl.ComponentTypes() {
		list = append(list, uint32(t))
	}
	return &vagrant_plugin_sdk.PluginInfo_ComponentList{
		Component: list,
	}, nil
}

func (s *pluginInfoServer) Name(
	ctx context.Context,
	_ *empty.Empty,
) (*vagrant_plugin_sdk.PluginInfo_Name, error) {
	if err := isImplemented(s, "plugin info"); err != nil {
		return nil, err
	}

	return &vagrant_plugin_sdk.PluginInfo_Name{
		Name: s.Impl.Name(),
	}, nil
}
