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

// GuestPlugin implements plugin.Plugin (specifically GRPCPlugin) for
// the Guest component type.
type GuestPlugin struct {
	plugin.NetRPCUnsupportedPlugin

	Impl component.Guest // Impl is the concrete implementation
	*BasePlugin
}

func (p *GuestPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	bs := p.NewServer(broker, p.Impl)
	vagrant_plugin_sdk.RegisterGuestServiceServer(s, &guestServer{
		Impl:       p.Impl,
		BaseServer: bs,
		capabilityServer: &capabilityServer{
			BaseServer:     bs,
			CapabilityImpl: p.Impl,
			typ:            "guest",
		},
	})
	return nil
}

func (p *GuestPlugin) GRPCClient(
	ctx context.Context,
	broker *plugin.GRPCBroker,
	c *grpc.ClientConn,
) (interface{}, error) {
	client := vagrant_plugin_sdk.NewGuestServiceClient(c)
	bc := p.NewClient(ctx, broker, client.(SeederClient))
	return &guestClient{
		client:     client,
		BaseClient: bc,
		capabilityClient: &capabilityClient{
			client:     client,
			BaseClient: bc,
		},
	}, nil
}

// guestClient is an implementation of component.Guest over gRPC.
type guestClient struct {
	*BaseClient
	*capabilityClient
	client vagrant_plugin_sdk.GuestServiceClient
}

func (c *guestClient) GetCapabilityClient() *capabilityClient {
	return c.capabilityClient
}

func (c *guestClient) GuestDetectFunc() interface{} {
	spec, err := c.client.DetectSpec(c.Ctx, &emptypb.Empty{})
	if err != nil {
		return funcErr(err)
	}
	spec.Result = nil
	cb := func(ctx context.Context, args funcspec.Args) (bool, error) {
		resp, err := c.client.Detect(ctx, &vagrant_plugin_sdk.FuncSpec_Args{Args: args})
		if err != nil {
			return false, err
		}
		return resp.Detected, nil
	}
	return c.GenerateFunc(spec, cb)
}

func (c *guestClient) Detect(t core.Target) (bool, error) {
	f := c.GuestDetectFunc()
	raw, err := c.CallDynamicFunc(f, (*bool)(nil),
		argmapper.Typed(c.Ctx),
		argmapper.Typed(t),
	)
	if err != nil {
		return false, err
	}

	return raw.(bool), nil
}

func (c *guestClient) ParentFunc() interface{} {
	spec, err := c.client.ParentSpec(c.Ctx, &emptypb.Empty{})
	if err != nil {
		return funcErr(err)
	}
	spec.Result = nil
	cb := func(ctx context.Context, args funcspec.Args) (string, error) {
		ctx, _ = joincontext.Join(c.Ctx, ctx)
		resp, err := c.client.Parent(ctx, &vagrant_plugin_sdk.FuncSpec_Args{Args: args})
		if err != nil {
			return "", err
		}
		return resp.Parent, nil
	}

	return c.GenerateFunc(spec, cb)
}

func (c *guestClient) Parent() (string, error) {
	f := c.ParentFunc()
	raw, err := c.CallDynamicFunc(f, (*string)(nil),
		argmapper.Typed(c.Ctx),
	)
	if err != nil {
		return "", err
	}

	return raw.(string), nil
}

// guestServer is a gRPC server that the client talks to and calls a
// real implementation of the component.
type guestServer struct {
	*BaseServer
	*capabilityServer

	Impl component.Guest
}

func (s *guestServer) DetectSpec(
	ctx context.Context,
	args *emptypb.Empty,
) (*vagrant_plugin_sdk.FuncSpec, error) {
	if err := isImplemented(s, "guest"); err != nil {
		return nil, err
	}

	return s.GenerateSpec(s.Impl.GuestDetectFunc())
}

func (s *guestServer) Detect(
	ctx context.Context,
	args *vagrant_plugin_sdk.FuncSpec_Args,
) (*vagrant_plugin_sdk.Platform_DetectResp, error) {
	raw, err := s.CallDynamicFunc(s.Impl.GuestDetectFunc(), (*bool)(nil), args.Args,
		argmapper.Typed(ctx),
	)

	if err != nil {
		s.Logger.Error("guest detect failed",
			"error", err,
		)

		return nil, err
	}

	return &vagrant_plugin_sdk.Platform_DetectResp{Detected: raw.(bool)}, nil
}

func (s *guestServer) ParentSpec(
	ctx context.Context,
	_ *emptypb.Empty,
) (*vagrant_plugin_sdk.FuncSpec, error) {
	if err := isImplemented(s, s.typ); err != nil {
		return nil, err
	}

	return s.GenerateSpec(s.Impl.ParentFunc())
}

func (s *guestServer) Parent(
	ctx context.Context,
	args *vagrant_plugin_sdk.FuncSpec_Args,
) (*vagrant_plugin_sdk.Platform_ParentResp, error) {
	raw, err := s.CallDynamicFunc(s.Impl.ParentFunc(), (*string)(nil),
		args.Args, argmapper.Typed(ctx))

	if err != nil {
		s.Logger.Error("guest parent failed",
			"error", err,
		)

		return nil, err
	}

	return &vagrant_plugin_sdk.Platform_ParentResp{
		Parent: raw.(string)}, nil
}

var (
	_ plugin.Plugin                         = (*GuestPlugin)(nil)
	_ plugin.GRPCPlugin                     = (*GuestPlugin)(nil)
	_ vagrant_plugin_sdk.GuestServiceServer = (*guestServer)(nil)
	_ component.Guest                       = (*guestClient)(nil)
	_ core.Guest                            = (*guestClient)(nil)
	_ component.CapabilityPlatform          = (*guestClient)(nil)
	_ core.CapabilityPlatform               = (*guestClient)(nil)
	_ core.Seeder                           = (*guestClient)(nil)
)
