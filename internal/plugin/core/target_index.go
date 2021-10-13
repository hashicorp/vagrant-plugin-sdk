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
			Logger:  p.Logger.Named("core.target-index"),
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
			Logger:  p.Logger.Named("core.target-index"),
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

func (t *targetIndexClient) Delete(uuid string) (err error) {
	_, err = t.client.Delete(
		t.ctx,
		&vagrant_plugin_sdk.TargetIndex_TargetIdentifier{Id: uuid},
	)
	return
}

func (t *targetIndexClient) Get(uuid string) (entry core.Target, err error) {
	target, err := t.client.Get(
		t.ctx,
		&vagrant_plugin_sdk.TargetIndex_TargetIdentifier{Id: uuid},
	)
	if err != nil {
		return nil, err
	}
	m, err := t.Map(
		target,
		(*core.Target)(nil),
		argmapper.Typed(t.ctx),
	)
	return m.(core.Target), err
}

func (t *targetIndexClient) Includes(uuid string) (exists bool, err error) {
	incl, err := t.client.Includes(
		t.ctx,
		&vagrant_plugin_sdk.TargetIndex_TargetIdentifier{Id: uuid},
	)
	return incl.Exists, err
}

func (t *targetIndexClient) Set(entry core.Target) (updatedEntry core.Target, err error) {
	targetArg, err := t.Map(
		entry,
		(*vagrant_plugin_sdk.Args_Target)(nil),
		argmapper.Typed(t.ctx),
	)
	if err != nil {
		return nil, err
	}
	targetOut, err := t.client.Set(t.ctx, targetArg.(*vagrant_plugin_sdk.Args_Target))
	if err != nil {
		return nil, err
	}
	m, err := t.Map(
		targetOut,
		(*core.Target)(nil),
		argmapper.Typed(t.ctx),
	)
	return m.(core.Target), err
}

func (t *targetIndexClient) All() (targets []core.Target, err error) {
	argTargets, err := t.client.All(t.ctx, &empty.Empty{})
	if err != nil {
		return nil, err
	}
	targets = []core.Target{}
	for _, argTarget := range argTargets.Targets {
		result, err := t.Map(
			argTarget,
			(*core.Target)(nil),
			argmapper.Typed(t.ctx),
		)
		if err != nil {
			return nil, err
		}
		targets = append(targets, result.(core.Target))
	}
	return
}

// Target Index Server

func (t *targetIndexServer) Delete(
	ctx context.Context,
	in *vagrant_plugin_sdk.TargetIndex_TargetIdentifier,
) (empty *empty.Empty, err error) {
	err = t.Impl.Delete(in.Id)
	return
}

func (t *targetIndexServer) Get(
	ctx context.Context,
	in *vagrant_plugin_sdk.TargetIndex_TargetIdentifier,
) (target *vagrant_plugin_sdk.Args_Target, err error) {
	tar, err := t.Impl.Get(in.Id)
	if err != nil {
		return nil, err
	}
	result, err := t.Map(tar, (**vagrant_plugin_sdk.Args_Target)(nil))
	if err != nil {
		return nil, err
	}

	return result.(*vagrant_plugin_sdk.Args_Target), err
}

func (t *targetIndexServer) Includes(
	ctx context.Context,
	in *vagrant_plugin_sdk.TargetIndex_TargetIdentifier,
) (result *vagrant_plugin_sdk.TargetIndex_IncludesResponse, err error) {
	resp, err := t.Impl.Includes(in.Id)
	if err != nil {
		return nil, err
	}
	result = &vagrant_plugin_sdk.TargetIndex_IncludesResponse{
		Exists: resp,
	}
	return
}

func (t *targetIndexServer) Set(
	ctx context.Context,
	in *vagrant_plugin_sdk.Args_Target,
) (target *vagrant_plugin_sdk.Args_Target, err error) {
	targetIn, err := t.Map(in, (*core.Target)(nil),
		argmapper.Typed(ctx))

	targetOut, err := t.Impl.Set(targetIn.(core.Target))
	if err != nil {
		return nil, err
	}
	result, err := t.Map(targetOut, (**vagrant_plugin_sdk.Args_Target)(nil))
	if err != nil {
		return nil, err
	}

	return result.(*vagrant_plugin_sdk.Args_Target), nil
}

func (t *targetIndexServer) All(
	ctx context.Context,
	_ *empty.Empty,
) (resp *vagrant_plugin_sdk.TargetIndex_AllResponse, err error) {
	targets, err := t.Impl.All()
	argsTargets := []*vagrant_plugin_sdk.Args_Target{}
	for _, target := range targets {
		result, err := t.Map(target, (**vagrant_plugin_sdk.Args_Target)(nil))
		if err != nil {
			return nil, err
		}
		argsTargets = append(argsTargets, result.(*vagrant_plugin_sdk.Args_Target))
	}
	resp = &vagrant_plugin_sdk.TargetIndex_AllResponse{
		Targets: argsTargets,
	}
	return
}

var (
	_ plugin.Plugin     = (*TargetIndexPlugin)(nil)
	_ plugin.GRPCPlugin = (*TargetIndexPlugin)(nil)
	_ core.TargetIndex  = (*targetIndexClient)(nil)
)
