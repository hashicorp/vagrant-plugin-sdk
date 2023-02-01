// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package dynamic

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/hashicorp/go-argmapper"
	"github.com/hashicorp/go-hclog"
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

	callFn, err := argmapper.BuildFunc(vsIn, vsOut, cb,
		argmapper.FuncName("Value mapping -> "+strings.TrimPrefix(fmt.Sprintf("%T", expectedType), "*")))
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
// generic interface may be expected (like proto.Message).
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

// This converts a value to another value using the provided mappers
// without any required type information. It is similar to the UnknownMap
// function but the expected type is not required. Due to this missing
// information, however, this function will be slower as it iterates
// the entire mapper list provided and attempts to call any function
// who's input signature includes a matching type for the given value.
func BlindMap(
	value interface{}, // value to map
	mappers []*argmapper.Func, // list of mappers to utilize
	args ...argmapper.Arg, // any extra argument to use when mapping
) (interface{}, error) {
	if value == nil {
		return nil, nil
	}

	t := reflect.TypeOf(value)
	for _, m := range mappers {
		for _, typ := range m.Input().Signature() {
			if t == typ || t.AssignableTo(typ) {
				margs := append(args, argmapper.Typed(value), argmapper.Logger(Logger))
				r := m.Call(margs...)
				if r.Err() == nil {
					return r.Out(0), nil
				} else {
					hclog.L().Info("failed to run blind map function", "source-type", hclog.Fmt("%T", value), "error", r.Err())
				}
			}
		}
	}

	return nil, fmt.Errorf("failed to map type (%T)", value)
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
	return BlindMap(input, WellKnownTypeFns)
}
