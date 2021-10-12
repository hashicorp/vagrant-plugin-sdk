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
	"github.com/hashicorp/vagrant-plugin-sdk/core"
	"github.com/hashicorp/vagrant-plugin-sdk/docs"
	"github.com/hashicorp/vagrant-plugin-sdk/internal/funcspec"
	"github.com/hashicorp/vagrant-plugin-sdk/proto/vagrant_plugin_sdk"
)

// HostPlugin implements plugin.Plugin (specifically GRPCPlugin) for
// the Host component type.
type HostPlugin struct {
	plugin.NetRPCUnsupportedPlugin

	Impl    component.Host    // Impl is the concrete implementation
	Mappers []*argmapper.Func // Mappers
	Logger  hclog.Logger      // Logger
}

func (p *HostPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	bs := &baseServer{
		base: &base{
			Mappers: p.Mappers,
			Logger:  p.Logger,
			Broker:  broker,
		},
	}
	vagrant_plugin_sdk.RegisterHostServiceServer(s, &hostServer{
		Impl:       p.Impl,
		baseServer: bs,
		capabilityServer: &capabilityServer{
			baseServer:     bs,
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
	bc := &baseClient{
		ctx: context.Background(),
		base: &base{
			Mappers: p.Mappers,
			Logger:  p.Logger,
			Broker:  broker,
		},
	}
	client := vagrant_plugin_sdk.NewHostServiceClient(c)
	return &hostClient{
		client:     client,
		baseClient: bc,
		capabilityClient: &capabilityClient{
			client:     client,
			baseClient: bc,
		},
	}, nil
}

type hostClient struct {
	*baseClient
	*capabilityClient
	client vagrant_plugin_sdk.HostServiceClient
}

func (c *hostClient) Config() (interface{}, error) {
	return configStructCall(c.ctx, c.client)
}

func (c *hostClient) ConfigSet(v interface{}) error {
	return configureCall(c.ctx, c.client, v)
}

func (c *hostClient) Documentation() (*docs.Documentation, error) {
	return documentationCall(c.ctx, c.client)
}

func (c *hostClient) HostDetectFunc() interface{} {
	spec, err := c.client.DetectSpec(c.ctx, &empty.Empty{})
	if err != nil {
		return funcErr(err)
	}
	spec.Result = nil
	cb := func(ctx context.Context, args funcspec.Args) (bool, error) {
		ctx, _ = joincontext.Join(c.ctx, ctx)
		resp, err := c.client.Detect(ctx, &vagrant_plugin_sdk.FuncSpec_Args{Args: args})
		if err != nil {
			return false, err
		}
		return resp.Detected, nil
	}

	return c.generateFunc(spec, cb)
}

func (c *hostClient) Detect(statebag core.StateBag) (bool, error) {
	f := c.HostDetectFunc()
	raw, err := c.callDynamicFunc(f, (*bool)(nil),
		argmapper.Typed(c.ctx),
		argmapper.Typed(statebag),
	)
	if err != nil {
		return false, err
	}

	return raw.(bool), nil
}

func (c *hostClient) ParentsFunc() interface{} {
	spec, err := c.client.ParentsSpec(c.ctx, &empty.Empty{})
	if err != nil {
		return funcErr(err)
	}
	spec.Result = nil
	cb := func(ctx context.Context, args funcspec.Args) ([]string, error) {
		ctx, _ = joincontext.Join(c.ctx, ctx)
		resp, err := c.client.Parents(ctx, &vagrant_plugin_sdk.FuncSpec_Args{Args: args})
		if err != nil {
			return nil, err
		}
		return resp.Parents, nil
	}

	return c.generateFunc(spec, cb)
}

func (c *hostClient) Parents() ([]string, error) {
	f := c.ParentsFunc()
	raw, err := c.callDynamicFunc(f, (*[]string)(nil),
		argmapper.Typed(c.ctx),
	)
	if err != nil {
		return nil, err
	}

	return raw.([]string), nil
}

type hostServer struct {
	*baseServer
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

	return s.generateSpec(s.Impl.HostDetectFunc())
}

func (s *hostServer) Detect(
	ctx context.Context,
	args *vagrant_plugin_sdk.FuncSpec_Args,
) (*vagrant_plugin_sdk.Platform_DetectResp, error) {
	raw, err := s.callDynamicFunc(s.Impl.HostDetectFunc(), (*bool)(nil),
		args.Args, argmapper.Typed(ctx))

	if err != nil {
		return nil, err
	}

	return &vagrant_plugin_sdk.Platform_DetectResp{
		Detected: raw.(bool)}, nil
}

func (s *hostServer) ParentsSpec(
	ctx context.Context,
	_ *empty.Empty,
) (*vagrant_plugin_sdk.FuncSpec, error) {
	if err := isImplemented(s, s.typ); err != nil {
		return nil, err
	}

	return s.generateSpec(s.Impl.ParentsFunc())
}

func (s *hostServer) Parents(
	ctx context.Context,
	args *vagrant_plugin_sdk.FuncSpec_Args,
) (*vagrant_plugin_sdk.Platform_ParentsResp, error) {
	raw, err := s.callDynamicFunc(s.Impl.ParentsFunc(), (*[]string)(nil),
		args.Args, argmapper.Typed(ctx))

	if err != nil {
		return nil, err
	}

	return &vagrant_plugin_sdk.Platform_ParentsResp{
		Parents: raw.([]string)}, nil
}

var (
	_ plugin.Plugin                        = (*HostPlugin)(nil)
	_ plugin.GRPCPlugin                    = (*HostPlugin)(nil)
	_ vagrant_plugin_sdk.HostServiceServer = (*hostServer)(nil)
	_ component.Host                       = (*hostClient)(nil)
	_ core.Host                            = (*hostClient)(nil)
)
