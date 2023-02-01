// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package plugin

import (
	"context"

	"github.com/LK4D4/joincontext"
	"github.com/hashicorp/go-argmapper"
	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/hashicorp/vagrant-plugin-sdk/component"
	"github.com/hashicorp/vagrant-plugin-sdk/core"
	"github.com/hashicorp/vagrant-plugin-sdk/internal/funcspec"
	"github.com/hashicorp/vagrant-plugin-sdk/proto/vagrant_plugin_sdk"
)

// SyncedFolderPlugin implements plugin.Plugin (specifically GRPCPlugin) for
// the SyncedFolder component type.
type SyncedFolderPlugin struct {
	plugin.NetRPCUnsupportedPlugin

	Impl component.SyncedFolder // Impl is the concrete implementation
	*BasePlugin
}

func (p *SyncedFolderPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	bs := p.NewServer(broker, p.Impl)
	vagrant_plugin_sdk.RegisterSyncedFolderServiceServer(s, &syncedFolderServer{
		Impl:       p.Impl,
		BaseServer: bs,
		capabilityServer: &capabilityServer{
			BaseServer:     bs,
			CapabilityImpl: p.Impl,
			typ:            "synced_folder",
		},
	})
	return nil
}

func (p *SyncedFolderPlugin) GRPCClient(
	ctx context.Context,
	broker *plugin.GRPCBroker,
	c *grpc.ClientConn,
) (interface{}, error) {
	client := vagrant_plugin_sdk.NewSyncedFolderServiceClient(c)
	bc := p.NewClient(ctx, broker, client.(SeederClient))

	return &syncedFolderClient{
		BaseClient: bc,
		client:     client,
		capabilityClient: &capabilityClient{
			client:     client,
			BaseClient: bc,
		},
	}, nil
}

// syncedFolderClient is an implementation of component.SyncedFolder over gRPC.
type syncedFolderClient struct {
	*BaseClient
	*capabilityClient
	client vagrant_plugin_sdk.SyncedFolderServiceClient
}

func (c *syncedFolderClient) UsableFunc() interface{} {
	spec, err := c.client.UsableSpec(c.Ctx, &emptypb.Empty{})
	if err != nil {
		return funcErr(err)
	}
	spec.Result = nil
	cb := func(ctx context.Context, args funcspec.Args) (bool, error) {
		ctx, _ = joincontext.Join(c.Ctx, ctx)
		resp, err := c.client.Usable(ctx, &vagrant_plugin_sdk.FuncSpec_Args{Args: args})
		if err != nil {
			return false, err
		}
		return resp.Usable, nil
	}

	return c.GenerateFunc(spec, cb)
}

func (c *syncedFolderClient) Usable(machine core.Machine) (bool, error) {
	f := c.UsableFunc()
	raw, err := c.CallDynamicFunc(f, (*bool)(nil),
		argmapper.Typed(c.Ctx),
		argmapper.Typed(machine),
	)
	if err != nil {
		return false, err
	}

	return raw.(bool), nil
}

func (c *syncedFolderClient) EnableFunc() interface{} {
	spec, err := c.client.EnableSpec(c.Ctx, &emptypb.Empty{})
	if err != nil {
		return funcErr(err)
	}
	spec.Result = nil
	cb := func(ctx context.Context, args funcspec.Args) error {
		ctx, _ = joincontext.Join(c.Ctx, ctx)
		_, err := c.client.Enable(ctx, &vagrant_plugin_sdk.FuncSpec_Args{Args: args})
		return err
	}

	return c.GenerateFunc(spec, cb)
}

func (c *syncedFolderClient) Enable(machine core.Machine, folders []*core.Folder, opts ...interface{}) error {
	f := c.EnableFunc()
	_, err := c.CallDynamicFunc(f, false,
		argmapper.Typed(c.Ctx),
		argmapper.Typed(machine),
		argmapper.Typed(folders),
		argmapper.Typed(argmapper.Typed(&component.Direct{Arguments: opts})),
	)
	return err
}

func (c *syncedFolderClient) PrepareFunc() interface{} {
	spec, err := c.client.PrepareSpec(c.Ctx, &emptypb.Empty{})
	if err != nil {
		return funcErr(err)
	}
	spec.Result = nil
	cb := func(ctx context.Context, args funcspec.Args) error {
		ctx, _ = joincontext.Join(c.Ctx, ctx)
		_, err := c.client.Prepare(ctx, &vagrant_plugin_sdk.FuncSpec_Args{Args: args})
		return err
	}

	return c.GenerateFunc(spec, cb)
}

func (c *syncedFolderClient) Prepare(machine core.Machine, folders []*core.Folder, opts ...interface{}) error {
	f := c.PrepareFunc()
	_, err := c.CallDynamicFunc(f, false,
		argmapper.Typed(c.Ctx),
		argmapper.Typed(machine),
		argmapper.Typed(folders),
		argmapper.Typed(argmapper.Typed(&component.Direct{Arguments: opts})),
	)
	return err
}

func (c *syncedFolderClient) DisableFunc() interface{} {
	spec, err := c.client.DisableSpec(c.Ctx, &emptypb.Empty{})
	if err != nil {
		return funcErr(err)
	}
	spec.Result = nil
	cb := func(ctx context.Context, args funcspec.Args) error {
		ctx, _ = joincontext.Join(c.Ctx, ctx)
		_, err := c.client.Disable(ctx, &vagrant_plugin_sdk.FuncSpec_Args{Args: args})
		return err
	}

	return c.GenerateFunc(spec, cb)
}

func (c *syncedFolderClient) Disable(machine core.Machine, folders []*core.Folder, opts ...interface{}) error {
	f := c.DisableFunc()
	_, err := c.CallDynamicFunc(f, false,
		argmapper.Typed(c.Ctx),
		argmapper.Typed(machine),
		argmapper.Typed(folders),
		argmapper.Typed(argmapper.Typed(&component.Direct{Arguments: opts})),
	)
	return err
}

func (c *syncedFolderClient) CleanupFunc() interface{} {
	spec, err := c.client.CleanupSpec(c.Ctx, &emptypb.Empty{})
	if err != nil {
		return funcErr(err)
	}
	spec.Result = nil
	cb := func(ctx context.Context, args funcspec.Args) error {
		ctx, _ = joincontext.Join(c.Ctx, ctx)
		_, err := c.client.Cleanup(ctx, &vagrant_plugin_sdk.FuncSpec_Args{Args: args})
		return err
	}

	return c.GenerateFunc(spec, cb)
}

func (c *syncedFolderClient) Cleanup(machine core.Machine, opts ...interface{}) error {
	f := c.CleanupFunc()
	_, err := c.CallDynamicFunc(f, false,
		argmapper.Typed(c.Ctx),
		argmapper.Typed(machine),
		argmapper.Typed(argmapper.Typed(&component.Direct{Arguments: opts})),
	)
	return err
}

// syncedFolderServer is a gRPC server that the client talks to and calls a
// real implementation of the component.
type syncedFolderServer struct {
	*BaseServer
	*capabilityServer

	Impl component.SyncedFolder
	vagrant_plugin_sdk.UnsafeSyncedFolderServiceServer
}

func (s *syncedFolderServer) UsableSpec(
	ctx context.Context,
	_ *emptypb.Empty,
) (*vagrant_plugin_sdk.FuncSpec, error) {
	if err := isImplemented(s, s.typ); err != nil {
		return nil, err
	}

	return s.GenerateSpec(s.Impl.UsableFunc())
}

func (s *syncedFolderServer) Usable(
	ctx context.Context,
	args *vagrant_plugin_sdk.FuncSpec_Args,
) (*vagrant_plugin_sdk.SyncedFolder_UsableResp, error) {
	raw, err := s.CallDynamicFunc(s.Impl.UsableFunc(), (*bool)(nil),
		args.Args, argmapper.Typed(ctx))

	if err != nil {
		s.Logger.Error("synced folder usable check failed",
			"error", err,
		)
		return nil, err
	}

	return &vagrant_plugin_sdk.SyncedFolder_UsableResp{
		Usable: raw.(bool)}, nil
}

func (s *syncedFolderServer) EnableSpec(
	ctx context.Context,
	_ *emptypb.Empty,
) (*vagrant_plugin_sdk.FuncSpec, error) {
	if err := isImplemented(s, s.typ); err != nil {
		return nil, err
	}

	return s.GenerateSpec(s.Impl.EnableFunc())
}

func (s *syncedFolderServer) Enable(
	ctx context.Context,
	args *vagrant_plugin_sdk.FuncSpec_Args,
) (*emptypb.Empty, error) {
	_, err := s.CallDynamicFunc(s.Impl.EnableFunc(), false,
		args.Args, argmapper.Typed(ctx))

	if err != nil {
		s.Logger.Error("synced folder enable failed",
			"error", err,
		)
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *syncedFolderServer) PrepareSpec(
	ctx context.Context,
	_ *emptypb.Empty,
) (*vagrant_plugin_sdk.FuncSpec, error) {
	if err := isImplemented(s, s.typ); err != nil {
		return nil, err
	}

	return s.GenerateSpec(s.Impl.PrepareFunc())
}

func (s *syncedFolderServer) Prepare(
	ctx context.Context,
	args *vagrant_plugin_sdk.FuncSpec_Args,
) (*emptypb.Empty, error) {
	_, err := s.CallDynamicFunc(s.Impl.PrepareFunc(), false,
		args.Args, argmapper.Typed(ctx))

	if err != nil {
		s.Logger.Error("synced folder prepare failed",
			"error", err,
		)
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *syncedFolderServer) DisableSpec(
	ctx context.Context,
	_ *emptypb.Empty,
) (*vagrant_plugin_sdk.FuncSpec, error) {
	if err := isImplemented(s, s.typ); err != nil {
		return nil, err
	}

	return s.GenerateSpec(s.Impl.DisableFunc())
}

func (s *syncedFolderServer) Disable(
	ctx context.Context,
	args *vagrant_plugin_sdk.FuncSpec_Args,
) (*emptypb.Empty, error) {
	_, err := s.CallDynamicFunc(s.Impl.DisableFunc(), false,
		args.Args, argmapper.Typed(ctx))

	if err != nil {
		s.Logger.Error("synced folder disable failed",
			"error", err,
		)
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *syncedFolderServer) CleanupSpec(
	ctx context.Context,
	_ *emptypb.Empty,
) (*vagrant_plugin_sdk.FuncSpec, error) {
	if err := isImplemented(s, s.typ); err != nil {
		return nil, err
	}

	return s.GenerateSpec(s.Impl.CleanupFunc())
}

func (s *syncedFolderServer) Cleanup(
	ctx context.Context,
	args *vagrant_plugin_sdk.FuncSpec_Args,
) (*emptypb.Empty, error) {
	_, err := s.CallDynamicFunc(s.Impl.CleanupFunc(), false,
		args.Args, argmapper.Typed(ctx))

	if err != nil {
		s.Logger.Error("synced folder cleanup failed",
			"error", err,
		)
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

var (
	_ plugin.Plugin                                = (*SyncedFolderPlugin)(nil)
	_ plugin.GRPCPlugin                            = (*SyncedFolderPlugin)(nil)
	_ vagrant_plugin_sdk.SyncedFolderServiceServer = (*syncedFolderServer)(nil)
	_ component.SyncedFolder                       = (*syncedFolderClient)(nil)
	_ core.SyncedFolder                            = (*syncedFolderClient)(nil)
	_ component.CapabilityPlatform                 = (*syncedFolderClient)(nil)
	_ core.CapabilityPlatform                      = (*syncedFolderClient)(nil)
	_ core.Seeder                                  = (*syncedFolderClient)(nil)
)
