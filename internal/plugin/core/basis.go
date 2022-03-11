package core

import (
	"context"

	"github.com/hashicorp/go-argmapper"
	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/hashicorp/vagrant-plugin-sdk/core"
	"github.com/hashicorp/vagrant-plugin-sdk/datadir"
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
		BaseClient: p.NewClient(ctx, broker, nil),
	}, nil
}

func (p *BasisPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	vagrant_plugin_sdk.RegisterBasisServiceServer(s, &basisServer{
		Impl:       p.Impl,
		BaseServer: p.NewServer(broker, nil),
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
	vagrant_plugin_sdk.UnsafeBasisServiceServer
}

func (p *basisClient) Boxes() (b core.BoxCollection, err error) {
	r, err := p.client.Boxes(p.Ctx, &emptypb.Empty{})
	if err != nil {
		return
	}

	result, err := p.Map(r, (*core.BoxCollection)(nil),
		argmapper.Typed(p.Ctx),
	)
	if err == nil {
		b = result.(core.BoxCollection)
	}

	return
}

func (p *basisClient) CWD() (path string, err error) {
	r, err := p.client.CWD(p.Ctx, &emptypb.Empty{})
	if err == nil {
		path = r.Path
	}

	return
}

func (p *basisClient) DefaultPrivateKey() (path string, err error) {
	r, err := p.client.DefaultPrivateKey(p.Ctx, &emptypb.Empty{})
	if err == nil {
		path = r.Path
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

func (p *basisClient) TargetIndex() (index core.TargetIndex, err error) {
	r, err := p.client.TargetIndex(p.Ctx, &emptypb.Empty{})
	if err != nil {
		return
	}

	result, err := p.Map(r, (*core.TargetIndex)(nil),
		argmapper.Typed(p.Ctx))
	if err == nil {
		index = result.(core.TargetIndex)
	}
	return
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

func (p *basisServer) Boxes(
	ctx context.Context,
	_ *emptypb.Empty,
) (r *vagrant_plugin_sdk.Args_BoxCollection, err error) {
	boxCollection, err := p.Impl.Boxes()
	if err != nil {
		return
	}

	result, err := p.Map(boxCollection, (**vagrant_plugin_sdk.Args_BoxCollection)(nil))
	if err == nil {
		r = result.(*vagrant_plugin_sdk.Args_BoxCollection)
	}

	return
}

func (p *basisServer) CWD(
	ctx context.Context,
	_ *emptypb.Empty,
) (*vagrant_plugin_sdk.Args_Path, error) {
	c, err := p.Impl.CWD()
	if err != nil {
		return nil, err
	}

	return &vagrant_plugin_sdk.Args_Path{
		Path: c,
	}, nil
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

func (p *basisServer) DefaultPrivateKey(
	ctx context.Context,
	_ *emptypb.Empty,
) (*vagrant_plugin_sdk.Args_Path, error) {
	c, err := p.Impl.DefaultPrivateKey()
	if err != nil {
		return nil, err
	}

	return &vagrant_plugin_sdk.Args_Path{
		Path: c,
	}, nil
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

func (p *basisServer) TargetIndex(
	ctx context.Context,
	_ *emptypb.Empty,
) (r *vagrant_plugin_sdk.Args_TargetIndex, err error) {
	idx, err := p.Impl.TargetIndex()
	if err != nil {
		return nil, err
	}

	result, err := p.Map(idx, (**vagrant_plugin_sdk.Args_TargetIndex)(nil),
		argmapper.Typed(ctx))
	if err == nil {
		r = result.(*vagrant_plugin_sdk.Args_TargetIndex)
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

var (
	_ plugin.Plugin                         = (*BasisPlugin)(nil)
	_ plugin.GRPCPlugin                     = (*BasisPlugin)(nil)
	_ core.Basis                            = (*basisClient)(nil)
	_ vagrant_plugin_sdk.BasisServiceServer = (*basisServer)(nil)
	_ core.Seeder                           = (*basisClient)(nil)
)
