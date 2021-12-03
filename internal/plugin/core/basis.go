package core

import (
	"context"

	"github.com/hashicorp/go-argmapper"
	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/hashicorp/vagrant-plugin-sdk/component"
	"github.com/hashicorp/vagrant-plugin-sdk/core"
	"github.com/hashicorp/vagrant-plugin-sdk/datadir"
	"github.com/hashicorp/vagrant-plugin-sdk/internal-shared/dynamic"
	vplugin "github.com/hashicorp/vagrant-plugin-sdk/internal/plugin"
	"github.com/hashicorp/vagrant-plugin-sdk/proto/vagrant_plugin_sdk"
	"github.com/hashicorp/vagrant-plugin-sdk/terminal"
)

type BasisPlugin struct {
	plugin.NetRPCUnsupportedPlugin

	Impl core.Basis
	*vplugin.BasePlugin
}

func (p *BasisPlugin) GRPCClient(
	ctx context.Context,
	broker *plugin.GRPCBroker,
	c *grpc.ClientConn,
) (interface{}, error) {
	return &basisClient{
		client:     vagrant_plugin_sdk.NewBasisServiceClient(c),
		BaseClient: p.NewClient(ctx, broker),
	}, nil
}

func (p *BasisPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	vagrant_plugin_sdk.RegisterBasisServiceServer(s, &basisServer{
		Impl:       p.Impl,
		BaseServer: p.NewServer(broker),
	})
	return nil
}

type basisClient struct {
	*vplugin.BaseClient

	client vagrant_plugin_sdk.BasisServiceClient
}

type basisServer struct {
	*vplugin.BaseServer

	Impl core.Basis
	vagrant_plugin_sdk.UnimplementedBasisServiceServer
}

func (p *basisClient) UI() (ui terminal.UI, err error) {
	r, err := p.client.UI(p.Ctx, &emptypb.Empty{})
	if err != nil {
		return
	}

	result, err := p.Map(r, (*terminal.UI)(nil), argmapper.Typed(p.Ctx))
	if err == nil {
		ui = result.(terminal.UI)
	}

	return
}

func (p *basisClient) DataDir() (dir *datadir.Basis, err error) {
	r, err := p.client.DataDir(p.Ctx, &emptypb.Empty{})
	if err != nil {
		return
	}

	result, err := p.Map(r, (**datadir.Basis)(nil))
	if err == nil {
		dir = result.(*datadir.Basis)
	}

	return
}

func (p *basisClient) Host() (h core.Host, err error) {
	r, err := p.client.Host(p.Ctx, &emptypb.Empty{})
	if err != nil {
		return
	}

	result, err := p.Map(r, (*core.Host)(nil), argmapper.Typed(p.Ctx))
	if err == nil {
		h = result.(core.Host)
	}

	return
}

func (p *basisClient) Plugins(types ...string) (h []*core.NamedPlugin, err error) {
	r, err := p.client.Plugins(p.Ctx, &vagrant_plugin_sdk.Basis_PluginsRequest{
		Types: types,
	})
	if err != nil {
		return
	}

	result := []*core.NamedPlugin{}
	for _, plugin := range r.Plugins {
		typ, err := component.FindType(plugin.Type)
		if err != nil {
			return nil, err
		}
		plg, err := p.Map(r, typ, argmapper.Typed(p.Ctx))
		if err != nil {
			return nil, err
		}
		result = append(result, &core.NamedPlugin{
			Name:   plugin.Name,
			Plugin: plg,
		})
	}

	return result, nil
}

func (p *basisServer) DataDir(
	ctx context.Context,
	_ *emptypb.Empty,
) (r *vagrant_plugin_sdk.Args_DataDir_Basis, err error) {
	d, err := p.Impl.DataDir()
	if err != nil {
		return
	}

	result, err := p.Map(d, (**vagrant_plugin_sdk.Args_DataDir_Basis)(nil))
	if err == nil {
		r = result.(*vagrant_plugin_sdk.Args_DataDir_Basis)
	}

	return
}

func (p *basisServer) Host(
	ctx context.Context,
	_ *emptypb.Empty,
) (r *vagrant_plugin_sdk.Args_Host, err error) {
	d, err := p.Impl.Host()
	if err != nil {
		return
	}

	result, err := p.Map(d, (**vagrant_plugin_sdk.Args_Host)(nil),
		argmapper.Typed(ctx))
	if err == nil {
		r = result.(*vagrant_plugin_sdk.Args_Host)
	}

	return
}

func (t *basisServer) UI(
	ctx context.Context,
	_ *emptypb.Empty,
) (r *vagrant_plugin_sdk.Args_TerminalUI, err error) {
	d, err := t.Impl.UI()
	if err != nil {
		return
	}

	result, err := t.Map(d, (**vagrant_plugin_sdk.Args_TerminalUI)(nil),
		argmapper.Typed(ctx))
	if err == nil {
		r = result.(*vagrant_plugin_sdk.Args_TerminalUI)
	}

	return
}

func (t *basisServer) Plugins(
	ctx context.Context,
	in *vagrant_plugin_sdk.Basis_PluginsRequest,
) (r *vagrant_plugin_sdk.Basis_PluginsResponse, err error) {
	plugins, err := t.Impl.Plugins(in.Types...)
	if err != nil {
		return
	}

	result := []*vagrant_plugin_sdk.Basis_Plugin{}
	for _, plugin := range plugins {
		val, err := dynamic.UnknownMap(plugin.Plugin, (*proto.Message)(nil), t.Mappers,
			argmapper.Typed(t.Internal()),
			argmapper.Typed(ctx),
			argmapper.Typed(t.Logger),
		)
		if err != nil {
			return nil, err
		}
		pluginProto, err := dynamic.EncodeAny(val.(proto.Message))
		if err != nil {
			return nil, err
		}
		result = append(result, &vagrant_plugin_sdk.Basis_Plugin{
			Name:   plugin.Name,
			Type:   plugin.Type,
			Plugin: pluginProto,
		})
	}

	return &vagrant_plugin_sdk.Basis_PluginsResponse{
		Plugins: result,
	}, nil
}

var (
	_ plugin.Plugin     = (*BasisPlugin)(nil)
	_ plugin.GRPCPlugin = (*BasisPlugin)(nil)
	_ core.Basis        = (*basisClient)(nil)
)
