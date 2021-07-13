package dynamic

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/hashicorp/go-argmapper"

	"github.com/hashicorp/vagrant-plugin-sdk/proto/vagrant_plugin_sdk"
)

// Wraps an argmapper Func with a FuncSpec
// definition to allow easy reuse with
// wrapped clients
type SpecAndFunc struct {
	Func *argmapper.Func
	Spec *vagrant_plugin_sdk.FuncSpec
}

// Convert a value to an expected type. Converter functions should be
// included in the args list. It is important to note that the expectedType
// is a pointer to the desired type (including interfaces). For example,
// if an `int` is wanted, the expectedType would be `(*int)(nil)`.
func Map(
	resultValue, // value to be converted
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
		argmapper.Typed(resultValue),
		argmapper.Logger(Logger))

	if err = vsOut.FromResult(callFn.Call(args...)); err != nil {
		return nil, err
	}

	return vsOut.Typed(typ).Value.Interface(), nil
}

// Calls the function provided and converts the
// result to an expected type. If no type conversion
// is required, a `false` value for the expectedType
// will return the raw interface return value.
func CallFunc(
	f interface{}, // function to call
	expectedType interface{}, // nil pointer of expected return type
	mappers []*argmapper.Func, // argmapper funcs to be used as converters
	args ...argmapper.Arg, // list of argmapper arguments
) (interface{}, error) {
	var rawFunc *argmapper.Func

	if sf, ok := f.(*SpecAndFunc); ok {
		rawFunc = sf.Func
	} else if af, ok := f.(*argmapper.Func); ok {
		rawFunc = af
	} else {
		var err error
		rawFunc, err = argmapper.NewFunc(f,
			argmapper.Logger(Logger),
		)

		if err != nil {
			return nil, err
		}
	}

	args = append(args,
		argmapper.ConverterFunc(mappers...),
		argmapper.Logger(Logger))

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
		return nil, fmt.Errorf("Failed to convert %T to %T (%s)",
			raw,
			strings.TrimPrefix(fmt.Sprintf("%T", expectedType), "*"), // remove leading * so we display actual expected type
			err.Error())
	}

	return final, nil
}
