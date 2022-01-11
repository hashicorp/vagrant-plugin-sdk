package plugin

import (
	"context"

	"github.com/LK4D4/joincontext"
	"github.com/hashicorp/go-argmapper"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
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
}

type capabilityComponent interface {
	HasCapabilityFunc() interface{}
	HasCapability(name string) (bool, error)
	CapabilityFunc(name string) interface{}
	Capability(name string, args ...interface{}) (interface{}, error)
}

type capabilityClient struct {
	*BaseClient
	client capabilityPlatform
}

func (c *capabilityClient) HasCapabilityFunc() interface{} {
	spec, err := c.client.HasCapabilitySpec(c.Ctx, &emptypb.Empty{})
	if err != nil {
		return funcErr(err)
	}
	spec.Result = nil

	cb := func(ctx context.Context, args funcspec.Args) (bool, error) {
		new_ctx, _ := joincontext.Join(c.Ctx, ctx)
		resp, err := c.client.HasCapability(new_ctx, &vagrant_plugin_sdk.FuncSpec_Args{Args: args})

		if err != nil {
			return false, err
		}

		if !resp.HasCapability && c.parentComponent != nil {
			// Check the parent plugin for the capability
			parentComponent := c.parentComponent.(capabilityComponent)
			f := parentComponent.HasCapabilityFunc()
			parentRequestArgs := []argmapper.Arg{argmapper.Typed(new_ctx)}
			for _, a := range args {
				parentRequestArgs = append(parentRequestArgs, argmapper.Typed(a.Value))
			}
			raw, err := dynamic.CallFunc(f, (*bool)(nil), c.Mappers, parentRequestArgs...)
			return raw.(bool), err
		}

		return resp.HasCapability, nil
	}
	return c.GenerateFunc(spec, cb)
}

func (c *capabilityClient) HasCapability(name string) (bool, error) {
	f := c.HasCapabilityFunc()
	n := &component.NamedCapability{Capability: name}
	raw, err := c.CallDynamicFunc(f, (*bool)(nil),
		argmapper.Typed(n),
		argmapper.Typed(c.Ctx),
	)
	if err != nil {
		return false, err
	}

	c.Logger.Info("check for capability complete",
		"capability_name", name,
		"result", raw,
	)
	return raw.(bool), nil
}

func (c *capabilityClient) CapabilityFunc(name string) interface{} {
	if c.parentComponent != nil {
		ok, err := c.parentComponent.(capabilityComponent).HasCapability(name)

		if err != nil {
			c.Logger.Error("parent capability check failed",
				"capability_name", name,
			)
		}

		if ok {
			c.Logger.Trace("capability detected on parent plugin",
				"capability_name", name,
			)

			return c.parentComponent.(capabilityComponent).CapabilityFunc(name)
		}
	}

	spec, err := c.client.CapabilitySpec(c.Ctx,
		&vagrant_plugin_sdk.Platform_Capability_NamedRequest{Name: name})
	if err != nil {
		return funcErr(err)
	}
	spec.Result = nil
	cb := func(ctx context.Context, args funcspec.Args) (interface{}, error) {
		ctx, _ = joincontext.Join(c.Ctx, ctx)
		resp, err := c.client.Capability(ctx,
			&vagrant_plugin_sdk.Platform_Capability_NamedRequest{
				FuncArgs: &vagrant_plugin_sdk.FuncSpec_Args{Args: args},
				Name:     name,
			},
		)

		if err != nil {
			c.Logger.Error("failure encountered while running capability",
				"name", name,
			)

			return nil, err
		}

		if resp.Result == nil {
			return nil, nil
		}

		// Result will be returned as an Any so decode it
		_, val, err := dynamic.DecodeAny(resp.Result)
		if err != nil {
			c.Logger.Error("failure while attempting to decode capability result",
				"name", name,
				"result", resp.Result,
				"error", err,
			)

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
				argmapper.Typed(c.Internal()),
				argmapper.Typed(c.Ctx),
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
	return c.GenerateFunc(spec, cb)
}

func (c *capabilityClient) Capability(name string, args ...interface{}) (interface{}, error) {
	f := c.CapabilityFunc(name)

	return c.CallDynamicFunc(f, false,
		argmapper.Typed(&component.Direct{Arguments: args}),
		argmapper.Typed(args...),
		argmapper.Typed(c.Ctx),
	)
}

type capabilityServer struct {
	*BaseServer
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

	return s.GenerateSpec(s.CapabilityImpl.HasCapabilityFunc())
}

func (s *capabilityServer) HasCapability(
	ctx context.Context,
	args *vagrant_plugin_sdk.FuncSpec_Args,
) (*vagrant_plugin_sdk.Platform_Capability_CheckResp, error) {
	raw, err := s.CallDynamicFunc(s.CapabilityImpl.HasCapabilityFunc(), (*bool)(nil),
		args.Args, argmapper.Typed(ctx))

	if err != nil {
		s.Logger.Error("capability check failed",
			"error", err,
		)

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

	return s.GenerateSpec(s.CapabilityImpl.CapabilityFunc(req.Name))
}

func (s *capabilityServer) Capability(
	ctx context.Context,
	args *vagrant_plugin_sdk.Platform_Capability_NamedRequest,
) (*vagrant_plugin_sdk.Platform_Capability_Resp, error) {
	v, err := s.CallDynamicFunc(
		s.CapabilityImpl.CapabilityFunc(args.Name),
		false,
		args.FuncArgs.Args,
		argmapper.Typed(ctx),
	)

	if err != nil {
		s.Logger.Error("failed to call capability",
			"name", args.Name)

		return nil, err
	}

	if v == nil {
		return &vagrant_plugin_sdk.Platform_Capability_Resp{}, nil
	}

	var val interface{}
	val, err = dynamic.MapToWellKnownProto(v)
	if err != nil {
		val, err = dynamic.UnknownMap(v, (*proto.Message)(nil), s.Mappers,
			argmapper.Typed(s.Internal()),
			argmapper.Typed(ctx),
			argmapper.Typed(s.Logger),
		)

		if err != nil {
			s.Logger.Error("failed to convert result value",
				"error", err)

			return nil, err
		}
	}

	var result *anypb.Any
	if val == nil {
		result = nil
	} else {
		result, err = dynamic.EncodeAny(val.(proto.Message))
		if err != nil {
			s.Logger.Error("failed to encode capability response message",
				"value", val,
				"error", err,
			)

			return nil, err
		}
	}

	return &vagrant_plugin_sdk.Platform_Capability_Resp{
		Result: result,
	}, nil
}

var (
	_ component.CapabilityPlatform = (*capabilityClient)(nil)
)
