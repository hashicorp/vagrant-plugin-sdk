// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package funcspec

import (
	"context"
	"reflect"

	"github.com/hashicorp/go-argmapper"
	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"

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

	args = append([]argmapper.Arg{
		argmapper.FilterInput(filterAllowedInputTypes),
		//		argmapper.FilterOutput(filterAllowedOutputTypes),
	}, args...)

	// Build our function
	f, err := argmapper.NewFunc(fn,
		argmapper.Logger(dynamic.Logger))
	if err != nil {
		return nil, err
	}

	// Grab the input set of the function and build up our funcspec
	result := vagrant_plugin_sdk.FuncSpec{Name: f.Name()}
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
		inputFn, err := argmapper.BuildFunc(set, set, cb,
			argmapper.Logger(dynamic.Logger))
		if err != nil {
			return nil, err
		}

		// Now redefine the function with our input and output
		// filters included
		reFn, err := inputFn.Redefine(args...)
		if err != nil {
			return nil, err
		}

		// Collect all the inputs for the redefined function
		// that are protos
		for _, reVal := range reFn.Input().Values() {
			if !filterProto(reVal) {
				continue
			}
			// Set the name on the value. If extra values are
			// required, the name won't have an impact as it
			// will fallback to a typed match
			result.Args = append(result.Args,
				&vagrant_plugin_sdk.FuncSpec_Value{
					Name: i.Name,
					Type: typeToMessage(reVal.Type),
				},
			)
		}
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
	return string(reflect.Zero(typ).Interface().(proto.Message).ProtoReflect().Descriptor().FullName())
}

var (
	filterAllowedInputTypes = argmapper.FilterOr(
		argmapper.FilterType(reflect.TypeOf((*context.Context)(nil)).Elem()),
		argmapper.FilterType(reflect.TypeOf((*proto.Message)(nil)).Elem()),
		argmapper.FilterType(reflect.TypeOf((*string)(nil)).Elem()),
	)

	filterAllowedOutputTypes = argmapper.FilterOr(
		argmapper.FilterType(reflect.TypeOf((*proto.Message)(nil)).Elem()),
		argmapper.FilterType(reflect.TypeOf((*bool)(nil)).Elem()),
		argmapper.FilterType(reflect.TypeOf((*string)(nil)).Elem()),
		argmapper.FilterType(reflect.TypeOf((*int64)(nil)).Elem()),
		argmapper.FilterType(reflect.TypeOf((**component.CommandFlag)(nil)).Elem()),
		argmapper.FilterType(reflect.TypeOf((**component.CommandInfo)(nil)).Elem()),
		argmapper.FilterType(reflect.TypeOf((*hclog.Logger)(nil)).Elem()),
		//		argmapper.FilterType(reflect.TypeOf((*interface{})(nil)).Elem()),
	)

	filterProto = argmapper.FilterType(reflect.TypeOf((*proto.Message)(nil)).Elem())
)
