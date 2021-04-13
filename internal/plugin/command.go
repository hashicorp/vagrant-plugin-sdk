package plugin

import (
	"context"

	"github.com/LK4D4/joincontext"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/hashicorp/go-argmapper"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/hashicorp/vagrant-plugin-sdk/component"
	"github.com/hashicorp/vagrant-plugin-sdk/core"
	"github.com/hashicorp/vagrant-plugin-sdk/docs"
	"github.com/hashicorp/vagrant-plugin-sdk/internal-shared/protomappers"
	"github.com/hashicorp/vagrant-plugin-sdk/internal/funcspec"
	"github.com/hashicorp/vagrant-plugin-sdk/proto/vagrant_plugin_sdk"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// CommandPlugin implements plugin.Plugin (specifically GRPCPlugin) for
// the Command component type.
type CommandPlugin struct {
	plugin.NetRPCUnsupportedPlugin

	Impl    component.Command // Impl is the concrete implementation
	Mappers []*argmapper.Func // Mappers
	Logger  hclog.Logger      // Logger
}

func (p *CommandPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	vagrant_plugin_sdk.RegisterCommandServiceServer(s, &commandServer{
		Impl: p.Impl,
		baseServer: &baseServer{
			base: &base{
				Mappers: p.Mappers,
				Logger:  p.Logger,
				Broker:  broker,
			},
		},
	})
	return nil
}

func (p *CommandPlugin) GRPCClient(
	ctx context.Context,
	broker *plugin.GRPCBroker,
	c *grpc.ClientConn,
) (interface{}, error) {
	return &commandClient{
		client: vagrant_plugin_sdk.NewCommandServiceClient(c),
		baseClient: &baseClient{
			ctx: context.Background(),
			base: &base{
				Mappers: p.Mappers,
				Logger:  p.Logger,
				Broker:  broker,
			},
		},
	}, nil
}

// commandClient is an implementation of component.Command over gRPC.
type commandClient struct {
	*baseClient

	client vagrant_plugin_sdk.CommandServiceClient
}

func (c *commandClient) Config() (interface{}, error) {
	return configStructCall(c.ctx, c.client)
}

func (c *commandClient) ConfigSet(v interface{}) error {
	return configureCall(c.ctx, c.client, v)
}

func (c *commandClient) Documentation() (*docs.Documentation, error) {
	return documentationCall(c.ctx, c.client)
}

func (c *commandClient) CommandFunc() interface{} {
	//TODO
	return nil
}

func (c *commandClient) Name() (name []string, err error) {
	meta, ok := metadata.FromOutgoingContext(c.ctx)
	if !ok {
		return
	}
	name = append(name, meta["plugin_name"]...)
	name = append(name, meta["command"]...)
	return
}

func (c *commandClient) CommandInfoFunc() interface{} {
	spec, err := c.client.CommandInfoSpec(c.ctx, &empty.Empty{})
	if err != nil {
		return funcErr(err)
	}
	spec.Result = nil
	cb := func(ctx context.Context, args funcspec.Args) (*vagrant_plugin_sdk.Command_CommandInfoResp, error) {
		ctx, _ = joincontext.Join(c.ctx, ctx)
		resp, err := c.client.CommandInfo(ctx, &vagrant_plugin_sdk.FuncSpec_Args{Args: args})
		if err != nil {
			return nil, err
		}
		return resp, nil
	}
	return c.generateFunc(spec, cb)
}

func (c *commandClient) CommandInfo() (*core.CommandInfo, error) {
	f := c.CommandInfoFunc()
	raw, err := c.callRemoteDynamicFunc(c.ctx, nil, (**vagrant_plugin_sdk.Command_CommandInfoResp)(nil), f)
	if err != nil {
		return nil, err
	}

	commandInfo := raw.(*vagrant_plugin_sdk.Command_CommandInfoResp)
	commandName, err := c.Name()
	if err != nil {
		return nil, err
	}
	flags, err := protomappers.Flags(commandInfo.Flags)
	return &core.CommandInfo{
		Name:     commandName,
		Help:     commandInfo.Help,
		Synopsis: commandInfo.Synopsis,
		Flags:    flags,
	}, nil
}

func (c *commandClient) ExecuteFunc() interface{} {
	spec, err := c.client.ExecuteSpec(c.ctx, &empty.Empty{})
	if err != nil {
		return funcErr(err)
	}
	spec.Result = nil
	cb := func(ctx context.Context, args funcspec.Args) (int64, error) {
		ctx, _ = joincontext.Join(c.ctx, ctx)
		resp, err := c.client.Execute(ctx, &vagrant_plugin_sdk.FuncSpec_Args{Args: args})
		if err != nil {
			return -1, err
		}
		return resp.ExitCode, nil
	}
	return c.generateFunc(spec, cb)
}

func (c *commandClient) Execute(name string) (int64, error) {
	f := c.ExecuteFunc()
	raw, err := c.callRemoteDynamicFunc(c.ctx, nil, (*int64)(nil), f)
	if err != nil {
		return -1, err
	}

	return raw.(int64), nil
}

func (c *commandClient) SubcommandsFunc() interface{} {
	spec, err := c.client.SubcommandSpec(c.ctx, &empty.Empty{})
	if err != nil {
		return funcErr(err)
	}
	spec.Result = nil
	cb := func(ctx context.Context, args funcspec.Args) (*vagrant_plugin_sdk.Command_SubcommandResp, error) {
		ctx, _ = joincontext.Join(c.ctx, ctx)
		resp, err := c.client.Subcommands(ctx, &vagrant_plugin_sdk.FuncSpec_Args{Args: args})
		if err != nil {
			return nil, err
		}
		return resp, nil
	}
	return c.generateFunc(spec, cb)
}

func (c *commandClient) Subcommands() ([]core.Command, error) {
	f := c.SubcommandsFunc()
	raw, err := c.callRemoteDynamicFunc(c.ctx, nil, (**vagrant_plugin_sdk.Command_SubcommandResp)(nil), f)
	if err != nil {
		return nil, err
	}

	res := []core.Command{}
	subcommands := raw.(*vagrant_plugin_sdk.Command_SubcommandResp).Commands
	for _, cmd := range subcommands {
		sc_client := &commandClient{
			client: c.client,
			baseClient: &baseClient{
				ctx: c.ctx,
				base: &base{
					Mappers: c.Mappers,
					Logger:  c.Logger,
					Broker:  c.Broker,
				},
			},
		}
		sc_client.SetRequestMetadata("command", cmd)
		res = append(res, sc_client)
	}
	if err != nil {
		return nil, err
	}
	return res, nil
}

// commandServer is a gRPC server that the client talks to and calls a
// real implementation of the component.
type commandServer struct {
	*baseServer

	Impl component.Command
	vagrant_plugin_sdk.UnimplementedCommandServiceServer
}

func (s *commandServer) ConfigStruct(
	ctx context.Context,
	empty *empty.Empty,
) (*vagrant_plugin_sdk.Config_StructResp, error) {
	return configStruct(s.Impl)
}

func (s *commandServer) Configure(
	ctx context.Context,
	req *vagrant_plugin_sdk.Config_ConfigureRequest,
) (*empty.Empty, error) {
	return configure(s.Impl, req)
}

func (s *commandServer) Documentation(
	ctx context.Context,
	empty *empty.Empty,
) (*vagrant_plugin_sdk.Config_Documentation, error) {
	return documentation(s.Impl)
}

func (s *commandServer) CommandInfoSpec(
	ctx context.Context,
	_ *empty.Empty,
) (*vagrant_plugin_sdk.FuncSpec, error) {
	if err := isImplemented(s, "command"); err != nil {
		return nil, err
	}

	return s.generateSpec(s.Impl.CommandInfoFunc())
}

func (s *commandServer) CommandInfo(
	ctx context.Context,
	args *vagrant_plugin_sdk.FuncSpec_Args,
) (*vagrant_plugin_sdk.Command_CommandInfoResp, error) {
	raw, err := s.callLocalDynamicFunc(
		s.Impl.CommandInfoFunc(),
		args.Args,
		(*core.CommandInfo)(nil),
		argmapper.Typed(ctx),
	)

	if err != nil {
		return nil, err
	}

	commandInfo, err := protomappers.CommandInfoProto(raw.(*core.CommandInfo))
	return commandInfo, nil
}

func (s *commandServer) ExecuteSpec(
	ctx context.Context,
	_ *empty.Empty,
) (*vagrant_plugin_sdk.FuncSpec, error) {
	if err := isImplemented(s, "command"); err != nil {
		return nil, err
	}

	return s.generateSpec(s.Impl.ExecuteFunc())
}

func (s *commandServer) Execute(
	ctx context.Context,
	args *vagrant_plugin_sdk.FuncSpec_Args,
) (*vagrant_plugin_sdk.Command_ExecuteResp, error) {
	raw, err := s.callUncheckedLocalDynamicFunc(
		s.Impl.ExecuteFunc(),
		args.Args,
		argmapper.Typed(ctx),
	)

	if err != nil {
		return nil, err
	}

	result := &vagrant_plugin_sdk.Command_ExecuteResp{
		ExitCode: raw.(int64),
	}
	return result, nil
}

func (s *commandServer) SubcommandSpec(
	ctx context.Context,
	_ *empty.Empty,
) (*vagrant_plugin_sdk.FuncSpec, error) {
	if err := isImplemented(s, "command"); err != nil {
		return nil, err
	}

	return s.generateSpec(s.Impl.SubcommandsFunc())
}

func (s *commandServer) Subcommands(
	ctx context.Context,
	args *vagrant_plugin_sdk.FuncSpec_Args,
) (*vagrant_plugin_sdk.Command_SubcommandResp, error) {
	raw, err := s.callLocalDynamicFunc(
		s.Impl.SubcommandsFunc(),
		args.Args,
		([]string)(nil),
		argmapper.Typed(ctx),
	)

	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}
	msg := &vagrant_plugin_sdk.Command_SubcommandResp{
		Commands: raw.([]string),
	}

	return msg, nil
}

var (
	_ plugin.Plugin                           = (*CommandPlugin)(nil)
	_ plugin.GRPCPlugin                       = (*CommandPlugin)(nil)
	_ vagrant_plugin_sdk.CommandServiceServer = (*commandServer)(nil)
	_ component.Command                       = (*commandClient)(nil)
	_ core.Command                            = (*commandClient)(nil)
)
