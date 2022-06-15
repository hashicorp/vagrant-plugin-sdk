package funcspec

import (
	"context"
	"fmt"
	"reflect"

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

	inputs := []argmapper.Value{}

	// Take each input for the defined function and redefine it
	// based on the filter to get valid proto inputs. We do this
	// individually so that we can ensure named parameters are
	// properly retained
	for _, i := range f.Input().Values() {
		// Create a function whose input and output is the same type
		// as the current input being processed
		set, err := argmapper.NewValueSet([]argmapper.Value{i})
		if err != nil {
			return nil, err
		}
		cb := func(in, out *argmapper.ValueSet) error { return nil }
		inputFn, err := argmapper.BuildFunc(set, set, cb, argmapper.Logger(dynamic.Logger))
		if err != nil {
			return nil, err
		}

		// Now redefine the function including our input filter
		reFn, err := inputFn.Redefine(
			append(
				args,
				argmapper.FilterInput(inputFilter),
			)...,
		)
		if err != nil {
			return nil, err
		}

		// Collect all the inputs for the redefined function
		// that are protos
		reI := []argmapper.Value{}
		for _, ri := range reFn.Input().Values() {
			if !filterProto(ri) {
				continue
			}
			reI = append(reI, ri)
		}

		// If there are no inputs, then the original function
		// input is being satisfied by provided arguments
		if len(reI) == 0 {
			continue
		}

		// If we have more than one input, this is unexpected and
		// we force an error.
		//
		// TODO(spox): This seems to be true for existing cases, but
		// 			   it's easy to imagine a situation where the
		// 			   conversion could result in extra inputs.
		if len(reI) != 1 {
			return nil, fmt.Errorf(
				"expected funcspec input redefine size to be 1, got: %d (fn: %s)",
				f.Name(),
				len(reI),
			)
		}

		// Grab the new input, set the name value and store as a new input
		v := reI[0]
		v.Name = i.Name
		inputs = append(inputs, v)
	}

	// Grab the input set of the function and build up our funcspec
	result := vagrant_plugin_sdk.FuncSpec{Name: f.Name()}
	for _, v := range inputs {
		if !filterProto(v) {
			continue
		}

		// if inputFilter(v) && v.Name != "" {
		// 	result.Args = append(result.Args, &vagrant_plugin_sdk.FuncSpec_Value{
		// 		Name: v.Name,
		// 	})
		// 	continue
		// }

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
	cliOptType       = reflect.TypeOf((**component.CommandFlag)(nil)).Elem()
	commandInfoType  = reflect.TypeOf((**component.CommandInfo)(nil)).Elem()
	configDataType   = reflect.TypeOf((**component.ConfigData)(nil)).Elem()
	interfaceType    = reflect.TypeOf((*interface{})(nil)).Elem()
)
