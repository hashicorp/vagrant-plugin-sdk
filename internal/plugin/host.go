package plugin

import (
	"context"
	"errors"

	"github.com/LK4D4/joincontext"
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

// HostPlugin implements plugin.Plugin (specifically GRPCPlugin) for
// the Host component type.
type HostPlugin struct {
	plugin.NetRPCUnsupportedPlugin

	Impl component.Host // Impl is the concrete implementation
	*BasePlugin
}

func (p *HostPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	bs := p.NewServer(broker)
	vagrant_plugin_sdk.RegisterHostServiceServer(s, &hostServer{
		Impl:       p.Impl,
		BaseServer: bs,
		capabilityServer: &capabilityServer{
			BaseServer:     bs,
			CapabilityImpl: p.Impl,
			typ:            "host",
		},
	})
	return nil
}

func (p *HostPlugin) GRPCClient(
	ctx context.Context,
	broker *plugin.GRPCBroker,
	c *grpc.ClientConn,
) (interface{}, error) {
	bc := p.NewClient(ctx, broker)
	client := vagrant_plugin_sdk.NewHostServiceClient(c)
	return &hostClient{
		client:     client,
		BaseClient: bc,
		capabilityClient: &capabilityClient{
			client:     client,
			BaseClient: bc,
		},
	}, nil
}

type hostClient struct {
	*BaseClient
	*capabilityClient
	client vagrant_plugin_sdk.HostServiceClient
}

func (c *hostClient) GetCapabilityClient() *capabilityClient {
	return c.capabilityClient
}

func (c *hostClient) Config() (interface{}, error) {
	return configStructCall(c.Ctx, c.client)
}

func (c *hostClient) ConfigSet(v interface{}) error {
	return configureCall(c.Ctx, c.client, v)
}

func (c *hostClient) Documentation() (*docs.Documentation, error) {
	return documentationCall(c.Ctx, c.client)
}

func (c *hostClient) HostDetectFunc() interface{} {
	spec, err := c.client.DetectSpec(c.Ctx, &empty.Empty{})
	if err != nil {
		return funcErr(err)
	}
	spec.Result = nil
	cb := func(ctx context.Context, args funcspec.Args) (bool, error) {
		ctx, _ = joincontext.Join(c.Ctx, ctx)
		resp, err := c.client.Detect(ctx, &vagrant_plugin_sdk.FuncSpec_Args{Args: args})
		if err != nil {
			return false, err
		}
		return resp.Detected, nil
	}

	return c.GenerateFunc(spec, cb)
}

func (c *hostClient) Detect(statebag core.StateBag) (bool, error) {
	f := c.HostDetectFunc()
	raw, err := c.CallDynamicFunc(f, (*bool)(nil),
		argmapper.Typed(c.Ctx),
		argmapper.Typed(statebag),
	)
	if err != nil {
		return false, err
	}

	return raw.(bool), nil
}

func (c *hostClient) ParentFunc() interface{} {
	spec, err := c.client.ParentSpec(c.Ctx, &empty.Empty{})
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

func (c *hostClient) Parent() (string, error) {
	f := c.ParentFunc()
	raw, err := c.CallDynamicFunc(f, (*string)(nil),
		argmapper.Typed(c.Ctx),
	)
	if err != nil {
		return "", err
	}

	return raw.(string), nil
}

type hostServer struct {
	*BaseServer
	*capabilityServer

	Impl component.Host
}

func (s *hostServer) ConfigStruct(
	ctx context.Context,
	empty *empty.Empty,
) (*vagrant_plugin_sdk.Config_StructResp, error) {
	return configStruct(s.Impl)
}

func (s *hostServer) Configure(
	ctx context.Context,
	req *vagrant_plugin_sdk.Config_ConfigureRequest,
) (*empty.Empty, error) {
	return configure(s.Impl, req)
}

func (s *hostServer) Documentation(
	ctx context.Context,
	empty *empty.Empty,
) (*vagrant_plugin_sdk.Config_Documentation, error) {
	return documentation(s.Impl)
}

func (s *hostServer) DetectSpec(
	ctx context.Context,
	_ *empty.Empty,
) (*vagrant_plugin_sdk.FuncSpec, error) {
	if err := isImplemented(s, "host"); err != nil {
		return nil, err
	}

	if err := isImplemented(s.Impl, "host"); err != nil {
		return nil, err
	}

	return s.GenerateSpec(s.Impl.HostDetectFunc())
}

func (s *hostServer) Detect(
	ctx context.Context,
	args *vagrant_plugin_sdk.FuncSpec_Args,
) (*vagrant_plugin_sdk.Platform_DetectResp, error) {
	raw, err := s.CallDynamicFunc(s.Impl.HostDetectFunc(), (*bool)(nil),
		args.Args, argmapper.Typed(ctx))

	if err != nil {
		return nil, err
	}

	return &vagrant_plugin_sdk.Platform_DetectResp{
		Detected: raw.(bool)}, nil
}

func (s *hostServer) ParentSpec(
	ctx context.Context,
	_ *empty.Empty,
) (*vagrant_plugin_sdk.FuncSpec, error) {
	if err := isImplemented(s, s.typ); err != nil {
		return nil, errors.New("error with is implemented")
	}

	f, err := s.GenerateSpec(s.Impl.ParentFunc())
	if err != nil {
		return nil, errors.New("generating spec")
	}
	return f, err
}

func (s *hostServer) Parent(
	ctx context.Context,
	args *vagrant_plugin_sdk.FuncSpec_Args,
) (*vagrant_plugin_sdk.Platform_ParentResp, error) {
	raw, err := s.CallDynamicFunc(s.Impl.ParentFunc(), (*string)(nil),
		args.Args, argmapper.Typed(ctx))

	if err != nil {
		return nil, err
	}

	return &vagrant_plugin_sdk.Platform_ParentResp{
		Parent: raw.(string)}, nil
}

var (
	_ plugin.Plugin                        = (*HostPlugin)(nil)
	_ plugin.GRPCPlugin                    = (*HostPlugin)(nil)
	_ vagrant_plugin_sdk.HostServiceServer = (*hostServer)(nil)
	_ component.Host                       = (*hostClient)(nil)
	_ core.Host                            = (*hostClient)(nil)
	_ capabilityParent                     = (*hostClient)(nil)
)
