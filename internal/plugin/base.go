package plugin

import (
	"context"
	"net"

	"github.com/LK4D4/joincontext"
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

// BasePlugin contains the information which is common among
// all plugins. It should be embedded in every plugin type.
type BasePlugin struct {
	Cache   cacher.Cache      // Cache for mappers
	Mappers []*argmapper.Func // Mappers
	Logger  hclog.Logger      // Logger
	Wrapped bool              // Used to determine if wrapper
}

func (b *BasePlugin) Clone() *BasePlugin {
	return &BasePlugin{
		Cache:   b.Cache,
		Mappers: b.Mappers,
		Logger:  b.Logger,
		Wrapped: b.Wrapped,
	}
}

func (b *BasePlugin) NewClient(
	ctx context.Context,
	broker *plugin.GRPCBroker,
) *BaseClient {
	return &BaseClient{
		Ctx: ctx,
		Base: &Base{
			Broker:  broker,
			Cache:   b.Cache,
			Cleanup: &pluginargs.Cleanup{},
			Logger:  b.Logger,
			Mappers: b.Mappers,
			Wrapped: b.Wrapped,
		},
	}
}

func (b *BasePlugin) NewServer(
	broker *plugin.GRPCBroker,
) *BaseServer {
	return &BaseServer{
		Base: &Base{
			Broker:  broker,
			Cache:   b.Cache,
			Cleanup: &pluginargs.Cleanup{},
			Logger:  b.Logger,
			Mappers: b.Mappers,
			Wrapped: b.Wrapped,
		},
	}
}

// Base contains shared logic for all plugin server/client implementations.
// This should be embedded in every plugin server/client implementation using
// the specialized server and client types.
type Base struct {
	Broker  *plugin.GRPCBroker
	Logger  hclog.Logger
	Mappers []*argmapper.Func
	Cleanup *pluginargs.Cleanup
	Cache   cacher.Cache
	Wrapped bool
}

func (b *Base) Wrap() *BasePlugin {
	return &BasePlugin{
		Logger:  b.Logger,
		Mappers: b.Mappers,
		Cache:   b.Cache,
		Wrapped: true,
	}
}

// If this plugin is a wrapper
func (b *Base) IsWrapped() bool {
	return b.Wrapped
}

// Base client type
type BaseClient struct {
	*Base

	Ctx          context.Context
	target       net.Addr
	parentPlugin interface{}
}

// Base server type
type BaseServer struct {
	*Base
}

// internal returns a new pluginargs.Internal that can be used with
// dynamic calls. The Internal structure is an internal-only argument
// that is used to perform cleanup.
func (b *Base) internal() *pluginargs.Internal {
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

// Map a value to the expected type using registered mappers
// NOTE: The expected type must be a pointer, so an expected type
// of `*int` means an `int` is wanted. Expected type of `**int`
// means an `*int` is wanted, etc.
func (b *Base) Map(
	resultValue, // value to be converted
	expectedType interface{}, // nil pointer of desired type
	args ...argmapper.Arg, // list of argmapper arguments
) (interface{}, error) {
	args = append(args,
		argmapper.ConverterFunc(MapperFns...),
		argmapper.ConverterFunc(b.Mappers...),
		argmapper.Typed(b.internal()),
		argmapper.Typed(b.Logger),
	)

	return dynamic.Map(resultValue, expectedType, args...)
}

func (b *Base) SetCache(c cacher.Cache) {
	b.Cache = c
}

// Sets the parent plugins
func (b *BaseClient) SetParentPlugin(plugin interface{}) {
	b.parentPlugin = plugin
}

func (b *BaseClient) AppendMappers(mappers ...*argmapper.Func) {
	b.Mappers = append(b.Mappers, mappers...)
}

func (b *BaseClient) GenerateContext(ctx context.Context) (context.Context, context.CancelFunc) {
	return joincontext.Join(ctx, b.Ctx)
}

func (b *BaseClient) Close() error {
	return b.Cleanup.Close()
}

// Used internally to extract broker
func (b *BaseClient) GRPCBroker() *plugin.GRPCBroker {
	return b.Broker
}

// Sets a direct target which can be connected
// to when passing this client over proto.
func (b *BaseClient) SetTarget(t net.Addr) {
	b.target = t
}

// Provides the direct target being used
// by this client.
func (b *BaseClient) Target() net.Addr {
	return b.target
}

// This is here for internal usage on plugin setup
// to provide extra information to ruby Based plugins
func (b *BaseClient) SetRequestMetadata(key, value string) {
	b.Ctx = metadata.AppendToOutgoingContext(b.Ctx, key, value)
	b.Logger.Trace("new metadata has been set for outgoing requests",
		"key", key, "value", value)
}

// Generate a function from a provided spec
func (b *BaseClient) GenerateFunc(
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
func (b *BaseClient) CallDynamicFunc(
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
func (b *BaseServer) CallDynamicFunc(
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

// Generate a funcspec Based on the provided function
func (b *BaseServer) GenerateSpec(
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
