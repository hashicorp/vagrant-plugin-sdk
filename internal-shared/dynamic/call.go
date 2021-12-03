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
		rawFunc, err = argmapper.NewFunc(f)

		if err != nil {
			return nil, err
		}
	}

	args = append([]argmapper.Arg{
		argmapper.FuncName("Call -> " + strings.TrimPrefix(fmt.Sprintf("%T", expectedType), "*")),
		argmapper.ConverterFunc(mappers...),
		argmapper.Logger(Logger.Named("call")),
	}, args...)

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
