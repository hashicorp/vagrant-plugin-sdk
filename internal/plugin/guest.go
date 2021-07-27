package plugin

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/hashicorp/go-argmapper"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"

	"github.com/hashicorp/vagrant-plugin-sdk/component"
	"github.com/hashicorp/vagrant-plugin-sdk/docs"
	"github.com/hashicorp/vagrant-plugin-sdk/internal/funcspec"
	"github.com/hashicorp/vagrant-plugin-sdk/proto/vagrant_plugin_sdk"
)

// GuestPlugin implements plugin.Plugin (specifically GRPCPlugin) for
// the Guest component type.
type GuestPlugin struct {
	plugin.NetRPCUnsupportedPlugin

	Impl    component.Guest   // Impl is the concrete implementation
	Mappers []*argmapper.Func // Mappers
	Logger  hclog.Logger      // Logger
}

func (p *GuestPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	bs := &baseServer{
		base: &base{
			Mappers: p.Mappers,
			Logger:  p.Logger,
			Broker:  broker,
		},
	}
	vagrant_plugin_sdk.RegisterGuestServiceServer(s, &hostServer{
		Impl:       p.Impl,
		baseServer: bs,
		capabilityServer: &capabilityServer{
			baseServer:     bs,
			CapabilityImpl: p.Impl,
			typ:            "guest",
		},
	})
	return nil
}

func (p *GuestPlugin) GRPCClient(
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
	client := vagrant_plugin_sdk.NewGuestServiceClient(c)
	return &hostClient{
		client:     client,
		baseClient: bc,
		capabilityClient: &capabilityClient{
			client:     client,
			baseClient: bc,
		},
	}, nil
}

// guestClient is an implementation of component.Guest over gRPC.
type guestClient struct {
	*baseClient
	*capabilityClient
	client vagrant_plugin_sdk.GuestServiceClient
}

func (c *guestClient) Config() (interface{}, error) {
	return configStructCall(c.ctx, c.client)
}

func (c *guestClient) ConfigSet(v interface{}) error {
	return configureCall(c.ctx, c.client, v)
}

func (c *guestClient) Documentation() (*docs.Documentation, error) {
	return documentationCall(c.ctx, c.client)
}

func (c *guestClient) DetectFunc() interface{} {
	spec, err := c.client.DetectSpec(c.ctx, &empty.Empty{})
	if err != nil {
		return funcErr(err)
	}
	spec.Result = nil
	cb := func(ctx context.Context, args funcspec.Args) (bool, error) {
		resp, err := c.client.Detect(ctx, &vagrant_plugin_sdk.FuncSpec_Args{Args: args})
		if err != nil {
			return false, err
		}
		return resp.Detected, nil
	}
	return c.generateFunc(spec, cb)
}

func (c *guestClient) Detect() (bool, error) {
	f := c.DetectFunc()
	raw, err := c.callDynamicFunc(f, (*bool)(nil),
		argmapper.Typed(c.ctx),
	)
	if err != nil {
		return false, err
	}

	return raw.(bool), nil
}

// guestServer is a gRPC server that the client talks to and calls a
// real implementation of the component.
type guestServer struct {
	*baseServer
	*capabilityServer

	Impl component.Guest
}

func (s *guestServer) ConfigStruct(
	ctx context.Context,
	empty *empty.Empty,
) (*vagrant_plugin_sdk.Config_StructResp, error) {
	return configStruct(s.Impl)
}

func (s *guestServer) Configure(
	ctx context.Context,
	req *vagrant_plugin_sdk.Config_ConfigureRequest,
) (*empty.Empty, error) {
	return configure(s.Impl, req)
}

func (s *guestServer) Documentation(
	ctx context.Context,
	empty *empty.Empty,
) (*vagrant_plugin_sdk.Config_Documentation, error) {
	return documentation(s.Impl)
}

func (s *guestServer) DetectSpec(
	ctx context.Context,
	args *empty.Empty,
) (*vagrant_plugin_sdk.FuncSpec, error) {
	if err := isImplemented(s, "guest"); err != nil {
		return nil, err
	}

	return s.generateSpec(s.Impl.DetectFunc())
}

func (s *guestServer) Detect(
	ctx context.Context,
	args *vagrant_plugin_sdk.FuncSpec_Args,
) (*vagrant_plugin_sdk.Platform_DetectResp, error) {
	raw, err := s.callDynamicFunc(s.Impl.DetectFunc(), (*bool)(nil), args.Args,
		argmapper.Typed(ctx),
	)

	if err != nil {
		return nil, err
	}

	return &vagrant_plugin_sdk.Platform_DetectResp{Detected: raw.(bool)}, nil
}

var (
	_ plugin.Plugin                         = (*GuestPlugin)(nil)
	_ plugin.GRPCPlugin                     = (*GuestPlugin)(nil)
	_ vagrant_plugin_sdk.GuestServiceServer = (*guestServer)(nil)
	_ component.Guest                       = (*guestClient)(nil)
)
