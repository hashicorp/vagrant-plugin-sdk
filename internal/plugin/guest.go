package plugin

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/hashicorp/go-argmapper"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/hashicorp/vagrant-plugin-sdk/component"
	"github.com/hashicorp/vagrant-plugin-sdk/core"
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
	vagrant_plugin_sdk.RegisterGuestServiceServer(s, &guestServer{
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

func (p *GuestPlugin) GRPCClient(
	ctx context.Context,
	broker *plugin.GRPCBroker,
	c *grpc.ClientConn,
) (interface{}, error) {
	return &guestClient{
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

// guestClient is an implementation of component.Guest over gRPC.
type guestClient struct {
	*baseClient

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

func (c *guestClient) Detect(machine core.Machine) (bool, error) {
	f := c.DetectFunc()
	raw, err := c.callRemoteDynamicFunc(c.ctx, nil, (*bool)(nil), f,
		argmapper.Typed(machine))
	if err != nil {
		return false, err
	}

	return raw.(bool), nil
}

func (c *guestClient) HasCapabilityFunc(capName string) interface{} {
	spec, err := c.client.HasCapabilitySpec(c.ctx, &vagrant_plugin_sdk.Guest_Capability_NamedRequest{Name: capName})
	if err != nil {
		return funcErr(err)
	}
	spec.Result = nil
	cb := func(ctx context.Context, args funcspec.Args) (bool, error) {
		resp, err := c.client.HasCapability(ctx, &vagrant_plugin_sdk.Guest_Capability_NamedRequest{Name: capName, FuncArgs: &vagrant_plugin_sdk.FuncSpec_Args{Args: args}})
		if err != nil {
			return false, err
		}
		return resp.HasCapability, nil
	}
	return c.generateFunc(spec, cb)
}

func (c *guestClient) HasCapability(machine core.Machine, capName string) (bool, error) {
	f := c.HasCapabilityFunc(capName)
	raw, err := c.callRemoteDynamicFunc(c.ctx, nil, (*bool)(nil), f,
		argmapper.Typed(machine),
		argmapper.Typed(capName),
		argmapper.Named("capabilityName", capName),
	)
	if err != nil {
		return false, err
	}

	return raw.(bool), nil
}

func (c *guestClient) CapabilityFunc(capName string) interface{} {
	spec, err := c.client.CapabilitySpec(c.ctx, &vagrant_plugin_sdk.Guest_Capability_NamedRequest{Name: capName})
	if err != nil {
		return funcErr(err)
	}
	spec.Result = nil
	cb := func(ctx context.Context, args funcspec.Args) (interface{}, error) {
		resp, err := c.client.Capability(ctx, &vagrant_plugin_sdk.Guest_Capability_NamedRequest{Name: capName, FuncArgs: &vagrant_plugin_sdk.FuncSpec_Args{Args: args}})
		if err != nil {
			return nil, err
		}
		return resp.Result, nil
	}
	return c.generateFunc(spec, cb)
}

// TODO(spox): need to determine what we want to do here with regards to cap results
func (c *guestClient) Capability(machine core.Machine, capName string, args ...interface{}) (interface{}, error) {
	f := c.CapabilityFunc(capName)
	margs := []argmapper.Arg{
		argmapper.Typed(machine),
	}
	for _, a := range args {
		margs = append(margs, argmapper.Typed(a))
	}
	raw, err := c.callRemoteDynamicFunc(c.ctx, nil, (interface{})(nil), f, margs...)
	if err != nil {
		return nil, err
	}

	return raw, nil
}

// guestServer is a gRPC server that the client talks to and calls a
// real implementation of the component.
type guestServer struct {
	*baseServer

	Impl component.Guest
	vagrant_plugin_sdk.UnimplementedGuestServiceServer
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
) (*vagrant_plugin_sdk.Guest_DetectResp, error) {
	raw, err := s.callLocalDynamicFunc(s.Impl.DetectFunc(), args.Args, (*bool)(nil),
		argmapper.Typed(ctx),
	)

	if err != nil {
		return nil, err
	}

	return &vagrant_plugin_sdk.Guest_DetectResp{Detected: raw.(bool)}, nil
}

func (s *guestServer) HasCapabilitySpec(
	ctx context.Context,
	args *vagrant_plugin_sdk.Guest_Capability_NamedRequest,
) (*vagrant_plugin_sdk.FuncSpec, error) {
	if err := isImplemented(s, "guest"); err != nil {
		return nil, err
	}

	return s.generateSpec(s.Impl.HasCapabilityFunc(args.Name))
}

func (s *guestServer) HasCapability(
	ctx context.Context,
	args *vagrant_plugin_sdk.Guest_Capability_NamedRequest,
) (*vagrant_plugin_sdk.Guest_Capability_CheckResp, error) {
	raw, err := s.callLocalDynamicFunc(s.Impl.HasCapabilityFunc(args.Name), args.FuncArgs.Args, (*bool)(nil),
		argmapper.Typed(ctx),
	)

	if err != nil {
		return nil, err
	}

	return &vagrant_plugin_sdk.Guest_Capability_CheckResp{HasCapability: raw.(bool)}, nil
}

func (s *guestServer) CapabilitySpec(
	ctx context.Context,
	args *vagrant_plugin_sdk.Guest_Capability_NamedRequest,
) (*vagrant_plugin_sdk.FuncSpec, error) {
	if err := isImplemented(s, "guest"); err != nil {
		return nil, err
	}

	return s.generateSpec(s.Impl.CapabilityFunc(args.Name))
}

func (s *guestServer) Capability(
	ctx context.Context,
	args *vagrant_plugin_sdk.Guest_Capability_NamedRequest,
) (*vagrant_plugin_sdk.Guest_Capability_Resp, error) {
	raw, err := s.callLocalDynamicFunc(s.Impl.CapabilityFunc(args.Name), args.FuncArgs.Args, (interface{})(nil),
		argmapper.Typed(ctx),
	)

	if err != nil {
		return nil, err
	}

	return &vagrant_plugin_sdk.Guest_Capability_Resp{Result: raw.(*anypb.Any)}, nil
}

var (
	_ plugin.Plugin                         = (*GuestPlugin)(nil)
	_ plugin.GRPCPlugin                     = (*GuestPlugin)(nil)
	_ vagrant_plugin_sdk.GuestServiceServer = (*guestServer)(nil)
	_ component.Guest                       = (*guestClient)(nil)
)
