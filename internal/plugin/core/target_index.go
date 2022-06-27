package core

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/hashicorp/go-argmapper"
	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hashicorp/vagrant-plugin-sdk/core"
	vplugin "github.com/hashicorp/vagrant-plugin-sdk/internal/plugin"
	"github.com/hashicorp/vagrant-plugin-sdk/proto/vagrant_plugin_sdk"
)

type TargetIndexPlugin struct {
	plugin.NetRPCUnsupportedPlugin

	Impl core.TargetIndex
	*vplugin.BasePlugin
}

func (p *TargetIndexPlugin) GRPCClient(
	ctx context.Context,
	broker *plugin.GRPCBroker,
	c *grpc.ClientConn,
) (interface{}, error) {
	return &targetIndexClient{
		client:     vagrant_plugin_sdk.NewTargetIndexServiceClient(c),
		BaseClient: p.NewClient(ctx, broker, nil),
	}, nil
}

func (p *TargetIndexPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	vagrant_plugin_sdk.RegisterTargetIndexServiceServer(s, &targetIndexServer{
		Impl:       p.Impl,
		BaseServer: p.NewServer(broker, nil),
	})
	return nil
}

type targetIndexClient struct {
	*vplugin.BaseClient

	client vagrant_plugin_sdk.TargetIndexServiceClient
}

type targetIndexServer struct {
	*vplugin.BaseServer

	Impl core.TargetIndex
	vagrant_plugin_sdk.UnimplementedTargetIndexServiceServer
}

func (t *targetIndexClient) Delete(uuid string) (err error) {
	_, err = t.client.Delete(
		t.Ctx,
		&vagrant_plugin_sdk.TargetIndex_TargetIdentifier{Id: uuid},
	)
	return
}

func (t *targetIndexClient) Get(uuid string) (entry core.Target, err error) {
	target, err := t.client.Get(
		t.Ctx,
		&vagrant_plugin_sdk.TargetIndex_TargetIdentifier{Id: uuid},
	)
	if err != nil {
		return nil, err
	}
	m, err := t.Map(
		target,
		(*core.Target)(nil),
		argmapper.Typed(t.Ctx),
	)
	return m.(core.Target), err
}

func (t *targetIndexClient) Includes(uuid string) (exists bool, err error) {
	incl, err := t.client.Includes(
		t.Ctx,
		&vagrant_plugin_sdk.TargetIndex_TargetIdentifier{Id: uuid},
	)
	return incl.Exists, err
}

func (t *targetIndexClient) Set(entry core.Target) (updatedEntry core.Target, err error) {
	targetArg, err := t.Map(
		entry,
		(*vagrant_plugin_sdk.Args_Target)(nil),
		argmapper.Typed(t.Ctx),
	)
	if err != nil {
		return nil, err
	}
	targetOut, err := t.client.Set(t.Ctx, targetArg.(*vagrant_plugin_sdk.Args_Target))
	if err != nil {
		return nil, err
	}
	m, err := t.Map(
		targetOut,
		(*core.Target)(nil),
		argmapper.Typed(t.Ctx),
	)
	return m.(core.Target), err
}

func (t *targetIndexClient) All() (targets []core.Target, err error) {
	argTargets, err := t.client.All(t.Ctx, &empty.Empty{})
	if err != nil {
		return nil, err
	}
	targets = []core.Target{}
	for _, argTarget := range argTargets.Targets {
		result, err := t.Map(
			argTarget,
			(*core.Target)(nil),
			argmapper.Typed(t.Ctx),
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
	defer func() {
		if err != nil {
			t.Logger.Error("failed to delete target from index",
				"error", err,
			)
		}
	}()

	err = t.Impl.Delete(in.Id)
	return
}

func (t *targetIndexServer) Get(
	ctx context.Context,
	in *vagrant_plugin_sdk.TargetIndex_TargetIdentifier,
) (target *vagrant_plugin_sdk.Args_Target, err error) {
	defer func() {
		// Log errors, but not NotFound errors which happen during normal operations
		if err != nil && (status.Convert(err).Code() != codes.NotFound) {
			t.Logger.Error("failed to get target from index",
				"error", err,
			)
		}
	}()

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
	defer func() {
		if err != nil {
			t.Logger.Error("failed to check inclusion in index",
				"error", err,
			)
		}
	}()

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
	defer func() {
		if err != nil {
			t.Logger.Error("failed to set target index",
				"error", err,
			)
		}
	}()

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
	defer func() {
		if err != nil {
			t.Logger.Error("failed to get target index list",
				"error", err,
			)
		}
	}()

	targets, err := t.Impl.All()
	argsTargets := []*vagrant_plugin_sdk.Args_Target{}
	for _, target := range targets {
		var result interface{}
		result, err = t.Map(target, (**vagrant_plugin_sdk.Args_Target)(nil))
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
