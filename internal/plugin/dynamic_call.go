package plugin

import (
	"fmt"
	"reflect"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/hashicorp/go-argmapper"

	"github.com/hashicorp/vagrant-plugin-sdk/component"
	"github.com/hashicorp/vagrant-plugin-sdk/internal/funcspec"
	"github.com/hashicorp/vagrant-plugin-sdk/internal/plugincomponent"
)

func Map(resultValue, expectedType interface{}, args ...argmapper.Arg) (interface{}, error) {
	typPtr := reflect.TypeOf(expectedType)
	if typPtr.Kind() != reflect.Ptr {
		return nil, fmt.Errorf("expectedType must be nil pointer")
	}
	typ := typPtr.Elem()

	vIn := argmapper.Value{Type: typ}
	vOut := argmapper.Value{Type: typ}
	vsIn, err := argmapper.NewValueSet([]argmapper.Value{vIn})
	if err != nil {
		return nil, err
	}
	vsOut, err := argmapper.NewValueSet([]argmapper.Value{vOut})
	if err != nil {
		return nil, err
	}

	cb := func(in, out *argmapper.ValueSet) error {
		val := in.Typed(typ).Value.Interface()
		out.Typed(typ).Value = reflect.ValueOf(val)
		return nil
	}

	callFn, err := argmapper.BuildFunc(vsIn, vsOut, cb)
	if err != nil {
		return nil, err
	}

	args = append(args, argmapper.Typed(resultValue))

	if err = vsOut.FromResult(callFn.Call(args...)); err != nil {
		return nil, err
	}

	return vsOut.Typed(typ).Value.Interface(), nil
}

func callDynamicFunc(
	f interface{}, // function
	expectedType interface{}, // expected result type
	mappers []*argmapper.Func,
	args ...argmapper.Arg,
) (interface{}, error) {
	var rawFunc *argmapper.Func

	if sf, ok := f.(*component.SpicyFunc); ok {
		rawFunc = sf.Func
	} else if af, ok := f.(*argmapper.Func); ok {
		rawFunc = af
	} else {
		var err error
		rawFunc, err = argmapper.NewFunc(f,
			argmapper.Logger(plugincomponent.ArgmapperLogger),
		)

		if err != nil {
			return nil, err
		}
	}

	// Make sure we have access to our context and logger and default args
	args = append(args, argmapper.ConverterFunc(mappers...))

	// Build the chain and call it
	callResult := rawFunc.Call(args...)
	if err := callResult.Err(); err != nil {
		return nil, err
	}
	raw := callResult.Out(0)

	// If a false value is passed as the expectedType, then
	// no validation is performed and we just return the value
	// that we got
	typPtr := reflect.TypeOf(expectedType)
	if typPtr.Kind() == reflect.Bool && expectedType.(bool) == false {
		return raw, nil
	}

	final, err := Map(raw, expectedType, args...)
	if err != nil {
		return nil, fmt.Errorf("Failed to convert %T to %T (%s)", raw, expectedType, err.Error())
	}

	return final, nil
}

// callDynamicFunc calls a dynamic (mapper-based) function with the
// given input arguments. This is a helper that is expected to be used
// by most component gRPC servers to implement their function calls.
func callDynamicFunc2(
	f interface{},
	args funcspec.Args,
	callArgs ...argmapper.Arg,
) (interface{}, error) {
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

	callArgs = append(callArgs,
		argmapper.Logger(plugincomponent.ArgmapperLogger))
	mapF, err := argmapper.NewFunc(f)
	if err != nil {
		return nil, err
	}

	result := mapF.Call(callArgs...)
	if err := result.Err(); err != nil {
		return nil, err
	}

	return result.Out(0), nil
}

// callDynamicFuncAny is callDynamicFunc that automatically encodes the
// result to an *any.Any.
func callDynamicFuncAny2(
	f interface{},
	args funcspec.Args,
	callArgs ...argmapper.Arg,
) (*any.Any, interface{}, error) {
	result, err := callDynamicFunc2(f, args, callArgs...)
	if err != nil {
		return nil, nil, err
	}

	// We expect the final result to always be a proto message so we can
	// send it back over the wire.
	//
	// NOTE(mitchellh): If we wanted to in the future, we can probably change
	// this to be any type that has a mapper that can take it to be a
	// proto.Message.
	msg, ok := result.(proto.Message)
	if !ok {
		return nil, nil, fmt.Errorf(
			"result of plugin-based function must be a proto.Message, got %T", msg)
	}

	anyVal, err := ptypes.MarshalAny(msg)
	return anyVal, result, err
}
