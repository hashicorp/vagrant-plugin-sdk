package plugin

import (
	"context"
	"fmt"
	"net"

	"github.com/LK4D4/joincontext"
	"github.com/hashicorp/go-argmapper"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/hashicorp/vagrant-plugin-sdk/core"
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

type SeederClient interface {
	Seed(ctx context.Context, in *vagrant_plugin_sdk.Args_Seeds, opts ...grpc.CallOption) (*emptypb.Empty, error)
	Seeds(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*vagrant_plugin_sdk.Args_Seeds, error)
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
	s SeederClient,
) *BaseClient {
	return &BaseClient{
		Ctx:    ctx,
		Seeder: s,
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
	impl interface{},
) *BaseServer {
	return &BaseServer{
		impl:       impl,
		seedValues: &vagrant_plugin_sdk.Args_Seeds{},
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

	Ctx             context.Context
	Seeder          SeederClient
	addr            net.Addr
	parentComponent interface{}
}

// Base server type
type BaseServer struct {
	*Base

	impl       interface{}
	seedValues *vagrant_plugin_sdk.Args_Seeds
}

// internal returns a new pluginargs.Internal that can be used with
// dynamic calls. The Internal structure is an internal-only argument
// that is used to perform cleanup.
func (b *Base) Internal() *pluginargs.Internal {
	// if the cache isn't currently set, just create
	// a new cache instance and set it now
	if b.Cache == nil {
		b.Cache = cacher.New()
	}

	m := make([]*argmapper.Func, len(b.Mappers)+len(MapperFns))
	copy(m, b.Mappers)
	copy(m[len(b.Mappers):], MapperFns)
	return &pluginargs.Internal{
		Broker:  b.Broker,
		Mappers: m,
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
		argmapper.Typed(b.Internal()),
		argmapper.Typed(b.Logger),
	)

	return dynamic.Map(resultValue, expectedType, args...)
}

func (b *Base) SetCache(c cacher.Cache) {
	b.Cache = c
}

func (b *BaseClient) Seed(args *core.Seeds) error {
	if b.Seeder == nil {
		b.Logger.Trace("plugin does not implement seeder interface")
		return nil
	}

	cb := func(d *vagrant_plugin_sdk.Args_Seeds) error {
		_, err := b.Seeder.Seed(b.Ctx, d)
		return err
	}

	_, err := b.CallDynamicFunc(cb, false,
		argmapper.Typed(b.Ctx),
		argmapper.Typed(args),
	)

	return err
}

func (b *BaseClient) Seeds() (*core.Seeds, error) {
	if b.Seeder == nil {
		b.Logger.Trace("plugin does not implement seeder interface")
		return core.NewSeeds(), nil
	}

	r, err := b.Seeder.Seeds(b.Ctx, &emptypb.Empty{})
	if err != nil {
		b.Logger.Error("failed to get seed values",
			"error", err,
		)

		return nil, err
	}

	s, err := b.Map(r, (**core.Seeds)(nil), argmapper.Typed(b.Ctx))
	if err != nil {
		b.Logger.Error("failed to convert seeds value response to proper type",
			"value", r,
			"error", err,
		)

		return nil, err
	}

	return s.(*core.Seeds), nil
}

// Sets the parent component
func (b *BaseClient) SetParentComponent(c interface{}) {
	b.parentComponent = c
}

func (b *BaseClient) GetParentComponent() interface{} {
	return b.parentComponent
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

// Sets a direct addr which can be connected
// to when passing this client over proto.
func (b *BaseClient) SetAddr(t net.Addr) {
	b.addr = t
}

// Provides the direct addr being used
// by this client.
func (b *BaseClient) Addr() net.Addr {
	return b.addr
}

// This is here for internal usage on plugin setup
// to provide extra information to ruby Based plugins
func (b *BaseClient) SetRequestMetadata(key, value string) {
	md, ok := metadata.FromOutgoingContext(b.Ctx)
	if !ok {
		md = metadata.New(map[string]string{})
	}
	md[key] = []string{value}
	b.Ctx = metadata.NewOutgoingContext(b.Ctx, md)
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
		argmapper.Typed(b.Internal()))...,
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
	internal := b.Internal()
	// TODO(spox): We need to determine how to properly cleanup when connections
	//             may still exist after the dynamic call is complete
	//	defer internal.Cleanup.Close()

	if b.Seeder != nil {
		s, err := b.Seeds()
		if err != nil {
			b.Logger.Error("failed to fetch dynamic seed values",
				"error", err,
			)

			return nil, err
		}

		for _, v := range s.Typed {
			if a, ok := v.(*anypb.Any); ok {
				b.Logger.Info("seeding typed value into dynamic call",
					"type", hclog.Fmt("%T", v),
					"subtype", a.TypeUrl,
				)

				callArgs = append(callArgs, argmapper.TypedSubtype(a, a.TypeUrl))
			} else {
				b.Logger.Info("seeding typed value into dynamic call",
					"type", hclog.Fmt("%T", v),
				)

				callArgs = append(callArgs, argmapper.Typed(v))
			}
		}

		for k := range s.Named {
			v := s.Named[k]
			if a, ok := v.(*anypb.Any); ok {
				b.Logger.Info("seeding named value into dynamic call",
					"name", k,
					"type", hclog.Fmt("%T", v),
					"subtype", a.TypeUrl,
				)

				callArgs = append(callArgs, argmapper.NamedSubtype(k, a, a.TypeUrl))
			} else {
				b.Logger.Info("seeding named value into dynamic call",
					"name", k,
					"type", hclog.Fmt("%T", v),
				)

				callArgs = append(callArgs, argmapper.Named(k, v))
			}
		}
	}

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
	internal := b.Internal()

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
		argmapper.Typed(b.Internal()))...,
	)
	if err != nil {
		return f, err
	}
	return f, err
}

func (b *BaseServer) Seed(
	ctx context.Context,
	seeds *vagrant_plugin_sdk.Args_Seeds,
) (*emptypb.Empty, error) {
	if b.impl == nil {
		b.Logger.Trace("plugin does not implement seeder interface")
		return &emptypb.Empty{}, nil
	}

	if !b.IsWrapped() {
		b.seedValues = seeds
		return &emptypb.Empty{}, nil
	}

	seeder, ok := b.impl.(core.Seeder)
	if !ok {
		b.Logger.Error("plugin implementation does not provide core.Seeder",
			"impl", b.impl,
		)

		return nil, fmt.Errorf("implementation does not support value seeds")
	}

	v, err := dynamic.Map(seeds, (**core.Seeds)(nil),
		argmapper.Typed(ctx, b.Internal(), b.Logger),
		argmapper.ConverterFunc(b.Mappers...),
	)

	if err != nil {
		b.Logger.Error("failed to store seed values",
			"error", err,
		)

		return nil, err
	}

	err = seeder.Seed(v.(*core.Seeds))

	if err != nil {
		b.Logger.Error("failed to store seed values",
			"error", err,
		)
	}
	return &emptypb.Empty{}, err
}

func (b *BaseServer) Seeds(
	ctx context.Context,
	_ *emptypb.Empty,
) (*vagrant_plugin_sdk.Args_Seeds, error) {
	if b.impl == nil {
		b.Logger.Trace("plugin does not implement seeder interface")
		return &vagrant_plugin_sdk.Args_Seeds{}, nil
	}

	if !b.IsWrapped() {
		return b.seedValues, nil
	}

	seeder, ok := b.impl.(core.Seeder)
	if !ok {
		b.Logger.Error("plugin implementation does not provide core.Seeder",
			"impl", b.impl,
		)

		return nil, fmt.Errorf("implementation does not support value seeds")
	}

	s, err := seeder.Seeds()
	if err != nil {
		b.Logger.Error("failed to fetch seed values",
			"error", err,
		)

		return nil, err
	}

	r, err := dynamic.Map(s,
		(**vagrant_plugin_sdk.Args_Seeds)(nil),
		argmapper.Typed(ctx, b.Internal(), b.Logger),
		argmapper.ConverterFunc(b.Mappers...),
	)

	if err != nil {
		b.Logger.Error("failed to convert seed values into proto message",
			"values", s,
			"error", err,
		)

		return nil, err
	}

	return r.(*vagrant_plugin_sdk.Args_Seeds), nil
}
