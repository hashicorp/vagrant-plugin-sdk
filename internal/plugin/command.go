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

// CommandPlugin implements plugin.Plugin (specifically GRPCPlugin) for
// the Command component type.
type CommandPlugin struct {
	plugin.NetRPCUnsupportedPlugin

	Impl component.Command // Impl is the concrete implementation
	*BasePlugin
}

func (p *CommandPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	vagrant_plugin_sdk.RegisterCommandServiceServer(s,
		&commandServer{
			Impl:       p.Impl,
			BaseServer: p.NewServer(broker, p.Impl),
		},
	)
	return nil
}

func (p *CommandPlugin) GRPCClient(
	ctx context.Context,
	broker *plugin.GRPCBroker,
	c *grpc.ClientConn,
) (interface{}, error) {
	cl := vagrant_plugin_sdk.NewCommandServiceClient(c)
	return &commandClient{
		client:     cl,
		BaseClient: p.NewClient(ctx, broker, cl.(SeederClient)),
	}, nil
}

// commandClient is an implementation of component.Command over gRPC.
type commandClient struct {
	*BaseClient

	client vagrant_plugin_sdk.CommandServiceClient
}

func (c *commandClient) CommandInfoFunc() interface{} {
	spec, err := c.client.CommandInfoSpec(c.Ctx, &emptypb.Empty{})
	if err != nil {
		return funcErr(err)
	}
	spec.Result = nil
	cb := func(ctx context.Context, args funcspec.Args) (*vagrant_plugin_sdk.Command_CommandInfo, error) {
		ctx, _ = joincontext.Join(c.Ctx, ctx)
		resp, err := c.client.CommandInfo(
			ctx, &vagrant_plugin_sdk.FuncSpec_Args{Args: args},
		)
		if err != nil {
			return nil, err
		}
		return resp.CommandInfo, nil
	}
	return c.GenerateFunc(spec, cb)
}

func (c *commandClient) CommandInfo() (*component.CommandInfo, error) {
	f := c.CommandInfoFunc()
	raw, err := c.CallDynamicFunc(f, (**component.CommandInfo)(nil),
		argmapper.Typed(c.Ctx))
	if err != nil {
		return nil, err
	}

	return raw.(*component.CommandInfo), err
}

func (c *commandClient) ExecuteFunc(cliArgs []string) interface{} {
	spec, err := c.client.ExecuteSpec(c.Ctx, &vagrant_plugin_sdk.Command_ExecuteSpecReq{
		CommandArgs: cliArgs})
	if err != nil {
		return funcErr(err)
	}
	spec.Result = nil
	cb := func(ctx context.Context, args funcspec.Args) (int32, error) {
		ctx, _ = joincontext.Join(c.Ctx, ctx)
		executeArgs := &vagrant_plugin_sdk.Command_ExecuteReq{
			Spec:        &vagrant_plugin_sdk.FuncSpec_Args{Args: args},
			CommandArgs: cliArgs,
		}
		resp, err := c.client.Execute(ctx, executeArgs)
		if err != nil {
			return -1, err
		}
		return resp.ExitCode, nil
	}
	return c.GenerateFunc(spec, cb)
}

func (c *commandClient) Execute(cliArgs []string) (int32, error) {
	f := c.ExecuteFunc(cliArgs)

	raw, err := c.CallDynamicFunc(f, (*int32)(nil),
		argmapper.Typed(c.Ctx))
	if err != nil {
		return -1, err
	}

	return raw.(int32), nil
}

// commandServer is a gRPC server that the client talks to and calls a
// real implementation of the component.
type commandServer struct {
	*BaseServer

	Impl component.Command
	vagrant_plugin_sdk.UnsafeCommandServiceServer
}

func (s *commandServer) CommandInfoSpec(
	ctx context.Context,
	_ *emptypb.Empty,
) (*vagrant_plugin_sdk.FuncSpec, error) {
	if err := isImplemented(s, "command"); err != nil {
		return nil, err
	}

	return s.GenerateSpec(s.Impl.CommandInfoFunc())
}

func (s *commandServer) CommandInfo(
	ctx context.Context,
	req *vagrant_plugin_sdk.FuncSpec_Args,
) (*vagrant_plugin_sdk.Command_CommandInfoResp, error) {
	raw, err := s.CallDynamicFunc(s.Impl.CommandInfoFunc(),
		(**vagrant_plugin_sdk.Command_CommandInfo)(nil),
		req.Args,
		argmapper.Typed(ctx),
	)

	if err != nil {
		s.Logger.Error("command info failed",
			"error", err,
		)

		return nil, err
	}

	return &vagrant_plugin_sdk.Command_CommandInfoResp{
		CommandInfo: raw.(*vagrant_plugin_sdk.Command_CommandInfo),
	}, nil
}

func (s *commandServer) ExecuteSpec(
	ctx context.Context,
	req *vagrant_plugin_sdk.Command_ExecuteSpecReq,
) (*vagrant_plugin_sdk.FuncSpec, error) {
	if err := isImplemented(s, "command"); err != nil {
		return nil, err
	}
	return s.GenerateSpec(s.Impl.ExecuteFunc(req.CommandArgs))
}

func (s *commandServer) Execute(
	ctx context.Context,
	req *vagrant_plugin_sdk.Command_ExecuteReq,
) (*vagrant_plugin_sdk.Command_ExecuteResp, error) {
	raw, err := s.CallDynamicFunc(s.Impl.ExecuteFunc(req.CommandArgs),
		(*int32)(nil),
		req.Spec.Args,
		argmapper.Typed(ctx),
	)

	if err != nil {
		s.Logger.Error("command execution failed",
			"error", err,
		)

		return nil, err
	}

	result := &vagrant_plugin_sdk.Command_ExecuteResp{
		ExitCode: raw.(int32),
	}
	return result, nil
}

var (
	_ plugin.Plugin                           = (*CommandPlugin)(nil)
	_ plugin.GRPCPlugin                       = (*CommandPlugin)(nil)
	_ vagrant_plugin_sdk.CommandServiceServer = (*commandServer)(nil)
	_ component.Command                       = (*commandClient)(nil)
	_ core.Command                            = (*commandClient)(nil)
	_ core.Seeder                             = (*commandClient)(nil)
)
