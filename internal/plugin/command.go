package plugin

import (
	"context"

	"github.com/LK4D4/joincontext"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/hashicorp/go-argmapper"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"

	"github.com/hashicorp/vagrant-plugin-sdk/component"
	"github.com/hashicorp/vagrant-plugin-sdk/docs"
	"github.com/hashicorp/vagrant-plugin-sdk/internal/funcspec"
	pb "github.com/hashicorp/vagrant-plugin-sdk/proto/gen"
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
	pb.RegisterCommandServiceServer(s, &commandServer{
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
		client: pb.NewCommandServiceClient(c),
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

	client pb.CommandServiceClient
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

func (c *commandClient) SynopsisFunc() interface{} {
	spec, err := c.client.SynopsisSpec(c.ctx, &empty.Empty{})
	if err != nil {
		return funcErr(err)
	}
	spec.Result = nil
	cb := func(ctx context.Context, args funcspec.Args) (string, error) {
		ctx, _ = joincontext.Join(c.ctx, ctx)
		resp, err := c.client.Synopsis(ctx, &pb.FuncSpec_Args{Args: args})
		if err != nil {
			return "", err
		}
		return resp.Synopsis, nil
	}
	return c.generateFunc(spec, cb)
}

func (c *commandClient) Synopsis() (string, error) {
	f := c.SynopsisFunc()
	raw, err := c.callRemoteDynamicFunc(c.ctx, nil, (*string)(nil), f)
	if err != nil {
		return "", err
	}
	return raw.(string), nil
}

func (c *commandClient) HelpFunc() interface{} {
	spec, err := c.client.HelpSpec(c.ctx, &empty.Empty{})
	if err != nil {
		return funcErr(err)
	}
	spec.Result = nil
	cb := func(ctx context.Context, args funcspec.Args) (string, error) {
		ctx, _ = joincontext.Join(c.ctx, ctx)
		resp, err := c.client.Help(ctx, &pb.FuncSpec_Args{Args: args})
		if err != nil {
			return "", err
		}
		return resp.Help, nil
	}
	return c.generateFunc(spec, cb)
}

func (c *commandClient) Help() (string, error) {
	f := c.HelpFunc()
	raw, err := c.callRemoteDynamicFunc(c.ctx, nil, (*string)(nil), f)
	if err != nil {
		return "", err
	}
	return raw.(string), nil
}

func (c *commandClient) FlagsFunc() interface{} {
	spec, err := c.client.FlagsSpec(c.ctx, &empty.Empty{})
	if err != nil {
		return funcErr(err)
	}
	spec.Result = nil
	cb := func(ctx context.Context, args funcspec.Args) (string, error) {
		ctx, _ = joincontext.Join(c.ctx, ctx)
		resp, err := c.client.Flags(ctx, &pb.FuncSpec_Args{Args: args})
		if err != nil {
			return "", err
		}
		return resp.Flags, nil
	}
	return c.generateFunc(spec, cb)
}

func (c *commandClient) Flags() (string, error) {
	f := c.FlagsFunc()
	raw, err := c.callRemoteDynamicFunc(c.ctx, nil, (*string)(nil), f)
	if err != nil {
		return "", err
	}
	return raw.(string), nil
}

func (c *commandClient) ExecuteFunc() interface{} {
	spec, err := c.client.ExecuteSpec(c.ctx, &empty.Empty{})
	if err != nil {
		return funcErr(err)
	}
	spec.Result = nil
	cb := func(ctx context.Context, args funcspec.Args) (int64, error) {
		ctx, _ = joincontext.Join(c.ctx, ctx)
		resp, err := c.client.Execute(ctx, &pb.FuncSpec_Args{Args: args})
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

// commandServer is a gRPC server that the client talks to and calls a
// real implementation of the component.
type commandServer struct {
	*baseServer

	Impl component.Command
}

func (s *commandServer) ConfigStruct(
	ctx context.Context,
	empty *empty.Empty,
) (*pb.Config_StructResp, error) {
	return configStruct(s.Impl)
}

func (s *commandServer) Configure(
	ctx context.Context,
	req *pb.Config_ConfigureRequest,
) (*empty.Empty, error) {
	return configure(s.Impl, req)
}

func (s *commandServer) Documentation(
	ctx context.Context,
	empty *empty.Empty,
) (*pb.Config_Documentation, error) {
	return documentation(s.Impl)
}

func (s *commandServer) SynopsisSpec(
	ctx context.Context,
	_ *empty.Empty,
) (*pb.FuncSpec, error) {
	if err := isImplemented(s, "command"); err != nil {
		return nil, err
	}

	return s.generateSpec(s.Impl.SynopsisFunc())
}

func (s *commandServer) Synopsis(
	ctx context.Context,
	args *pb.FuncSpec_Args,
) (*pb.Command_SynopsisResp, error) {
	raw, err := s.callLocalDynamicFunc(
		s.Impl.SynopsisFunc(), args.Args, argmapper.Typed(ctx),
	)

	if err != nil {
		return nil, err
	}

	return &pb.Command_SynopsisResp{Synopsis: raw.(string)}, nil
}

func (s *commandServer) HelpSpec(
	ctx context.Context,
	_ *empty.Empty,
) (*pb.FuncSpec, error) {
	if err := isImplemented(s, "command"); err != nil {
		return nil, err
	}

	return s.generateSpec(s.Impl.HelpFunc())
}

func (s *commandServer) Help(
	ctx context.Context,
	args *pb.FuncSpec_Args,
) (*pb.Command_HelpResp, error) {
	raw, err := s.callLocalDynamicFunc(
		s.Impl.HelpFunc(), args.Args, argmapper.Typed(ctx),
	)

	if err != nil {
		return nil, err
	}

	return &pb.Command_HelpResp{Help: raw.(string)}, nil
}

func (s *commandServer) FlagsSpec(
	ctx context.Context,
	_ *empty.Empty,
) (*pb.FuncSpec, error) {
	if err := isImplemented(s, "command"); err != nil {
		return nil, err
	}

	return s.generateSpec(s.Impl.FlagsFunc())
}

func (s *commandServer) Flags(
	ctx context.Context,
	args *pb.FuncSpec_Args,
) (*pb.Command_FlagsResp, error) {
	raw, err := s.callLocalDynamicFunc(
		s.Impl.FlagsFunc(), args.Args, argmapper.Typed(ctx),
	)

	if err != nil {
		return nil, err
	}

	return &pb.Command_FlagsResp{Flags: raw.(string)}, nil
}

func (s *commandServer) ExecuteSpec(
	ctx context.Context,
	_ *empty.Empty,
) (*pb.FuncSpec, error) {
	if err := isImplemented(s, "command"); err != nil {
		return nil, err
	}

	return s.generateSpec(s.Impl.ExecuteFunc())
}

func (s *commandServer) Execute(
	ctx context.Context,
	args *pb.FuncSpec_Args,
) (*pb.Command_ExecuteResp, error) {
	raw, err := s.callLocalDynamicFunc(
		s.Impl.ExecuteFunc(), args.Args, argmapper.Typed(ctx),
	)

	if err != nil {
		return nil, err
	}

	return &pb.Command_ExecuteResp{ExitCode: raw.(int64)}, nil
}

var (
	_ plugin.Plugin           = (*CommandPlugin)(nil)
	_ plugin.GRPCPlugin       = (*CommandPlugin)(nil)
	_ pb.CommandServiceServer = (*commandServer)(nil)
	_ component.Command       = (*commandClient)(nil)
)
