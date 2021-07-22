package plugin

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/hashicorp/go-argmapper"
	"github.com/hashicorp/vagrant-plugin-sdk/component"
	"github.com/hashicorp/vagrant-plugin-sdk/internal/funcspec"
	"github.com/hashicorp/vagrant-plugin-sdk/proto/vagrant_plugin_sdk"
	"google.golang.org/grpc"
)

type capabilityHost interface {
	HasCapability(ctx context.Context, in *vagrant_plugin_sdk.FuncSpec_Args, opts ...grpc.CallOption) (*vagrant_plugin_sdk.Host_Capability_CheckResp, error)
	HasCapabilitySpec(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*vagrant_plugin_sdk.FuncSpec, error)
	Capability(ctx context.Context, in *vagrant_plugin_sdk.Host_Capability_NamedRequest, opts ...grpc.CallOption) (*vagrant_plugin_sdk.Host_Capability_Resp, error)
	CapabilitySpec(ctx context.Context, in *vagrant_plugin_sdk.Host_Capability_NamedRequest, opts ...grpc.CallOption) (*vagrant_plugin_sdk.FuncSpec, error)
}

type capabilityClient struct {
	*baseClient
	client capabilityHost
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
		&vagrant_plugin_sdk.Host_Capability_NamedRequest{Name: name})
	if err != nil {
		return funcErr(err)
	}
	spec.Result = nil
	cb := func(ctx context.Context, args funcspec.Args) (interface{}, error) {
		resp, err := c.client.Capability(ctx,
			&vagrant_plugin_sdk.Host_Capability_NamedRequest{
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
