package plugin

import (
	"context"

	"github.com/LK4D4/joincontext"
	"github.com/hashicorp/go-argmapper"
	"github.com/hashicorp/go-plugin"
	"github.com/hashicorp/vagrant-plugin-sdk/component"
	"github.com/hashicorp/vagrant-plugin-sdk/core"
	"github.com/hashicorp/vagrant-plugin-sdk/internal/funcspec"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/hashicorp/vagrant-plugin-sdk/proto/vagrant_plugin_sdk"
)

// PushPlugin implements plugin.Plugin (specifically GRPCPlugin) for
// the Push component type.
type PushPlugin struct {
	plugin.NetRPCUnsupportedPlugin

	Impl component.Push // Impl is the concrete implementation
	*BasePlugin
}

func (p *PushPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	vagrant_plugin_sdk.RegisterPushServiceServer(s, &pushServer{
		Impl:       p.Impl,
		BaseServer: p.NewServer(broker, p.Impl),
	})
	return nil
}

func (p *PushPlugin) GRPCClient(
	ctx context.Context,
	broker *plugin.GRPCBroker,
	c *grpc.ClientConn,
) (interface{}, error) {
	cl := vagrant_plugin_sdk.NewPushServiceClient(c)
	return &pushClient{
		client:     cl,
		BaseClient: p.NewClient(ctx, broker, cl.(SeederClient)),
	}, nil
}

// pushClient is an implementation of component.Push over gRPC.
type pushClient struct {
	*BaseClient

	client vagrant_plugin_sdk.PushServiceClient
}

// this meets the component.Push interface
func (c *pushClient) PushFunc() interface{} {
	spec, err := c.client.PushSpec(c.Ctx, &emptypb.Empty{})
	if err != nil {
		return funcErr(err)
	}
	spec.Result = nil
	cb := func(ctx context.Context, args funcspec.Args) error {
		ctx, _ = joincontext.Join(c.Ctx, ctx)
		_, err := c.client.Push(ctx, &vagrant_plugin_sdk.FuncSpec_Args{Args: args})
		if err != nil {
			return err
		}
		return nil
	}
	return c.GenerateFunc(spec, cb)
}

// this meets the core.Push interfacce
func (c *pushClient) Push() error {
	f := c.PushFunc()
	_, err := c.CallDynamicFunc(f, false, argmapper.Typed(c.Ctx))
	return err
}

// pushServer is a gRPC server that the client talks to and calls a
// real implementation of the component.
type pushServer struct {
	*BaseServer

	Impl component.Push
	vagrant_plugin_sdk.UnsafePushServiceServer
}

func (s *pushServer) PushSpec(
	ctx context.Context,
	args *emptypb.Empty,
) (*vagrant_plugin_sdk.FuncSpec, error) {
	if err := isImplemented(s, "push"); err != nil {
		return nil, err
	}

	return s.GenerateSpec(s.Impl.PushFunc())
}

func (s *pushServer) Push(
	ctx context.Context,
	args *vagrant_plugin_sdk.FuncSpec_Args,
) (*emptypb.Empty, error) {
	_, err := s.CallDynamicFunc(s.Impl.PushFunc(), (*int32)(nil), args.Args,
		argmapper.Typed(ctx))

	return &emptypb.Empty{}, err
}

var (
	_ plugin.Plugin                        = (*PushPlugin)(nil)
	_ plugin.GRPCPlugin                    = (*PushPlugin)(nil)
	_ vagrant_plugin_sdk.PushServiceServer = (*pushServer)(nil)
	_ component.Push                       = (*pushClient)(nil)
	_ core.Push                            = (*pushClient)(nil)
	_ core.Seeder                          = (*pushClient)(nil)
)
