package core

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/hashicorp/go-argmapper"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"

	"github.com/hashicorp/vagrant-plugin-sdk/core"
	"github.com/hashicorp/vagrant-plugin-sdk/internal/pluginargs"
	"github.com/hashicorp/vagrant-plugin-sdk/proto/vagrant_plugin_sdk"
)

type TargetIndexPlugin struct {
	plugin.NetRPCUnsupportedPlugin

	Impl    core.TargetIndex
	Mappers []*argmapper.Func
	Logger  hclog.Logger
}

func (p *TargetIndexPlugin) GRPCClient(
	ctx context.Context,
	broker *plugin.GRPCBroker,
	c *grpc.ClientConn,
) (interface{}, error) {
	return &targetIndexClient{
		client: vagrant_plugin_sdk.NewTargetIndexServiceClient(c),
		ctx:    ctx,
		base: &base{
			Mappers: p.Mappers,
			Logger:  p.Logger,
			Broker:  broker,
			Cleanup: &pluginargs.Cleanup{},
		},
	}, nil
}

func (p *TargetIndexPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	vagrant_plugin_sdk.RegisterTargetIndexServiceServer(s, &targetIndexServer{
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

type targetIndexClient struct {
	*base

	ctx    context.Context
	client vagrant_plugin_sdk.TargetIndexServiceClient
}

type targetIndexServer struct {
	*base

	Impl core.TargetIndex
	vagrant_plugin_sdk.UnimplementedTargetIndexServiceServer
}

func (t *targetIndexClient) Delete(machine core.Machine) (err error) {
	id, err := machine.ID()
	_, err = t.client.Delete(t.ctx, &vagrant_plugin_sdk.Ref_Target{
		ResourceId: id,
	})
	return
}

func (t *targetIndexClient) Get(uuid string) (entry core.Machine, err error) {
	return
}

func (t *targetIndexClient) Includes(uuid string) (exists bool, err error) {
	return
}

func (t *targetIndexClient) Set(entry core.Machine) (updatedEntry core.Machine, err error) {
	return
}

func (t *targetIndexClient) Recover(entry core.Machine) (updatedEntry core.Machine, err error) {
	return
}

func (t *targetIndexServer) Delete(
	ctx context.Context,
	in *vagrant_plugin_sdk.Ref_Target,
) (empty *empty.Empty, err error) {
	return
}

func (t *targetIndexServer) Get(
	ctx context.Context,
	in *vagrant_plugin_sdk.TargetIndex_GetRequest,
) (target *vagrant_plugin_sdk.Args_Target, err error) {
	return
}

func (t *targetIndexServer) Includes(
	ctx context.Context,
	in *vagrant_plugin_sdk.TargetIndex_IncludesRequest,
) (target *vagrant_plugin_sdk.Args_Target, err error) {
	return
}

func (t *targetIndexServer) Set(
	ctx context.Context,
	in *vagrant_plugin_sdk.Args_Target,
) (target *vagrant_plugin_sdk.Args_Target, err error) {
	return
}

func (t *targetIndexServer) Recover(
	ctx context.Context,
	in *vagrant_plugin_sdk.Args_Target,
) (target *vagrant_plugin_sdk.Args_Target, err error) {
	return
}

var (
	_ plugin.Plugin     = (*TargetIndexPlugin)(nil)
	_ plugin.GRPCPlugin = (*TargetIndexPlugin)(nil)
	_ core.TargetIndex  = (*targetIndexClient)(nil)
)
