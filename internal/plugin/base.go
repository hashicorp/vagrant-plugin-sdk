package plugin

import (
	"context"
	"fmt"
	"reflect"

	"github.com/hashicorp/go-argmapper"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"

	"github.com/hashicorp/vagrant-plugin-sdk/internal/funcspec"
	"github.com/hashicorp/vagrant-plugin-sdk/internal/pluginargs"
	pb "github.com/hashicorp/vagrant-plugin-sdk/proto/gen"
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
}

type baseClient struct {
	*base
}

type baseServer struct {
	*base
}

// internal returns a new pluginargs.Internal that can be used with
// dynamic calls. The Internal structure is an internal-only argument
// that is used to perform cleanup.
func (b *base) internal() *pluginargs.Internal {
	return &pluginargs.Internal{
		Broker:  b.Broker,
		Mappers: b.Mappers,
		Cleanup: &pluginargs.Cleanup{},
	}
}

func (b *baseClient) callRemoteDynamicFunc(
	ctx context.Context,
	mappers []*argmapper.Func,
	result interface{}, // expected result type
	f interface{}, // function
	args ...argmapper.Arg,
) (interface{}, error) {
	// We allow f to be a *mapper.Func because our plugin system creates
	// a func directly due to special argument types.
	// TODO: test
	rawFunc, ok := f.(*argmapper.Func)
	if !ok {
		var err error
		rawFunc, err = argmapper.NewFunc(f, argmapper.Logger(b.Logger))
		if err != nil {
			return nil, err
		}
	}

	// Make sure we have access to our context and logger and default args
	args = append(args,
		argmapper.ConverterFunc(b.Mappers...),
		argmapper.ConverterFunc(mappers...),
		argmapper.Typed(
			ctx,
			b.Logger,
		),

		// argmapper.Named("labels", &component.LabelSet{Labels: c.labels}),
	)

	// Build the chain and call it
	callResult := rawFunc.Call(args...)
	if err := callResult.Err(); err != nil {
		return nil, err
	}
	raw := callResult.Out(0)

	// If we don't have an expected result type, then just return as-is.
	// Otherwise, we need to verify the result type matches properly.
	if result == nil {
		return raw, nil
	}

	// Verify
	interfaceType := reflect.TypeOf(result).Elem()
	if rawType := reflect.TypeOf(raw); !rawType.Implements(interfaceType) {
		return nil, status.Errorf(codes.FailedPrecondition,
			"operation expected result type %s, got %s",
			interfaceType.String(),
			rawType.String())
	}

	return raw, nil
}

func (b *baseClient) generateFunc(spec *pb.FuncSpec, cbFn interface{}, args ...argmapper.Arg) interface{} {
	return funcspec.Func(spec, cbFn, append(args,
		argmapper.Logger(b.Logger),
		argmapper.Typed(b.internal()))...,
	)
}

func (b *baseServer) callLocalDynamicFunc(
	f interface{},
	args funcspec.Args,
	result interface{}, // expected result type
	callArgs ...argmapper.Arg,
) (interface{}, error) {
	internal := b.internal()
	defer internal.Cleanup.Close()

	callArgs = append(callArgs,
		argmapper.ConverterFunc(b.Mappers...),
		argmapper.Logger(b.Logger),
		argmapper.Typed(internal),
	)

	// Decode our *any.Any values.
	for _, arg := range args {
		anyVal := arg.Value

		name, err := ptypes.AnyMessageName(anyVal)
		if err != nil {
			return nil, err
		}

		typ := proto.MessageType(name)
		if typ == nil {
			return nil, fmt.Errorf("cannot decode type: %s", name)
		}

		// Allocate the message type. If it is a pointer we want to
		// allocate the actual structure and not the pointer to the structure.
		if typ.Kind() == reflect.Ptr {
			typ = typ.Elem()
		}
		v := reflect.New(typ)
		v.Elem().Set(reflect.Zero(typ))

		// Unmarshal directly into our newly allocated structure.
		if err := ptypes.UnmarshalAny(anyVal, v.Interface().(proto.Message)); err != nil {
			return nil, err
		}

		callArgs = append(callArgs,
			argmapper.NamedSubtype(arg.Name, v.Interface(), arg.Type),
		)
	}

	mapF, err := argmapper.NewFunc(f)
	if err != nil {
		return nil, err
	}

	callResult := mapF.Call(callArgs...)
	if err := callResult.Err(); err != nil {
		return nil, err
	}

	raw := callResult.Out(0)

	// TODO: Maybe this should just assert that the result
	// implements proto.Message
	// Verify
	interfaceType := reflect.TypeOf(result).Elem()
	if rawType := reflect.TypeOf(raw); !rawType.Implements(interfaceType) {
		return nil, status.Errorf(codes.FailedPrecondition,
			"operation expected result type %s, got %s",
			interfaceType.String(),
			rawType.String())
	}
	return raw, nil
}

func (b *baseServer) generateSpec(fn interface{}, args ...argmapper.Arg) (*pb.FuncSpec, error) {
	return funcspec.Spec(fn, append(args,
		argmapper.Logger(b.Logger),
		argmapper.ConverterFunc(b.Mappers...),
		argmapper.Typed(b.internal()))...,
	)
}
