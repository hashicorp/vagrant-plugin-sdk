package dynamic

import (
	"fmt"
	"reflect"

	"github.com/hashicorp/go-argmapper"
	"google.golang.org/protobuf/proto"
)

var WellKnownTypeFns []*argmapper.Func

// Convert a value to an expected type. Converter functions should be
// included in the args list. It is important to note that the expectedType
// is a pointer to the desired type (including interfaces). For example,
// if an `int` is wanted, the expectedType would be `(*int)(nil)`.
func Map(
	value, // value to be converted
	expectedType interface{}, // nil pointer of desired type
	args ...argmapper.Arg, // list of argmapper arguments (including converter funcs)
) (interface{}, error) {
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

	args = append(args,
		argmapper.Typed(value),
		argmapper.Logger(Logger.Named("map")))

	if err = vsOut.FromResult(callFn.Call(args...)); err != nil {
		return nil, err
	}

	return vsOut.Typed(typ).Value.Interface(), nil
}

// Convert a value to another value using provided mappers. This
// can be useful for converting a received value (like a proto)
// into an internal value. It works by searching the list of provided
// mappers for a functions that support an input argument of the
// type that matches (or satisfies) the given value. Extra args
// can be provided which are used when applying the mapping to
// convert the value. The expected type can be useful where a
// generic interface may be expected (like proto.Message). If
// the expected type is unknown, setting as an interface{} nil
// pointer will effectively disable the type check.
func UnknownMap(
	value, // value to map
	expectedType interface{}, // nil pointer of desired type
	mappers []*argmapper.Func, // list of mappers to utilize
	args ...argmapper.Arg, // any extra arguments to use when mapping
) (interface{}, error) {
	// If the value provided is nil, we are already done!
	if value == nil {
		return nil, nil
	}

	t := reflect.TypeOf(value)
	maps := []*argmapper.Func{}
	for _, m := range mappers {
		for _, typ := range m.Input().Signature() {
			if t == typ || t.AssignableTo(typ) {
				maps = append(maps, m)
				break
			}
		}
	}

	margs := append(args, argmapper.ConverterFunc(maps...))
	v, err := Map(value, expectedType, margs...)
	if err != nil {
		return nil, err
	}

	return v, nil
}

// Map a well known value to proto
func MapToWellKnownProto(
	input interface{}, // value to be converted to proto
) (proto.Message, error) {
	v, err := Map(input, (*proto.Message)(nil),
		argmapper.ConverterFunc(WellKnownTypeFns...),
	)
	if err != nil {
		return nil, err
	}

	return v.(proto.Message), nil
}

// Map a well known proto to value
func MapFromWellKnownProto(
	input proto.Message, // proto to be converted to value
) (interface{}, error) {
	v, err := Map(input, (*interface{})(nil),
		argmapper.ConverterFunc(WellKnownTypeFns...),
	)
	if err != nil {
		return nil, err
	}

	return v, nil
}
