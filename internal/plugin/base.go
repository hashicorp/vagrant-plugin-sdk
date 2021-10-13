package plugin

import (
	"context"
	"net"

	"github.com/hashicorp/go-argmapper"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/hashicorp/vagrant-plugin-sdk/internal-shared/cacher"
	"github.com/hashicorp/vagrant-plugin-sdk/internal-shared/dynamic"
	"github.com/hashicorp/vagrant-plugin-sdk/internal/funcspec"
	"github.com/hashicorp/vagrant-plugin-sdk/internal/pluginargs"
	"github.com/hashicorp/vagrant-plugin-sdk/proto/vagrant_plugin_sdk"
)

func isImplemented(t interface{}, name string) error {
	if t == nil {
		return status.Errorf(codes.Unimplemented, "plugin does not implement: "+name)
	}
	return nil
}

// base contains shared logic for all plugins. This should be embedded
// in every plugin implementation.
type base struct {
	Broker  *plugin.GRPCBroker
	Logger  hclog.Logger
	Mappers []*argmapper.Func
	Cleanup *pluginargs.Cleanup
	Cache   cacher.Cache
	Wrapped bool
}

func (b *base) IsWrapped() bool {
	return b.Wrapped
}

type baseClient struct {
	*base

	ctx    context.Context
	target net.Addr
}

type baseServer struct {
	*base
}

// internal returns a new pluginargs.Internal that can be used with
// dynamic calls. The Internal structure is an internal-only argument
// that is used to perform cleanup.
func (b *base) internal() *pluginargs.Internal {
	// if the cache isn't currently set, just create
	// a new cache instance and set it now
	if b.Cache == nil {
		b.Cache = cacher.New()
	}

	return &pluginargs.Internal{
		Broker:  b.Broker,
		Mappers: b.Mappers,
		Cleanup: b.Cleanup,
		Cache:   b.Cache,
		Logger:  b.Logger,
	}
}

func (b *baseClient) Close() error {
	return b.Cleanup.Close()
}

func (b *base) SetCache(c cacher.Cache) {
	b.Cache = c
}

// Used internally to extract broker
func (b *baseClient) GRPCBroker() *plugin.GRPCBroker {
	return b.Broker
}

// Sets a direct target which can be connected
// to when passing this client over proto.
func (b *baseClient) SetTarget(t net.Addr) {
	b.target = t
}

// Provides the direct target being used
// by this client.
func (b *baseClient) Target() net.Addr {
	return b.target
}

// This is here for internal usage on plugin setup
// to provide extra information to ruby based plugins
func (b *baseClient) SetRequestMetadata(key, value string) {
	b.ctx = metadata.AppendToOutgoingContext(b.ctx, key, value)
	b.Logger.Trace("new metadata has been set for outgoing requests",
		"key", key, "value", value)
}

// Generate a function from a provided spec
func (b *baseClient) generateFunc(
	spec *vagrant_plugin_sdk.FuncSpec, // spec for the function
	cbFn interface{}, // callback function
	args ...argmapper.Arg, // any extra argmapper args
) interface{} {
	return funcspec.Func(spec, cbFn, append(args,
		argmapper.Typed(b.internal()))...,
	)
}

// Calls the function provided and converts the
// result to an expected type. If no type conversion
// is required, a `false` value for the expectedType
// will return the raw interface return value. Automatically
// provided args include hclog.Logger and pluginargs.Internal
// typed arguments, registered mappers, and a custom logger
// for argmapper.
//
// NOTE: Provide a `false` value for expectedType if no
// type conversion is required.
func (b *baseClient) callDynamicFunc(
	f interface{}, // function to call
	expectedType interface{}, // nil pointer of expected return type
	callArgs ...argmapper.Arg, // any extra argmapper arguments to include
) (interface{}, error) {
	internal := b.internal()
	// TODO(spox): We need to determine how to properly cleanup when connections
	//             may still exist after the dynamic call is complete
	//	defer internal.Cleanup.Close()
	callArgs = append(callArgs,
		argmapper.Typed(internal),
		argmapper.Typed(b.Logger),
	)

	return dynamic.CallFunc(f, expectedType, b.Mappers, callArgs...)
}

// Calls the function provided and converts the
// result to an expected type. If no type conversion
// is required, a `false` value for the expectedType
// will return the raw interface return value. Automatically
// provided args include hclog.Logger and pluginargs.Internal
// typed arguments, registered mappers, and a custom logger
// for argmapper.
//
// NOTE: Provide a `false` value for expectedType if no
// type conversion is required.
func (b *baseServer) callDynamicFunc(
	f interface{}, // function to call
	expectedType interface{}, // nil pointer of expected return type
	args funcspec.Args, // funspec defined arguments
	callArgs ...argmapper.Arg, // any extra argmapper arguments to include
) (interface{}, error) {
	internal := b.internal()

	// Decode our *any.Any values.
	for _, arg := range args {
		anyVal := arg.Value

		_, v, err := dynamic.DecodeAny(anyVal)
		if err != nil {
			return nil, err
		}

		callArgs = append(callArgs,
			argmapper.NamedSubtype(arg.Name, v, arg.Type),
		)
	}
	callArgs = append(callArgs,
		argmapper.Typed(internal),
		argmapper.Typed(b.Logger),
	)

	return dynamic.CallFunc(f, expectedType, b.Mappers, callArgs...)
}

// Generate a funcspec based on the provided function
func (b *baseServer) generateSpec(
	fn interface{}, // function to generate funcspec
	args ...argmapper.Arg, // optional argmapper args
) (*vagrant_plugin_sdk.FuncSpec, error) {
	if f, ok := fn.(*dynamic.SpecAndFunc); ok {
		return f.Spec, nil
	}
	f, err := funcspec.Spec(fn, append(args,
		argmapper.ConverterFunc(b.Mappers...),
		argmapper.Typed(b.internal()))...,
	)
	if err != nil {
		return f, err
	}
	return f, err
}
