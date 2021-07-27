package plugin

import (
	"context"

	"github.com/LK4D4/joincontext"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/hashicorp/go-argmapper"
	"github.com/hashicorp/vagrant-plugin-sdk/component"
	"github.com/hashicorp/vagrant-plugin-sdk/internal/funcspec"
	"github.com/hashicorp/vagrant-plugin-sdk/proto/vagrant_plugin_sdk"
	"google.golang.org/grpc"
)

type capabilityPlatform interface {
	HasCapability(ctx context.Context, in *vagrant_plugin_sdk.FuncSpec_Args, opts ...grpc.CallOption) (*vagrant_plugin_sdk.Platform_Capability_CheckResp, error)
	HasCapabilitySpec(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*vagrant_plugin_sdk.FuncSpec, error)
	Capability(ctx context.Context, in *vagrant_plugin_sdk.Platform_Capability_NamedRequest, opts ...grpc.CallOption) (*vagrant_plugin_sdk.Platform_Capability_Resp, error)
	CapabilitySpec(ctx context.Context, in *vagrant_plugin_sdk.Platform_Capability_NamedRequest, opts ...grpc.CallOption) (*vagrant_plugin_sdk.FuncSpec, error)
	Parents(ctx context.Context, in *vagrant_plugin_sdk.FuncSpec_Args, opts ...grpc.CallOption) (*vagrant_plugin_sdk.Platform_ParentsResp, error)
	ParentsSpec(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*vagrant_plugin_sdk.FuncSpec, error)
}

type capabilityClient struct {
	*baseClient
	client capabilityPlatform
}

func (c *capabilityClient) HasCapabilityFunc() interface{} {
	spec, err := c.client.HasCapabilitySpec(c.ctx, &empty.Empty{})
	if err != nil {
		return funcErr(err)
	}
	spec.Result = nil

	cb := func(ctx context.Context, args funcspec.Args) (bool, error) {
		resp, err := c.client.HasCapability(ctx, &vagrant_plugin_sdk.FuncSpec_Args{Args: args})

		if err != nil {
			return false, err
		}
		return resp.HasCapability, nil
	}
	return c.generateFunc(spec, cb)
}

func (c *capabilityClient) HasCapability(name string) (bool, error) {
	f := c.HasCapabilityFunc()
	n := &component.NamedCapability{Capability: name}
	raw, err := c.callDynamicFunc(f, (*bool)(nil),
		argmapper.Typed(n),
		argmapper.Typed(c.ctx),
	)
	if err != nil {
		return false, err
	}

	return raw.(bool), nil
}

func (c *capabilityClient) CapabilityFunc(name string) interface{} {
	spec, err := c.client.CapabilitySpec(c.ctx,
		&vagrant_plugin_sdk.Platform_Capability_NamedRequest{Name: name})
	if err != nil {
		return funcErr(err)
	}
	spec.Result = nil
	cb := func(ctx context.Context, args funcspec.Args) (interface{}, error) {
		resp, err := c.client.Capability(ctx,
			&vagrant_plugin_sdk.Platform_Capability_NamedRequest{
				FuncArgs: &vagrant_plugin_sdk.FuncSpec_Args{Args: args},
				Name:     name,
			},
		)

		if err != nil {
			return nil, err
		}
		return resp.Result, nil
	}
	return c.generateFunc(spec, cb)
}

func (c *capabilityClient) Capability(name string, args ...interface{}) (interface{}, error) {
	f := c.CapabilityFunc(name)
	raw, err := c.callDynamicFunc(f, false,
		argmapper.Typed(args...),
		argmapper.Typed(c.ctx),
	)
	if err != nil {
		return nil, err
	}

	return raw, nil
}

func (c *capabilityClient) ParentsFunc() interface{} {
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

func (c *capabilityClient) Parents() ([]string, error) {
	f := c.ParentsFunc()
	raw, err := c.callDynamicFunc(f, (*[]string)(nil),
		argmapper.Typed(c.ctx),
	)
	if err != nil {
		return nil, err
	}

	return raw.([]string), nil
}

type capabilityServer struct {
	*baseServer
	CapabilityImpl component.CapabilityPlatform
	typ            string
}

func (s *capabilityServer) HasCapabilitySpec(
	ctx context.Context,
	_ *empty.Empty,
) (*vagrant_plugin_sdk.FuncSpec, error) {
	if err := isImplemented(s, s.typ); err != nil {
		return nil, err
	}

	return s.generateSpec(s.CapabilityImpl.HasCapabilityFunc())
}

func (s *capabilityServer) HasCapability(
	ctx context.Context,
	args *vagrant_plugin_sdk.FuncSpec_Args,
) (*vagrant_plugin_sdk.Platform_Capability_CheckResp, error) {
	raw, err := s.callDynamicFunc(s.CapabilityImpl.HasCapabilityFunc(), (*bool)(nil),
		args.Args, argmapper.Typed(ctx))

	if err != nil {
		return nil, err
	}

	return &vagrant_plugin_sdk.Platform_Capability_CheckResp{
		HasCapability: raw.(bool)}, nil
}

func (s *capabilityServer) CapabilitySpec(
	ctx context.Context,
	req *vagrant_plugin_sdk.Platform_Capability_NamedRequest,
) (*vagrant_plugin_sdk.FuncSpec, error) {
	if err := isImplemented(s, s.typ); err != nil {
		return nil, err
	}

	return s.generateSpec(s.CapabilityImpl.CapabilityFunc(req.Name))
}

func (s *capabilityServer) Capability(
	ctx context.Context,
	args *vagrant_plugin_sdk.Platform_Capability_NamedRequest,
) (*vagrant_plugin_sdk.Platform_Capability_Resp, error) {
	_, err := s.callDynamicFunc(s.CapabilityImpl.CapabilityFunc(args.Name), false,
		args.FuncArgs.Args, argmapper.Typed(ctx))

	if err != nil {
		return nil, err
	}

	return &vagrant_plugin_sdk.Platform_Capability_Resp{}, nil
}

func (s *capabilityServer) ParentsSpec(
	ctx context.Context,
	_ *empty.Empty,
) (*vagrant_plugin_sdk.FuncSpec, error) {
	if err := isImplemented(s, s.typ); err != nil {
		return nil, err
	}

	return s.generateSpec(s.CapabilityImpl.ParentsFunc())
}

func (s *capabilityServer) Parents(
	ctx context.Context,
	args *vagrant_plugin_sdk.FuncSpec_Args,
) (*vagrant_plugin_sdk.Platform_ParentsResp, error) {
	raw, err := s.callDynamicFunc(s.CapabilityImpl.ParentsFunc(), (*[]string)(nil),
		args.Args, argmapper.Typed(ctx))

	if err != nil {
		return nil, err
	}

	return &vagrant_plugin_sdk.Platform_ParentsResp{
		Parents: raw.([]string)}, nil
}

var (
	_ component.CapabilityPlatform = (*capabilityClient)(nil)
)
