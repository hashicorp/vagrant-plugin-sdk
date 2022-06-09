package funcspec

import (
	"fmt"
	"reflect"

	"github.com/hashicorp/go-argmapper"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

// anyConvGen is an argmapper.ConverterGenFunc that dynamically creates
// converters to *anypb.Any for types that implement proto.Message. This
// allows automatic conversion to *anypb.Any.
//
// This is automatically injected for all funcspec.Func calls.
func anyConvGen(v argmapper.Value) (*argmapper.Func, error) {
	anyType := reflect.TypeOf((*anypb.Any)(nil))
	protoMessageType := reflect.TypeOf((*proto.Message)(nil)).Elem()
	if !v.Type.Implements(protoMessageType) {
		return nil, nil
	}

	// We take this value as our input.
	inputSet, err := argmapper.NewValueSet([]argmapper.Value{v})
	if err != nil {
		return nil, err
	}

	// Generate an int with the subtype of the string value
	outputSet, err := argmapper.NewValueSet([]argmapper.Value{{
		Name:    v.Name,
		Type:    anyType,
		Subtype: string(reflect.Zero(v.Type).Interface().(proto.Message).ProtoReflect().Descriptor().FullName()),
	}})
	if err != nil {
		return nil, err
	}

	return argmapper.BuildFunc(inputSet, outputSet, func(in, out *argmapper.ValueSet) error {
		inputVal := inputSet.Typed(v.Type)
		// If there is no typed input, check the named inputs
		if inputVal == nil {
			inputVal = inputSet.Named(v.Name)
		}
		anyVal, err := anypb.New(inputVal.Value.Interface().(proto.Message))
		if err != nil {
			return err
		}

		// If there is no typed output, check the named inputs
		outputVal := outputSet.Typed(anyType)
		if outputVal == nil {
			outputVal = outputSet.Named(v.Name)
		}
		outputVal.Value = reflect.ValueOf(anyVal)
		return nil
	}, argmapper.FuncName(fmt.Sprintf("converter: %s -> *anypb.Any", v.Type)))

}
