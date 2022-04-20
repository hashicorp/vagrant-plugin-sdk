package core

import (
	"context"

	"github.com/hashicorp/go-argmapper"
	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/hashicorp/vagrant-plugin-sdk/core"
	"github.com/hashicorp/vagrant-plugin-sdk/datadir"
	"github.com/hashicorp/vagrant-plugin-sdk/helper/path"
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
	defer func() {
		if err != nil {
			p.Logger.Error("failed to get boxes",
				"error", err,
			)
		}
	}()

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

func (b *basisClient) CWD() (p path.Path, err error) {
	defer func() {
		if err != nil {
			b.Logger.Error("failed to get current working directory",
				"error", err,
			)
		}
	}()

	r, err := b.client.CWD(b.Ctx, &emptypb.Empty{})
	if err == nil {
		p = path.NewPath(r.Path)
	}

	return
}

func (b *basisClient) DefaultPrivateKey() (p path.Path, err error) {
	defer func() {
		if err != nil {
			b.Logger.Error("failed to get default private key",
				"error", err,
			)
		}
	}()

	r, err := b.client.DefaultPrivateKey(b.Ctx, &emptypb.Empty{})
	if err == nil {
		p = path.NewPath(r.Path)
	}

	return
}

func (p *basisClient) DataDir() (dir *datadir.Basis, err error) {
	defer func() {
		if err != nil {
			p.Logger.Error("failed to get data directory",
				"error", err,
			)
		}
	}()

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
	defer func() {
		if err != nil {
			p.Logger.Error("failed to get host",
				"error", err,
			)
		}
	}()

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

func (p *basisClient) ResourceId() (rid string, err error) {
	defer func() {
		if err != nil {
			p.Logger.Error("failed to get resource id",
				"error", err,
			)
		}
	}()

	r, err := p.client.ResourceId(p.Ctx, &emptypb.Empty{})
	if err == nil {
		rid = r.ResourceId
	}

	return
}

func (p *basisClient) TargetIndex() (index core.TargetIndex, err error) {
	defer func() {
		if err != nil {
			p.Logger.Error("failed to get target index",
				"error", err,
			)
		}
	}()

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
	defer func() {
		if err != nil {
			p.Logger.Error("failed to get ui",
				"error", err,
			)
		}
	}()

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
	defer func() {
		if err != nil {
			p.Logger.Error("failed to get boxes",
				"error", err,
			)
		}
	}()

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
		p.Logger.Error("failed to get current working directory",
			"error", err,
		)

		return nil, err
	}

	return &vagrant_plugin_sdk.Args_Path{
		Path: c.String(),
	}, nil
}

func (p *basisServer) DataDir(
	ctx context.Context,
	_ *emptypb.Empty,
) (r *vagrant_plugin_sdk.Args_DataDir_Basis, err error) {
	defer func() {
		if err != nil {
			p.Logger.Error("failed to get data directory",
				"error", err,
			)
		}
	}()

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
		p.Logger.Error("failed to get default private key",
			"error", err,
		)

		return nil, err
	}

	return &vagrant_plugin_sdk.Args_Path{
		Path: c.String(),
	}, nil
}

func (p *basisServer) Host(
	ctx context.Context,
	_ *emptypb.Empty,
) (r *vagrant_plugin_sdk.Args_Host, err error) {
	defer func() {
		if err != nil {
			p.Logger.Error("failed to get host",
				"error", err,
			)
		}
	}()

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

func (p *basisServer) ResourceId(
	ctx context.Context,
	_ *emptypb.Empty,
) (*vagrant_plugin_sdk.Basis_ResourceIdResponse, error) {
	rid, err := p.Impl.ResourceId()

	if err != nil {
		p.Logger.Error("resource id lookup failed",
			"error", err,
		)

		return nil, err
	}

	return &vagrant_plugin_sdk.Basis_ResourceIdResponse{
		ResourceId: rid,
	}, nil
}

func (p *basisServer) TargetIndex(
	ctx context.Context,
	_ *emptypb.Empty,
) (r *vagrant_plugin_sdk.Args_TargetIndex, err error) {
	defer func() {
		if err != nil {
			p.Logger.Error("failed to get target index",
				"error", err,
			)
		}
	}()

	idx, err := p.Impl.TargetIndex()
	if err != nil {
		return
	}

	result, err := p.Map(idx, (**vagrant_plugin_sdk.Args_TargetIndex)(nil),
		argmapper.Typed(ctx))
	if err == nil {
		r = result.(*vagrant_plugin_sdk.Args_TargetIndex)
	}
	return
}

func (p *basisServer) UI(
	ctx context.Context,
	_ *emptypb.Empty,
) (r *vagrant_plugin_sdk.Args_TerminalUI, err error) {
	defer func() {
		if err != nil {
			p.Logger.Error("failed to get ui",
				"error", err,
			)
		}
	}()

	d, err := p.Impl.UI()
	if err != nil {
		return
	}

	result, err := p.Map(d, (**vagrant_plugin_sdk.Args_TerminalUI)(nil),
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
