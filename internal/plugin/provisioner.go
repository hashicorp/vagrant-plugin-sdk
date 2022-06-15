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

// ProvisionerPlugin implements plugin.Plugin (specifically GRPCPlugin) for
// the Provisioner component type.
type ProvisionerPlugin struct {
	plugin.NetRPCUnsupportedPlugin

	Impl component.Provisioner // Impl is the concrete implementation
	*BasePlugin
}

func (p *ProvisionerPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	vagrant_plugin_sdk.RegisterProvisionerServiceServer(s, &provisionerServer{
		Impl:       p.Impl,
		BaseServer: p.NewServer(broker, p.Impl),
	})
	return nil
}

func (p *ProvisionerPlugin) GRPCClient(
	ctx context.Context,
	broker *plugin.GRPCBroker,
	c *grpc.ClientConn,
) (interface{}, error) {
	cl := vagrant_plugin_sdk.NewProvisionerServiceClient(c)
	return &provisionerClient{
		client:     cl,
		BaseClient: p.NewClient(ctx, broker, cl.(SeederClient)),
	}, nil
}

// provisionerClient is an implementation of component.Provisioner over gRPC.
type provisionerClient struct {
	*BaseClient

	client vagrant_plugin_sdk.ProvisionerServiceClient
}

func (c *provisionerClient) CleanupFunc() interface{} {
	spec, err := c.client.CleanupSpec(c.Ctx, &emptypb.Empty{})
	if err != nil {
		return funcErr(err)
	}
	spec.Result = nil
	cb := func(ctx context.Context, args funcspec.Args) error {
		ctx, _ = joincontext.Join(c.Ctx, ctx)
		_, err := c.client.Cleanup(ctx, &vagrant_plugin_sdk.FuncSpec_Args{Args: args})
		if err != nil {
			return err
		}
		return nil
	}
	return c.GenerateFunc(spec, cb)
}

func (c *provisionerClient) Cleanup(machine core.Machine, config *component.ConfigData) error {
	f := c.CleanupFunc()

	_, err := c.CallDynamicFunc(f, false,
		argmapper.Typed(c.Ctx),
		argmapper.Typed(machine),
		argmapper.Typed(config),
	)

	return err
}

func (c *provisionerClient) ConfigureFunc() interface{} {
	spec, err := c.client.ConfigureSpec(c.Ctx, &emptypb.Empty{})
	if err != nil {
		return funcErr(err)
	}
	spec.Result = nil
	cb := func(ctx context.Context, args funcspec.Args) error {
		ctx, _ = joincontext.Join(c.Ctx, ctx)
		_, err := c.client.Configure(ctx, &vagrant_plugin_sdk.FuncSpec_Args{Args: args})
		if err != nil {
			return err
		}
		return nil
	}
	return c.GenerateFunc(spec, cb)
}

func (c *provisionerClient) Configure(machine core.Machine, config, rootConfig *component.ConfigData) error {
	f := c.CleanupFunc()

	_, err := c.CallDynamicFunc(f, false,
		argmapper.Typed(c.Ctx),
		argmapper.Typed(machine),
		argmapper.Typed(config),
		argmapper.Typed(rootConfig),
	)

	return err
}

func (c *provisionerClient) ProvisionFunc() interface{} {
	spec, err := c.client.ProvisionSpec(c.Ctx, &emptypb.Empty{})
	if err != nil {
		return funcErr(err)
	}
	spec.Result = nil
	cb := func(ctx context.Context, args funcspec.Args) error {
		ctx, _ = joincontext.Join(c.Ctx, ctx)
		_, err := c.client.Provision(ctx, &vagrant_plugin_sdk.FuncSpec_Args{Args: args})
		if err != nil {
			return err
		}
		return nil
	}
	return c.GenerateFunc(spec, cb)
}

func (c *provisionerClient) Provision(machine core.Machine, config *component.ConfigData) error {
	f := c.ProvisionFunc()

	_, err := c.CallDynamicFunc(f, false,
		argmapper.Typed(c.Ctx),
		argmapper.Typed(machine),
		argmapper.Typed(config),
	)

	return err
}

// provisionerServer is a gRPC server that the client talks to and calls a
// real implementation of the component.
type provisionerServer struct {
	*BaseServer

	Impl component.Provisioner
	vagrant_plugin_sdk.UnsafeProvisionerServiceServer
}

func (s *provisionerServer) CleanupSpec(
	ctx context.Context,
	args *emptypb.Empty,
) (*vagrant_plugin_sdk.FuncSpec, error) {
	if err := isImplemented(s, "provisioner"); err != nil {
		return nil, err
	}

	return s.GenerateSpec(s.Impl.CleanupFunc())
}

func (s *provisionerServer) Cleanup(
	ctx context.Context,
	args *vagrant_plugin_sdk.FuncSpec_Args,
) (*emptypb.Empty, error) {
	_, err := s.CallDynamicFunc(s.Impl.CleanupFunc(), false, args.Args,
		argmapper.Typed(ctx))

	if err != nil {
		s.Logger.Error("Error while running Cleanup", "error", err)
	}

	return &emptypb.Empty{}, err
}

func (s *provisionerServer) ConfigureSpec(
	ctx context.Context,
	args *emptypb.Empty,
) (*vagrant_plugin_sdk.FuncSpec, error) {
	if err := isImplemented(s, "provisioner"); err != nil {
		return nil, err
	}

	return s.GenerateSpec(s.Impl.ConfigureFunc())
}

func (s *provisionerServer) Configure(
	ctx context.Context,
	args *vagrant_plugin_sdk.FuncSpec_Args,
) (*emptypb.Empty, error) {
	_, err := s.CallDynamicFunc(s.Impl.ConfigureFunc(), false, args.Args,
		argmapper.Typed(ctx))

	if err != nil {
		s.Logger.Error("Error while running Configure", "error", err)
	}

	return &emptypb.Empty{}, err
}

func (s *provisionerServer) ProvisionSpec(
	ctx context.Context,
	args *emptypb.Empty,
) (*vagrant_plugin_sdk.FuncSpec, error) {
	if err := isImplemented(s, "provisioner"); err != nil {
		return nil, err
	}

	return s.GenerateSpec(s.Impl.ProvisionFunc())
}

func (s *provisionerServer) Provision(
	ctx context.Context,
	args *vagrant_plugin_sdk.FuncSpec_Args,
) (*emptypb.Empty, error) {
	_, err := s.CallDynamicFunc(s.Impl.ProvisionFunc(), false, args.Args,
		argmapper.Typed(ctx))

	if err != nil {
		s.Logger.Error("Error while running Provision", "error", err)
	}

	return &emptypb.Empty{}, err
}

var (
	_ plugin.Plugin                               = (*ProvisionerPlugin)(nil)
	_ plugin.GRPCPlugin                           = (*ProvisionerPlugin)(nil)
	_ vagrant_plugin_sdk.ProvisionerServiceServer = (*provisionerServer)(nil)
	_ component.Provisioner                       = (*provisionerClient)(nil)
	_ core.Seeder                                 = (*provisionerClient)(nil)
	_ core.Provisioner                            = (*provisionerClient)(nil)
)
