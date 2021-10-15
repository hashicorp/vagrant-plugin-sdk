package plugin

import (
	"context"
	"fmt"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/hashicorp/go-argmapper"
	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"

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
	vagrant_plugin_sdk.RegisterProviderServiceServer(s, &providerServer{
		Impl:       p.Impl,
		BaseServer: p.NewServer(broker),
	})
	return nil
}

func (p *ProviderPlugin) GRPCClient(
	ctx context.Context,
	broker *plugin.GRPCBroker,
	c *grpc.ClientConn,
) (interface{}, error) {
	return &providerClient{
		client:     vagrant_plugin_sdk.NewProviderServiceClient(c),
		BaseClient: p.NewClient(ctx, broker),
	}, nil
}

// providerClient is an implementation of component.Provider over gRPC.
type providerClient struct {
	*BaseClient

	client vagrant_plugin_sdk.ProviderServiceClient
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
	spec, err := c.client.UsableSpec(c.Ctx, &empty.Empty{})
	if err != nil {
		return funcErr(err)
	}
	spec.Result = nil
	cb := func(ctx context.Context, args funcspec.Args) (bool, error) {
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

func (c *providerClient) InitFunc() interface{} {
	spec, err := c.client.InitSpec(c.Ctx, &empty.Empty{})
	if err != nil {
		return funcErr(err)
	}
	spec.Result = nil
	cb := func(ctx context.Context, args funcspec.Args) (bool, error) {
		_, err := c.client.Init(ctx, &vagrant_plugin_sdk.FuncSpec_Args{Args: args})
		if err != nil {
			return false, err
		}
		return true, nil
	}
	return c.GenerateFunc(spec, cb)
}

func (c *providerClient) Init(machine core.Machine) (bool, error) {
	f := c.InitFunc()
	_, err := c.CallDynamicFunc(f, (*bool)(nil),
		argmapper.Typed(machine),
		argmapper.Typed(c.Ctx),
	)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (c *providerClient) InstalledFunc() interface{} {
	spec, err := c.client.InstalledSpec(c.Ctx, &empty.Empty{})
	if err != nil {
		return funcErr(err)
	}
	spec.Result = nil
	cb := func(ctx context.Context, args funcspec.Args) (bool, error) {
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

func (c *providerClient) ActionUpFunc() interface{} {
	spec, err := c.client.ActionUpSpec(c.Ctx, &empty.Empty{})
	if err != nil {
		return funcErr(err)
	}
	spec.Result = nil
	cb := func(ctx context.Context, args funcspec.Args) (interface{}, error) {
		resp, err := c.client.ActionUp(ctx, &vagrant_plugin_sdk.FuncSpec_Args{Args: args})
		if err != nil {
			return false, err
		}
		return resp.Result, nil
	}
	return c.GenerateFunc(spec, cb)
}

func (c *providerClient) ActionUp() error {
	f := c.ActionUpFunc()
	_, err := c.CallDynamicFunc(f, false,
		argmapper.Typed(c.Ctx),
	)
	if err != nil {
		return err
	}

	return nil
}

// providerServer is a gRPC server that the client talks to and calls a
// real implementation of the component.
type providerServer struct {
	*BaseServer

	Impl component.Provider
	vagrant_plugin_sdk.UnimplementedProviderServiceServer
}

func (s *providerServer) ConfigStruct(
	ctx context.Context,
	empty *empty.Empty,
) (*vagrant_plugin_sdk.Config_StructResp, error) {
	return configStruct(s.Impl)
}

func (s *providerServer) Configure(
	ctx context.Context,
	req *vagrant_plugin_sdk.Config_ConfigureRequest,
) (*empty.Empty, error) {
	return configure(s.Impl, req)
}

func (s *providerServer) Documentation(
	ctx context.Context,
	empty *empty.Empty,
) (*vagrant_plugin_sdk.Config_Documentation, error) {
	return documentation(s.Impl)
}

func (s *providerServer) UsableSpec(
	ctx context.Context,
	args *empty.Empty,
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
	args *empty.Empty,
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

func (s *providerServer) ActionUpSpec(
	ctx context.Context,
	args *empty.Empty,
) (*vagrant_plugin_sdk.FuncSpec, error) {
	if err := isImplemented(s, "provider"); err != nil {
		return nil, err
	}

	return s.GenerateSpec(s.Impl.ActionUpFunc())
}

func (s *providerServer) ActionUp(
	ctx context.Context,
	args *vagrant_plugin_sdk.FuncSpec_Args,
) (*vagrant_plugin_sdk.Provider_ActionResp, error) {
	raw, err := s.CallDynamicFunc(s.Impl.ActionUpFunc(), (*proto.Message)(nil),
		args.Args, argmapper.Typed(ctx))

	if err != nil {
		return nil, err
	}
	// Expect the results to be proto.Messages
	msg, ok := raw.(proto.Message)
	if !ok {
		return nil, fmt.Errorf(
			"result of plugin-based function must be a proto.Message, got %T", msg)
	}
	anyVal, err := ptypes.MarshalAny(msg)

	// TODO: This maybe needs to be expanded
	return &vagrant_plugin_sdk.Provider_ActionResp{Result: anyVal}, nil
}

func (s *providerServer) InitSpec(
	ctx context.Context,
	args *empty.Empty,
) (*vagrant_plugin_sdk.FuncSpec, error) {
	if err := isImplemented(s, "provider"); err != nil {
		return nil, err
	}

	return s.GenerateSpec(s.Impl.InitFunc())
}

func (s *providerServer) Init(
	ctx context.Context,
	args *vagrant_plugin_sdk.FuncSpec_Args,
) (*empty.Empty, error) {
	_, err := s.CallDynamicFunc(s.Impl.InitFunc(), false,
		args.Args, argmapper.Typed(ctx))

	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

var (
	_ plugin.Plugin                            = (*ProviderPlugin)(nil)
	_ plugin.GRPCPlugin                        = (*ProviderPlugin)(nil)
	_ vagrant_plugin_sdk.ProviderServiceServer = (*providerServer)(nil)
	_ component.Provider                       = (*providerClient)(nil)
	_ component.Configurable                   = (*providerClient)(nil)
	_ component.Documented                     = (*providerClient)(nil)
)
