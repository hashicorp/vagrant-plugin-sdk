package dynamic

import (
	"fmt"
	"reflect"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/known/anypb"
)

// Decode an Any proto and return the type and value
func DecodeAny(
	input *anypb.Any, // value to decode
) (t reflect.Type, r interface{}, err error) {
	name := input.MessageName()

	typ, err := protoregistry.GlobalTypes.FindMessageByName(name)
	if err != nil {
		return t, nil, fmt.Errorf("cannot decode type: %s (%s)", name, err)
	}

	// Allocate the message type. If it is a pointer we want to
	// allocate the actual structure and not the pointer to the structure.
	v := typ.New()
	if err := input.UnmarshalTo(v.Interface().(proto.Message)); err != nil {
		return t, nil, err
	}
	r = v.Interface()
	t = reflect.TypeOf(r)

	return
}

// Encode a proto message to Any
func EncodeAny(
	input proto.Message, // proto to encode
) (*anypb.Any, error) {
	return anypb.New(input)
}
