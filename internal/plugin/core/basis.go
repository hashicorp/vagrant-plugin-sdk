package core

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/hashicorp/go-argmapper"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/hashicorp/vagrant-plugin-sdk/core"
	"github.com/hashicorp/vagrant-plugin-sdk/datadir"
	"github.com/hashicorp/vagrant-plugin-sdk/internal/pluginargs"
	"github.com/hashicorp/vagrant-plugin-sdk/proto/vagrant_plugin_sdk"
	"github.com/hashicorp/vagrant-plugin-sdk/terminal"
)

type BasisPlugin struct {
	plugin.NetRPCUnsupportedPlugin

	Impl    core.Basis
	Mappers []*argmapper.Func
	Logger  hclog.Logger
}

func (p *BasisPlugin) GRPCClient(
	ctx context.Context,
	broker *plugin.GRPCBroker,
	c *grpc.ClientConn,
) (interface{}, error) {
	return &basisClient{
		client: vagrant_plugin_sdk.NewBasisServiceClient(c),
		ctx:    ctx,
		base: &base{
			Mappers: p.Mappers,
			Logger:  p.Logger,
			Broker:  broker,
			Cleanup: &pluginargs.Cleanup{},
		},
	}, nil
}

func (p *BasisPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	vagrant_plugin_sdk.RegisterBasisServiceServer(s, &basisServer{
		Impl: p.Impl,
		base: &base{
			Mappers: p.Mappers,
			Logger:  p.Logger,
			Broker:  broker,
			Cleanup: &pluginargs.Cleanup{},
		},
	})
	return nil
}

type basisClient struct {
	*base

	ctx    context.Context
	client vagrant_plugin_sdk.BasisServiceClient
}

type basisServer struct {
	*base

	Impl core.Basis
	vagrant_plugin_sdk.UnimplementedBasisServiceServer
}

func (p *basisClient) UI() (ui terminal.UI, err error) {
	r, err := p.client.UI(p.ctx, &emptypb.Empty{})
	if err != nil {
		return
	}

	result, err := p.Map(r, (*terminal.UI)(nil))
	if err == nil {
		ui = result.(terminal.UI)
	}

	return
}

func (p *basisClient) DataDir() (dir *datadir.Basis, err error) {
	r, err := p.client.DataDir(p.ctx, &emptypb.Empty{})
	if err != nil {
		return
	}

	result, err := p.Map(r, (**datadir.Basis)(nil))
	if err == nil {
		dir = result.(*datadir.Basis)
	}

	return
}

func (p *basisServer) DataDir(
	ctx context.Context,
	_ *empty.Empty,
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

var (
	_ plugin.Plugin     = (*BasisPlugin)(nil)
	_ plugin.GRPCPlugin = (*BasisPlugin)(nil)
	_ core.Basis        = (*basisClient)(nil)
)
