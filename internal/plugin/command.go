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
)

// CommandPlugin implements plugin.Plugin (specifically GRPCPlugin) for
// the Command component type.
type CommandPlugin struct {
	plugin.NetRPCUnsupportedPlugin

	Impl    []component.Command // Impl is the concrete implementation
	Mappers []*argmapper.Func   // Mappers
	Logger  hclog.Logger        // Logger
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

func (c *commandClient) CommandInfoFunc() interface{} {
	// TODO: set this command string
	req := &vagrant_plugin_sdk.Command_SpecReq{CommandString: []string{}}
	spec, err := c.client.CommandInfoSpec(c.ctx, req)
	if err != nil {
		return funcErr(err)
	}
	spec.Result = nil
	cb := func(ctx context.Context, args funcspec.Args) (*vagrant_plugin_sdk.Command_CommandInfoResp, error) {
		ctx, _ = joincontext.Join(c.ctx, ctx)
		// TODO: make this take the name
		resp, err := c.client.CommandInfo(ctx, &vagrant_plugin_sdk.Command_CommandInfoReq{CommandString: []string{}})
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

	commandInfo, err := protomappers.CommandInfo(raw.(*vagrant_plugin_sdk.Command_CommandInfoResp).CommandInfo)
	return commandInfo, err
}

func (c *commandClient) ExecuteFunc() interface{} {
	// TODO:
	req := &vagrant_plugin_sdk.Command_SpecReq{CommandString: []string{"myplugin", "info"}}
	spec, err := c.client.ExecuteSpec(c.ctx, req)
	if err != nil {
		return funcErr(err)
	}
	spec.Result = nil
	cb := func(ctx context.Context, args funcspec.Args) (int64, error) {
		ctx, _ = joincontext.Join(c.ctx, ctx)
		funcspecArgs := &vagrant_plugin_sdk.FuncSpec_Args{Args: args}
		executeArgs := &vagrant_plugin_sdk.Command_ExecuteReq{
			Args:          funcspecArgs,
			CommandString: []string{"myplugin", "info"},
		}
		resp, err := c.client.Execute(ctx, executeArgs)
		if err != nil {
			return -1, err
		}
		return resp.ExitCode, nil
	}
	return c.generateFunc(spec, cb)
}

func (c *commandClient) Execute(name string) (int64, error) {
	f := c.ExecuteFunc()
	raw, err := c.callRemoteDynamicFunc(c.ctx, c.Mappers, (*int64)(nil), f)
	if err != nil {
		return -1, err
	}

	return raw.(int64), nil
}

// commandServer is a gRPC server that the client talks to and calls a
// real implementation of the component.
type commandServer struct {
	*baseServer

	Impl []component.Command
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

func (s *commandServer) FindCmd(
	ctx context.Context,
	cmds []component.Command,
	cmdName []string,
) (component.Command, error) {
	if len(cmdName) == 0 {
		return cmds[0], nil
	}
	for _, cmd := range cmds {
		raw, err := s.callLocalDynamicFunc(
			cmd.CommandInfoFunc(),
			nil,
			(*core.CommandInfo)(nil),
			argmapper.Typed(ctx),
		)
		if err != nil {
			return nil, err
		}
		cmdInfo := raw.(*core.CommandInfo)
		if cmdInfo.Name == cmdName[len(cmdName)-1] {
			return cmd, nil
		}
	}
	return cmds[0], nil
}

func (s *commandServer) Documentation(
	ctx context.Context,
	empty *empty.Empty,
) (*vagrant_plugin_sdk.Config_Documentation, error) {
	return documentation(s.Impl)
}

func (s *commandServer) CommandInfoSpec(
	ctx context.Context,
	req *vagrant_plugin_sdk.Command_SpecReq,
) (*vagrant_plugin_sdk.FuncSpec, error) {
	if err := isImplemented(s, "command"); err != nil {
		return nil, err
	}

	impl, err := s.FindCmd(ctx, s.Impl, req.CommandString)
	if err != nil {
		return nil, err
	}
	return s.generateSpec(impl.CommandInfoFunc())
}

func (s *commandServer) CommandInfo(
	ctx context.Context,
	req *vagrant_plugin_sdk.Command_CommandInfoReq,
) (*vagrant_plugin_sdk.Command_CommandInfoResp, error) {
	impl, err := s.FindCmd(ctx, s.Impl, req.CommandString)
	if err != nil {
		return nil, err
	}
	raw, err := s.callLocalDynamicFunc(
		impl.CommandInfoFunc(),
		nil,
		(*core.CommandInfo)(nil),
		argmapper.Typed(ctx),
	)

	if err != nil {
		return nil, err
	}

	commandInfo, err := protomappers.CommandInfoProto(raw.(*core.CommandInfo))
	return &vagrant_plugin_sdk.Command_CommandInfoResp{
		CommandInfo: commandInfo,
	}, nil
}

func (s *commandServer) ExecuteSpec(
	ctx context.Context,
	req *vagrant_plugin_sdk.Command_SpecReq,
) (*vagrant_plugin_sdk.FuncSpec, error) {
	if err := isImplemented(s, "command"); err != nil {
		return nil, err
	}

	impl, err := s.FindCmd(ctx, s.Impl, req.CommandString)
	if err != nil {
		return nil, err
	}
	return s.generateSpec(impl.ExecuteFunc())
}

func (s *commandServer) Execute(
	ctx context.Context,
	args *vagrant_plugin_sdk.Command_ExecuteReq,
) (*vagrant_plugin_sdk.Command_ExecuteResp, error) {
	impl, err := s.FindCmd(ctx, s.Impl, args.CommandString)
	if err != nil {
		return nil, err
	}
	raw, err := s.callUncheckedLocalDynamicFunc(
		impl.ExecuteFunc(),
		args.Args.Args,
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

var (
	_ plugin.Plugin                           = (*CommandPlugin)(nil)
	_ plugin.GRPCPlugin                       = (*CommandPlugin)(nil)
	_ vagrant_plugin_sdk.CommandServiceServer = (*commandServer)(nil)
	_ component.Command                       = (*commandClient)(nil)
	_ core.Command                            = (*commandClient)(nil)
)
