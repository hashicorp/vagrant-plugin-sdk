package funcspec

import (
	"context"
	"reflect"

	"github.com/DavidGamba/go-getoptions/option"
	"github.com/golang/protobuf/proto"
	"github.com/hashicorp/go-argmapper"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hashicorp/vagrant-plugin-sdk/component"
	"github.com/hashicorp/vagrant-plugin-sdk/internal-shared/dynamic"
	"github.com/hashicorp/vagrant-plugin-sdk/proto/vagrant_plugin_sdk"
)

func ArgSpec(f *argmapper.Func, args ...argmapper.Arg) (*vagrant_plugin_sdk.FuncSpec, error) {
	// Grab the input set of the function and build up our funcspec
	result := vagrant_plugin_sdk.FuncSpec{Name: f.Name()}
	for _, v := range f.Input().Values() {
		if reflect.Zero(v.Type).Interface() == nil {
			continue
		}
		result.Args = append(result.Args, &vagrant_plugin_sdk.FuncSpec_Value{
			Name: v.Name,
			Type: v.Subtype,
		})
	}

	// Grab the output set and store that
	for _, v := range f.Output().Values() {
		result.Result = append(result.Result, &vagrant_plugin_sdk.FuncSpec_Value{
			Name: v.Name,
			Type: typeToMessage(v.Type),
		})
	}

	return &result, nil
}

// Spec takes a function pointer and generates a FuncSpec from it. The
// function must only take arguments that are proto.Message implementations
// or have a chain of converters that directly convert to a proto.Message.
func Spec(fn interface{}, args ...argmapper.Arg) (*vagrant_plugin_sdk.FuncSpec, error) {
	if fn == nil {
		return nil, status.Errorf(codes.Unimplemented, "required plugin type not implemented")
	}

	filterProto := argmapper.FilterType(protoMessageType)

	outputFilter := argmapper.FilterOr(
		filterProto,
		argmapper.FilterType(boolType),
		argmapper.FilterType(stringType),
		argmapper.FilterType(intType),
		argmapper.FilterType(cliOptType),
		argmapper.FilterType(commandInfoType),
		argmapper.FilterType(interfaceType),
	)
	// Copy our args cause we're going to use append() and we don't
	// want to modify our caller.
	args = append([]argmapper.Arg{
		argmapper.FilterOutput(outputFilter),
	}, args...)

	// Build our function
	f, err := argmapper.NewFunc(fn,
		argmapper.Logger(dynamic.Logger))
	if err != nil {
		return nil, err
	}

	inputFilter := argmapper.FilterOr(
		argmapper.FilterType(contextType),
		argmapper.FilterType(stringType),
		filterProto,
	)

	// Redefine the function in terms of protobuf messages. "Redefine" changes
	// the inputs of a function to only require values that match our filter
	// function. In our case, that is protobuf messages.
	f, err = f.Redefine(append(args,
		argmapper.FilterInput(inputFilter),
	)...)
	if err != nil {
		return nil, err
	}

	// Grab the input set of the function and build up our funcspec
	result := vagrant_plugin_sdk.FuncSpec{Name: f.Name()}
	for _, v := range f.Input().Values() {
		if !filterProto(v) {
			if inputFilter(v) && v.Name != "" {
				result.Args = append(result.Args, &vagrant_plugin_sdk.FuncSpec_Value{
					Name: v.Name,
				})
			}
			continue
		}

		result.Args = append(result.Args, &vagrant_plugin_sdk.FuncSpec_Value{
			Name: v.Name,
			Type: typeToMessage(v.Type),
		})
	}

	// Grab the output set and store that
	for _, v := range f.Output().Values() {
		// We only advertise proto types in output since those are the only
		// types we can send across the plugin boundary.
		if !filterProto(v) {
			continue
		}

		result.Result = append(result.Result, &vagrant_plugin_sdk.FuncSpec_Value{
			Name: v.Name,
			Type: typeToMessage(v.Type),
		})
	}

	return &result, nil
}

func typeToMessage(typ reflect.Type) string {
	return proto.MessageName(reflect.Zero(typ).Interface().(proto.Message))
}

var (
	contextType      = reflect.TypeOf((*context.Context)(nil)).Elem()
	protoMessageType = reflect.TypeOf((*proto.Message)(nil)).Elem()
	boolType         = reflect.TypeOf((*bool)(nil)).Elem()
	stringType       = reflect.TypeOf((*string)(nil)).Elem()
	intType          = reflect.TypeOf((*int64)(nil)).Elem()
	cliOptType       = reflect.TypeOf((*[]*option.Option)(nil)).Elem()
	commandInfoType  = reflect.TypeOf((**component.CommandInfo)(nil)).Elem()
	interfaceType    = reflect.TypeOf((*interface{})(nil)).Elem()
)
