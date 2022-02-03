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
	"github.com/hashicorp/vagrant-plugin-sdk/docs"
	"github.com/hashicorp/vagrant-plugin-sdk/internal/funcspec"
	"github.com/hashicorp/vagrant-plugin-sdk/proto/vagrant_plugin_sdk"
)

// ProviderPlugin implements plugin.Plugin (specifically GRPCPlugin) for
// the Provider component type.
type ProviderPlugin struct {
	plugin.NetRPCUnsupportedPlugin

	Impl component.Provider // Impl is the concrete implementation
	*BasePlugin
}

func (p *ProviderPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	bs := p.NewServer(broker, p.Impl)
	vagrant_plugin_sdk.RegisterProviderServiceServer(s, &providerServer{
		Impl:       p.Impl,
		BaseServer: bs,
		capabilityServer: &capabilityServer{
			BaseServer:     bs,
			CapabilityImpl: p.Impl,
			typ:            "provider",
		},
	})
	return nil
}

func (p *ProviderPlugin) GRPCClient(
	ctx context.Context,
	broker *plugin.GRPCBroker,
	c *grpc.ClientConn,
) (interface{}, error) {
	client := vagrant_plugin_sdk.NewProviderServiceClient(c)
	bc := p.NewClient(ctx, broker, client.(SeederClient))
	return &providerClient{
		client:     client,
		BaseClient: bc,
		capabilityClient: &capabilityClient{
			client:     client,
			BaseClient: bc,
		},
	}, nil
}

// providerClient is an implementation of component.Provider over gRPC.
type providerClient struct {
	*BaseClient
	*capabilityClient
	client vagrant_plugin_sdk.ProviderServiceClient
}

func (c *providerClient) GetCapabilityClient() *capabilityClient {
	return c.capabilityClient
}

func (c *providerClient) Config() (interface{}, error) {
	return configStructCall(c.Ctx, c.client)
}

func (c *providerClient) ConfigSet(v interface{}) error {
	return configureCall(c.Ctx, c.client, v)
}

func (c *providerClient) Documentation() (*docs.Documentation, error) {
	return documentationCall(c.Ctx, c.client)
}

func (c *providerClient) UsableFunc() interface{} {
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
		return resp.IsUsable, nil
	}
	return c.GenerateFunc(spec, cb)
}

func (c *providerClient) Usable() (bool, error) {
	f := c.UsableFunc()
	raw, err := c.CallDynamicFunc(f, (*bool)(nil),
		argmapper.Typed(c.Ctx),
	)
	if err != nil {
		return false, err
	}

	return raw.(bool), nil
}

func (c *providerClient) InstalledFunc() interface{} {
	spec, err := c.client.InstalledSpec(c.Ctx, &emptypb.Empty{})
	if err != nil {
		return funcErr(err)
	}
	spec.Result = nil
	cb := func(ctx context.Context, args funcspec.Args) (bool, error) {
		ctx, _ = joincontext.Join(c.Ctx, ctx)
		resp, err := c.client.Installed(ctx, &vagrant_plugin_sdk.FuncSpec_Args{Args: args})
		if err != nil {
			return false, err
		}
		return resp.IsInstalled, nil
	}
	return c.GenerateFunc(spec, cb)
}

func (c *providerClient) Installed() (bool, error) {
	f := c.InstalledFunc()
	raw, err := c.CallDynamicFunc(f, (*bool)(nil),
		argmapper.Typed(c.Ctx),
	)
	if err != nil {
		return false, err
	}

	return raw.(bool), nil
}

func (c *providerClient) ActionFunc(name string) interface{} {
	spec, err := c.client.ActionSpec(c.Ctx, &vagrant_plugin_sdk.Provider_ActionRequest{
		Name: name,
	})
	if err != nil {
		return funcErr(err)
	}
	spec.Result = nil
	cb := func(ctx context.Context, args funcspec.Args) error {
		ctx, _ = joincontext.Join(c.Ctx, ctx)
		_, err := c.client.Action(ctx, &vagrant_plugin_sdk.Provider_ActionRequest{
			Name:     name,
			FuncArgs: &vagrant_plugin_sdk.FuncSpec_Args{Args: args},
		})
		return err
	}
	return c.GenerateFunc(spec, cb)
}

func (c *providerClient) Action(name string, args ...interface{}) error {
	f := c.ActionFunc(name)
	_, err := c.CallDynamicFunc(f, false,
		argmapper.Typed(&component.Direct{Arguments: args}),
		argmapper.Typed(args...),
		argmapper.Typed(c.Ctx),
	)
	return err
}

func (c *providerClient) MachineIdChangedFunc() interface{} {
	spec, err := c.client.MachineIdChangedSpec(c.Ctx, &emptypb.Empty{})
	if err != nil {
		return funcErr(err)
	}
	spec.Result = nil
	cb := func(ctx context.Context, args funcspec.Args) error {
		ctx, _ = joincontext.Join(c.Ctx, ctx)
		_, err := c.client.MachineIdChanged(ctx, &vagrant_plugin_sdk.FuncSpec_Args{Args: args})
		return err
	}
	return c.GenerateFunc(spec, cb)
}

func (c *providerClient) MachineIdChanged() error {
	f := c.MachineIdChangedFunc()
	_, err := c.CallDynamicFunc(f, false,
		argmapper.Typed(c.Ctx),
	)
	return err
}

func (c *providerClient) SshInfoFunc() interface{} {
	spec, err := c.client.SshInfoSpec(c.Ctx, &emptypb.Empty{})
	if err != nil {
		return funcErr(err)
	}
	spec.Result = nil
	cb := func(ctx context.Context, args funcspec.Args) (*core.SshInfo, error) {
		ctx, _ = joincontext.Join(c.Ctx, ctx)
		resp, err := c.client.SshInfo(ctx, &vagrant_plugin_sdk.FuncSpec_Args{Args: args})
		if err != nil {
			return nil, err
		}
		result, err := c.Map(resp, (**core.SshInfo)(nil),
			argmapper.Typed(c.Ctx))
		if err != nil {
			return nil, err
		}
		return result.(*core.SshInfo), nil
	}
	return c.GenerateFunc(spec, cb)
}

func (c *providerClient) SshInfo() (*core.SshInfo, error) {
	f := c.SshInfoFunc()
	raw, err := c.CallDynamicFunc(f, (**core.SshInfo)(nil),
		argmapper.Typed(c.Ctx),
	)
	if err != nil {
		return nil, err
	}

	return raw.(*core.SshInfo), nil
}

func (c *providerClient) StateFunc() interface{} {
	spec, err := c.client.StateSpec(c.Ctx, &emptypb.Empty{})
	if err != nil {
		return funcErr(err)
	}
	spec.Result = nil
	cb := func(ctx context.Context, args funcspec.Args) (*core.MachineState, error) {
		ctx, _ = joincontext.Join(c.Ctx, ctx)
		resp, err := c.client.State(ctx, &vagrant_plugin_sdk.FuncSpec_Args{Args: args})
		if err != nil {
			return nil, err
		}
		result, err := c.Map(resp, (**core.MachineState)(nil),
			argmapper.Typed(c.Ctx))
		if err != nil {
			return nil, err
		}
		return result.(*core.MachineState), nil
	}
	return c.GenerateFunc(spec, cb)
}

func (c *providerClient) State() (*core.MachineState, error) {
	f := c.StateFunc()
	raw, err := c.CallDynamicFunc(f, (**core.MachineState)(nil),
		argmapper.Typed(c.Ctx),
	)
	if err != nil {
		return nil, err
	}

	return raw.(*core.MachineState), nil
}

// providerServer is a gRPC server that the client talks to and calls a
// real implementation of the component.
type providerServer struct {
	*BaseServer
	*capabilityServer

	Impl component.Provider
}

func (s *providerServer) ConfigStruct(
	ctx context.Context,
	empty *emptypb.Empty,
) (*vagrant_plugin_sdk.Config_StructResp, error) {
	return configStruct(s.Impl)
}

func (s *providerServer) Configure(
	ctx context.Context,
	req *vagrant_plugin_sdk.Config_ConfigureRequest,
) (*emptypb.Empty, error) {
	return configure(s.Impl, req)
}

func (s *providerServer) Documentation(
	ctx context.Context,
	empty *emptypb.Empty,
) (*vagrant_plugin_sdk.Config_Documentation, error) {
	return documentation(s.Impl)
}

func (s *providerServer) UsableSpec(
	ctx context.Context,
	args *emptypb.Empty,
) (*vagrant_plugin_sdk.FuncSpec, error) {
	if err := isImplemented(s, "provider"); err != nil {
		return nil, err
	}

	return s.GenerateSpec(s.Impl.UsableFunc())
}

func (s *providerServer) Usable(
	ctx context.Context,
	args *vagrant_plugin_sdk.FuncSpec_Args,
) (*vagrant_plugin_sdk.Provider_UsableResp, error) {
	raw, err := s.CallDynamicFunc(s.Impl.UsableFunc(), (*bool)(nil),
		args.Args, argmapper.Typed(ctx))

	if err != nil {
		return nil, err
	}

	return &vagrant_plugin_sdk.Provider_UsableResp{
		IsUsable: raw.(bool)}, nil
}

func (s *providerServer) InstalledSpec(
	ctx context.Context,
	args *emptypb.Empty,
) (*vagrant_plugin_sdk.FuncSpec, error) {
	if err := isImplemented(s, "provider"); err != nil {
		return nil, err
	}

	return s.GenerateSpec(s.Impl.InstalledFunc())
}

func (s *providerServer) Installed(
	ctx context.Context,
	args *vagrant_plugin_sdk.FuncSpec_Args,
) (*vagrant_plugin_sdk.Provider_InstalledResp, error) {
	raw, err := s.CallDynamicFunc(s.Impl.InstalledFunc(), (*bool)(nil),
		args.Args, argmapper.Typed(ctx))

	if err != nil {
		return nil, err
	}

	return &vagrant_plugin_sdk.Provider_InstalledResp{
		IsInstalled: raw.(bool)}, nil
}

func (s *providerServer) ActionSpec(
	ctx context.Context,
	args *vagrant_plugin_sdk.Provider_ActionRequest,
) (*vagrant_plugin_sdk.FuncSpec, error) {
	if err := isImplemented(s, "provider"); err != nil {
		return nil, err
	}

	return s.GenerateSpec(s.Impl.ActionFunc(args.Name))
}

func (s *providerServer) Action(
	ctx context.Context,
	args *vagrant_plugin_sdk.Provider_ActionRequest,
) (*emptypb.Empty, error) {
	_, err := s.CallDynamicFunc(
		s.Impl.ActionFunc(args.Name),
		false,
		args.FuncArgs.Args,
		argmapper.Typed(ctx),
	)

	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *providerServer) MachineIdChangedSpec(
	ctx context.Context,
	args *emptypb.Empty,
) (*vagrant_plugin_sdk.FuncSpec, error) {
	if err := isImplemented(s, "provider"); err != nil {
		return nil, err
	}

	return s.GenerateSpec(s.Impl.MachineIdChangedFunc())
}

func (s *providerServer) MachineIdChanged(
	ctx context.Context,
	args *vagrant_plugin_sdk.FuncSpec_Args,
) (*emptypb.Empty, error) {
	_, err := s.CallDynamicFunc(s.Impl.MachineIdChangedFunc(), false,
		args.Args, argmapper.Typed(ctx))

	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *providerServer) SshInfoSpec(
	ctx context.Context,
	args *emptypb.Empty,
) (*vagrant_plugin_sdk.FuncSpec, error) {
	if err := isImplemented(s, "provider"); err != nil {
		return nil, err
	}

	return s.GenerateSpec(s.Impl.SshInfoFunc())
}

func (s *providerServer) SshInfo(
	ctx context.Context,
	args *vagrant_plugin_sdk.FuncSpec_Args,
) (*vagrant_plugin_sdk.SSHInfo, error) {
	raw, err := s.CallDynamicFunc(s.Impl.SshInfoFunc(), (**core.SshInfo)(nil),
		args.Args, argmapper.Typed(ctx))

	if err != nil {
		return nil, err
	}

	result, err := s.Map(
		raw,
		(**vagrant_plugin_sdk.SSHInfo)(nil),
		argmapper.Typed(ctx),
	)
	if err != nil {
		return nil, err
	}

	return result.(*vagrant_plugin_sdk.SSHInfo), nil
}

func (s *providerServer) StateSpec(
	ctx context.Context,
	args *emptypb.Empty,
) (*vagrant_plugin_sdk.FuncSpec, error) {
	if err := isImplemented(s, "provider"); err != nil {
		return nil, err
	}

	return s.GenerateSpec(s.Impl.StateFunc())
}

func (s *providerServer) State(
	ctx context.Context,
	args *vagrant_plugin_sdk.FuncSpec_Args,
) (*vagrant_plugin_sdk.Args_Target_Machine_State, error) {
	raw, err := s.CallDynamicFunc(s.Impl.StateFunc(), (**core.MachineState)(nil),
		args.Args, argmapper.Typed(ctx))

	if err != nil {
		return nil, err
	}

	result, err := s.Map(
		raw,
		(**vagrant_plugin_sdk.Args_Target_Machine_State)(nil),
		argmapper.Typed(ctx),
	)
	if err != nil {
		return nil, err
	}

	return result.(*vagrant_plugin_sdk.Args_Target_Machine_State), nil
}

var (
	_ plugin.Plugin                            = (*ProviderPlugin)(nil)
	_ plugin.GRPCPlugin                        = (*ProviderPlugin)(nil)
	_ vagrant_plugin_sdk.ProviderServiceServer = (*providerServer)(nil)
	_ component.Provider                       = (*providerClient)(nil)
	_ core.Provider                            = (*providerClient)(nil)
	_ component.CapabilityPlatform             = (*providerClient)(nil)
	_ core.CapabilityPlatform                  = (*providerClient)(nil)
	_ component.Configurable                   = (*providerClient)(nil)
	_ component.Documented                     = (*providerClient)(nil)
	_ core.Seeder                              = (*providerClient)(nil)
)
