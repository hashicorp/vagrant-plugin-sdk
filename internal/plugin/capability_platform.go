package plugin

import (
	"context"

	"github.com/LK4D4/joincontext"
	"github.com/hashicorp/go-argmapper"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/hashicorp/vagrant-plugin-sdk/component"
	"github.com/hashicorp/vagrant-plugin-sdk/internal-shared/dynamic"
	"github.com/hashicorp/vagrant-plugin-sdk/internal/funcspec"
	"github.com/hashicorp/vagrant-plugin-sdk/proto/vagrant_plugin_sdk"
)

type capabilityPlatform interface {
	HasCapability(ctx context.Context, in *vagrant_plugin_sdk.FuncSpec_Args, opts ...grpc.CallOption) (*vagrant_plugin_sdk.Platform_Capability_CheckResp, error)
	HasCapabilitySpec(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*vagrant_plugin_sdk.FuncSpec, error)
	Capability(ctx context.Context, in *vagrant_plugin_sdk.Platform_Capability_NamedRequest, opts ...grpc.CallOption) (*vagrant_plugin_sdk.Platform_Capability_Resp, error)
	CapabilitySpec(ctx context.Context, in *vagrant_plugin_sdk.Platform_Capability_NamedRequest, opts ...grpc.CallOption) (*vagrant_plugin_sdk.FuncSpec, error)
	Parents(ctx context.Context, in *vagrant_plugin_sdk.FuncSpec_Args, opts ...grpc.CallOption) (*vagrant_plugin_sdk.Platform_ParentsResp, error)
	ParentsSpec(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*vagrant_plugin_sdk.FuncSpec, error)
}

type capabilityClient struct {
	*baseClient
	client capabilityPlatform
}

type CapabilityArguments struct {
	Arguments []interface{}
}

func (c *capabilityClient) HasCapabilityFunc() interface{} {
	spec, err := c.client.HasCapabilitySpec(c.ctx, &emptypb.Empty{})
	if err != nil {
		return funcErr(err)
	}
	spec.Result = nil

	cb := func(ctx context.Context, args funcspec.Args) (bool, error) {
		ctx, _ = joincontext.Join(c.ctx, ctx)
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
		ctx, _ = joincontext.Join(c.ctx, ctx)
		resp, err := c.client.Capability(ctx,
			&vagrant_plugin_sdk.Platform_Capability_NamedRequest{
				FuncArgs: &vagrant_plugin_sdk.FuncSpec_Args{Args: args},
				Name:     name,
			},
		)

		if err != nil {
			return nil, err
		}

		// Result will be returned as an Any so decode it
		_, val, err := dynamic.DecodeAny(resp.Result)
		if err != nil {
			return nil, err
		}

		// Start by attempting to map the decoded value using
		// the well known type protos
		result, err := dynamic.MapFromWellKnownProto(val.(proto.Message))
		if err != nil {
			// And the actual result can be pretty much anything, so
			// attempt to map it into something usable. If we can't
			// map it, log the mapping failure but then just return
			// the decoded value
			result, err = dynamic.BlindMap(val, c.Mappers,
				argmapper.Typed(c.internal()),
				argmapper.Typed(c.ctx),
				argmapper.Typed(c.Logger),
			)

			if err != nil {
				c.Logger.Debug("failed to map decoded result from capability",
					"value", val,
					"error", err,
				)

				return val, err
			}
		}

		return result, nil
	}
	return c.generateFunc(spec, cb)
}

func (c *capabilityClient) Capability(name string, args ...interface{}) (interface{}, error) {
	f := c.CapabilityFunc(name)

	return c.callDynamicFunc(f, false,
		argmapper.Typed(args...),
		argmapper.Typed(&CapabilityArguments{Arguments: args}),
		argmapper.Typed(c.ctx),
	)
}

type capabilityServer struct {
	*baseServer
	CapabilityImpl component.CapabilityPlatform
	typ            string
}

func (s *capabilityServer) HasCapabilitySpec(
	ctx context.Context,
	_ *emptypb.Empty,
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
	v, err := s.callDynamicFunc(s.CapabilityImpl.CapabilityFunc(args.Name), false,
		args.FuncArgs.Args, argmapper.Typed(ctx))

	if err != nil {
		return nil, err
	}

	var val interface{}
	val, err = dynamic.MapToWellKnownProto(v)
	if err != nil {
		val, err = dynamic.UnknownMap(v, (*proto.Message)(nil), s.Mappers,
			argmapper.Typed(s.internal()),
			argmapper.Typed(ctx),
			argmapper.Typed(s.Logger),
		)

		if err != nil {
			return nil, err
		}
	}

	result, err := dynamic.EncodeAny(val.(proto.Message))
	if err != nil {
		return nil, err
	}

	return &vagrant_plugin_sdk.Platform_Capability_Resp{
		Result: result,
	}, nil
}

var (
	_ component.CapabilityPlatform = (*capabilityClient)(nil)
)
