package plugin

import (
	"context"
	"errors"
	"fmt"

	"github.com/LK4D4/joincontext"
	"github.com/hashicorp/go-argmapper"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/hashicorp/vagrant-plugin-sdk/component"
	"github.com/hashicorp/vagrant-plugin-sdk/core"
	"github.com/hashicorp/vagrant-plugin-sdk/internal-shared/dynamic"
	"github.com/hashicorp/vagrant-plugin-sdk/internal/funcspec"
	"github.com/hashicorp/vagrant-plugin-sdk/proto/vagrant_plugin_sdk"
)

type capabilityParent interface {
	GenerateContext(ctx context.Context) (context.Context, context.CancelFunc)
	PluginHasCapability() interface{}
	HasCapability(name string) (bool, error)
	GetCapabilityClient() *capabilityClient
}

type capabilityPlatform interface {
	HasCapability(ctx context.Context, in *vagrant_plugin_sdk.FuncSpec_Args, opts ...grpc.CallOption) (*vagrant_plugin_sdk.Platform_Capability_CheckResp, error)
	HasCapabilitySpec(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*vagrant_plugin_sdk.FuncSpec, error)
	Capability(ctx context.Context, in *vagrant_plugin_sdk.Platform_Capability_NamedRequest, opts ...grpc.CallOption) (*vagrant_plugin_sdk.Platform_Capability_Resp, error)
	CapabilitySpec(ctx context.Context, in *vagrant_plugin_sdk.Platform_Capability_NamedRequest, opts ...grpc.CallOption) (*vagrant_plugin_sdk.FuncSpec, error)
	Seeds(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*vagrant_plugin_sdk.Args_Direct, error)
	Seed(ctx context.Context, in *vagrant_plugin_sdk.Args_Direct, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

type capabilityClient struct {
	*BaseClient
	client capabilityPlatform
}

func (c *capabilityClient) GetCapabilityClient() *capabilityClient {
	return c
}

func (c *capabilityClient) Seed(args ...interface{}) error {
	cb := func(d *vagrant_plugin_sdk.Args_Direct) error {
		_, err := c.client.Seed(c.Ctx, d)
		return err
	}

	_, err := c.CallDynamicFunc(cb, false,
		argmapper.Typed(c.Ctx),
		argmapper.Typed(&component.Direct{Arguments: args}),
	)

	return err
}

func (c *capabilityClient) Seeds() ([]interface{}, error) {
	cb := func() (*vagrant_plugin_sdk.Args_Direct, error) {
		return c.client.Seeds(c.Ctx, &emptypb.Empty{})
	}

	r, err := c.CallDynamicFunc(cb,
		(**component.Direct)(nil),
		argmapper.Typed(c.Ctx),
	)

	if err != nil {
		return nil, err
	}

	return r.(*component.Direct).Arguments, nil
}

func (c *capabilityClient) getCapabilityFromParent(ctx context.Context, args funcspec.Args) (interface{}, error) {
	for _, p := range c.parentPlugins {
		parentPlugin := p.(capabilityParent)
		new_ctx, _ := parentPlugin.GenerateContext(ctx)
		f := parentPlugin.PluginHasCapability()
		parentRequestArgs := []argmapper.Arg{argmapper.Typed(new_ctx)}
		for _, a := range args {
			parentRequestArgs = append(parentRequestArgs, argmapper.Typed(a.Value))
		}
		raw, err := dynamic.CallFunc(f, (*bool)(nil), c.Mappers, parentRequestArgs...)
		if err != nil {
			return nil, err
		}
		if raw.(bool) {
			return p, nil
		}
	}
	return nil, errors.New("could not find capability in parent plugins")
}

func (c *capabilityClient) getCapabilityFromParent2(ctx context.Context, name string) (interface{}, error) {
	for _, p := range c.parentPlugins {
		parentPlugin := p.(capabilityParent)
		// new_ctx, _ := parentPlugin.GenerateContext(ctx)
		hasCap, err := parentPlugin.HasCapability(name)
		if err != nil {
			return nil, err
		}
		if hasCap {
			return p, nil
		}
	}
	return nil, errors.New("could not find capability in parent plugins")
}

func (c *capabilityClient) PluginHasCapability() interface{} {
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
		return resp.HasCapability, nil
	}
	return c.GenerateFunc(spec, cb)
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

		if !resp.HasCapability {
			p, _ := c.getCapabilityFromParent(ctx, args)
			if p != nil {
				return true, nil
			}
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

	return raw.(bool), nil
}

func (c *capabilityClient) CapabilityFunc(name string) interface{} {
	p, _ := c.getCapabilityFromParent2(c.Ctx, name)
	var pluginWithCapability interface{}
	if p == nil {
		pluginWithCapability = c
	} else {
		pluginWithCapability = p
	}

	client := pluginWithCapability.(capabilityParent).GetCapabilityClient()

	spec, err := client.client.CapabilitySpec(c.Ctx,
		&vagrant_plugin_sdk.Platform_Capability_NamedRequest{Name: name})
	if err != nil {
		return funcErr(err)
	}
	spec.Result = nil
	cb := func(ctx context.Context, args funcspec.Args) (interface{}, error) {
		ctx, _ = joincontext.Join(client.Ctx, ctx)
		resp, err := client.client.Capability(ctx,
			&vagrant_plugin_sdk.Platform_Capability_NamedRequest{
				FuncArgs: &vagrant_plugin_sdk.FuncSpec_Args{Args: args},
				Name:     name,
			},
		)

		if err != nil {
			c.Logger.Error("failure encountered while running capability",
				"name", name,
				"error", err,
			)

			return nil, err
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
				argmapper.Typed(c.internal()),
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

	ex, err := c.Seeds()
	if err != nil {
		c.Logger.Error("failed to fetch seed values for capability call",
			"name", name,
			"error", err,
		)

		return nil, err
	}

	return c.CallDynamicFunc(f, false,
		argmapper.Typed(&component.Direct{Arguments: args}),
		argmapper.Typed(args...),
		argmapper.Typed(ex...),
		argmapper.Typed(c.Ctx),
	)
}

type capabilityServer struct {
	*BaseServer
	CapabilityImpl component.CapabilityPlatform
	typ            string
	seeds          []*anypb.Any
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
			"error", err)

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
			s.Logger.Error("failed to convert result value",
				"error", err)

			return nil, err
		}
	}

	result, err := dynamic.EncodeAny(val.(proto.Message))
	if err != nil {
		s.Logger.Error("failed to encode capability response message",
			"value", val,
			"error", err,
		)

		return nil, err
	}

	return &vagrant_plugin_sdk.Platform_Capability_Resp{
		Result: result,
	}, nil
}

func (s *capabilityServer) Seed(
	ctx context.Context,
	args *vagrant_plugin_sdk.Args_Direct,
) (*emptypb.Empty, error) {
	if !s.IsWrapped() {
		s.seeds = args.List
		return &emptypb.Empty{}, nil
	}

	v, err := dynamic.Map(args, (**component.Direct)(nil),
		argmapper.Typed(ctx, s.internal(), s.Logger),
		argmapper.ConverterFunc(s.Mappers...),
	)

	if err != nil {
		s.Logger.Error("failed to store seed values",
			"error", err,
		)

		return nil, err
	}

	if seeder, ok := s.CapabilityImpl.(core.Seeder); ok {
		err = seeder.Seed(v.(*component.Direct).Arguments...)
		return &emptypb.Empty{}, err
	}

	return nil, fmt.Errorf("failed to properly seed plugin")
}

func (s *capabilityServer) Seeds(
	ctx context.Context,
	_ *emptypb.Empty,
) (*vagrant_plugin_sdk.Args_Direct, error) {
	if !s.IsWrapped() {
		return &vagrant_plugin_sdk.Args_Direct{
			List: s.seeds,
		}, nil
	}

	var seeder core.Seeder
	var ok bool

	if seeder, ok = s.CapabilityImpl.(core.Seeder); !ok {
		s.Logger.Error("plugin implementation does not support seeding")
		return nil, fmt.Errorf("plugin is not a valid seeder")
	}

	vals, err := seeder.Seeds()
	if err != nil {
		s.Logger.Error("failed to fetch seed values",
			"error", err,
		)

		return nil, err
	}

	r, err := dynamic.Map(
		&component.Direct{Arguments: vals},
		(**vagrant_plugin_sdk.Args_Direct)(nil),
		argmapper.Typed(ctx, s.internal(), s.Logger),
		argmapper.ConverterFunc(s.Mappers...),
	)

	if err != nil {
		s.Logger.Error("failed to convert seed values into proto message",
			"values", vals,
			"error", err,
		)

		return nil, err
	}

	return r.(*vagrant_plugin_sdk.Args_Direct), nil
}

var (
	_ component.CapabilityPlatform = (*capabilityClient)(nil)
)
