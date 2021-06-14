package component

import (
	"context"

	"github.com/LK4D4/joincontext"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/hashicorp/go-argmapper"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/hashicorp/vagrant-plugin-sdk/core"

	"github.com/hashicorp/vagrant-plugin-sdk/component"
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
	vagrant_plugin_sdk.RegisterHostServiceServer(s, &hostServer{
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

func (p *HostPlugin) GRPCClient(
	ctx context.Context,
	broker *plugin.GRPCBroker,
	c *grpc.ClientConn,
) (interface{}, error) {
	return &hostClient{
		client: vagrant_plugin_sdk.NewHostServiceClient(c),
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

// hostClient is an implementation of component.Host over gRPC.
type hostClient struct {
	*baseClient

	client vagrant_plugin_sdk.HostServiceClient
}

func (c *hostClient) Config() (v interface{}, err error) {
	return
}

func (c *hostClient) ConfigSet(v interface{}) (err error) {
	return
}

func (c *hostClient) Documentation() (d *docs.Documentation, err error) {
	return
}

func (c *hostClient) DetectFunc() interface{} {
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

func (c *hostClient) Detect() bool {
	f := c.DetectFunc()
	raw, err := c.callRemoteDynamicFunc(c.ctx, nil, (*bool)(nil), f)
	if err != nil {
		return false
	}

	return raw.(bool)
}

func (c *hostClient) HasCapabilityFunc() interface{} {
	spec, err := c.client.HasCapabilitySpec(c.ctx, &empty.Empty{})
	if err != nil {
		// TODO: stay calm, don't panic
		panic(err)
		// return funcErr(err)
	}
	spec.Result = nil
	cb := func(ctx context.Context, args funcspec.Args, capabilityName string) (bool, error) {
		resp, err := c.client.HasCapability(ctx, &vagrant_plugin_sdk.Host_Capability_NamedRequest{
			FuncArgs: &vagrant_plugin_sdk.FuncSpec_Args{Args: args},
			Name:     capabilityName,
		})
		if err != nil {
			return false, err
		}

		return resp.HasCapability, nil
	}
	f := c.generateFunc(spec, cb)
	return f
}

func (c *hostClient) HasCapability(name string) bool {
	f := c.HasCapabilityFunc()
	raw, err := c.callRemoteDynamicFunc(
		c.ctx,
		c.Mappers,
		(*bool)(nil),
		f,
		argmapper.Typed(name),
	)
	if err != nil {
		return false
	}

	return raw.(bool)
}

func (c *hostClient) CapabilityFunc(capName string) interface{} {
	spec, err := c.client.CapabilitySpec(c.ctx, &vagrant_plugin_sdk.Host_Capability_NamedRequest{Name: capName})
	if err != nil {
		// TODO: stay calm, don't panic
		panic(err)
		// return funcErr(err)
	}
	spec.Result = nil
	cb := func(ctx context.Context, args funcspec.Args) (*anypb.Any, error) {
		ctx, _ = joincontext.Join(c.ctx, ctx)
		capArgs := &vagrant_plugin_sdk.Host_Capability_NamedRequest{
			FuncArgs: &vagrant_plugin_sdk.FuncSpec_Args{Args: args},
			Name:     capName,
		}
		resp, err := c.client.Capability(ctx, capArgs)
		if err != nil {
			return nil, err
		}
		return resp.Result, nil
	}
	return c.generateFunc(spec, cb)
}

func (c *hostClient) Capability(name string, args ...argmapper.Arg) (interface{}, error) {
	f := c.CapabilityFunc(name)
	raw, err := c.callRemoteDynamicFunc(
		c.ctx,
		c.Mappers,
		(interface{})(nil),
		f,
		args...,
	)
	if err != nil {
		return false, err
	}

	return raw, nil
}

type hostServer struct {
	*baseServer

	Impl component.Host
	vagrant_plugin_sdk.UnimplementedHostServiceServer
}

func (s *hostServer) ConfigStruct(
	ctx context.Context,
	empty *empty.Empty,
) (c *vagrant_plugin_sdk.Config_StructResp, err error) {
	return
}

func (s *hostServer) Configure(
	ctx context.Context,
	req *vagrant_plugin_sdk.Config_ConfigureRequest,
) (e *empty.Empty, err error) {
	return
}

func (s *hostServer) Documentation(
	ctx context.Context,
	empty *empty.Empty,
) (c *vagrant_plugin_sdk.Config_Documentation, err error) {
	return
}

func (s *hostServer) DetectSpec(
	ctx context.Context,
	args *empty.Empty,
) (*vagrant_plugin_sdk.FuncSpec, error) {
	if err := isImplemented(s, "host"); err != nil {
		return nil, err
	}

	return s.generateSpec(s.Impl.DetectFunc())
}

func (s *hostServer) Detect(
	ctx context.Context,
	args *vagrant_plugin_sdk.FuncSpec_Args,
) (*vagrant_plugin_sdk.Host_DetectResp, error) {
	s.Logger.Debug("running the detect function on the server to call real implementation")
	raw, err := s.callLocalDynamicFunc(s.Impl.DetectFunc(), args.Args, (*bool)(nil),
		argmapper.Typed(ctx),
	)

	if err != nil {
		return nil, err
	}

	return &vagrant_plugin_sdk.Host_DetectResp{Detected: raw.(bool)}, nil
}

func (s *hostServer) HasCapabilitySpec(
	ctx context.Context,
	args *empty.Empty,
) (*vagrant_plugin_sdk.FuncSpec, error) {
	if err := isImplemented(s, "host"); err != nil {
		return nil, err
	}

	f := s.Impl.HasCapability
	return s.generateSpec(f)
}

func (s *hostServer) HasCapability(
	ctx context.Context,
	args *vagrant_plugin_sdk.Host_Capability_NamedRequest,
) (*vagrant_plugin_sdk.Host_Capability_CheckResp, error) {
	f := s.Impl.HasCapability

	raw, err := s.callUncheckedLocalDynamicFunc(
		f,
		args.FuncArgs.Args,
		argmapper.Typed(ctx),
		argmapper.Typed(args.Name),
	)

	if err != nil {
		return nil, err
	}
	return &vagrant_plugin_sdk.Host_Capability_CheckResp{HasCapability: raw.(bool)}, err
}

func (s *hostServer) CapabilitySpec(
	ctx context.Context,
	args *vagrant_plugin_sdk.Host_Capability_NamedRequest,
) (*vagrant_plugin_sdk.FuncSpec, error) {
	if err := isImplemented(s, "host"); err != nil {
		return nil, err
	}

	_, ok := s.Impl.(*hostClient)
	if ok {
		spec := s.Impl.CapabilityFunc(args.Name)
		g, e := s.generateArgSpec(spec.(*argmapper.Func))
		return g, e
	}
	spec := s.Impl.CapabilityFunc(args.Name)
	g, e := s.generateSpec(spec)
	return g, e
}

func (s *hostServer) Capability(
	ctx context.Context,
	args *vagrant_plugin_sdk.Host_Capability_NamedRequest,
) (*vagrant_plugin_sdk.Host_Capability_Resp, error) {
	_, ok := s.Impl.(*hostClient)
	if ok {
		fn := s.Impl.CapabilityFunc(args.Name)
		raw, err := s.callUncheckedLocalDynamicArgmapperFunc(
			fn.(*argmapper.Func),
			args.FuncArgs.Args,
			argmapper.Typed(ctx),
			argmapper.Typed(args.Name),
		)
		if err != nil {
			return nil, err
		}
		if raw != nil {
			return &vagrant_plugin_sdk.Host_Capability_Resp{Result: raw.(*anypb.Any)}, nil
		} else {
			return &vagrant_plugin_sdk.Host_Capability_Resp{Result: nil}, nil
		}
	} else {
		fn := s.Impl.CapabilityFunc(args.Name)
		raw, err := s.callUncheckedLocalDynamicFunc(fn, args.FuncArgs.Args, argmapper.Typed(ctx))
		if err != nil {
			return nil, err
		}
		if raw != nil {
			return &vagrant_plugin_sdk.Host_Capability_Resp{Result: raw.(*anypb.Any)}, nil
		} else {
			return &vagrant_plugin_sdk.Host_Capability_Resp{Result: nil}, nil
		}
	}
}

var (
	_ plugin.Plugin                        = (*HostPlugin)(nil)
	_ plugin.GRPCPlugin                    = (*HostPlugin)(nil)
	_ vagrant_plugin_sdk.HostServiceServer = (*hostServer)(nil)
	_ component.Host                       = (*hostClient)(nil)
	_ core.Host                            = (*hostClient)(nil)
)
