package funcspec

import (
	"github.com/golang/protobuf/ptypes/any"
	"github.com/hashicorp/go-argmapper"

	"github.com/hashicorp/vagrant-plugin-sdk/proto/vagrant_plugin_sdk"
)

// Args is a type that will be populated with all the expected args of
// the FuncSpec. This can be used in the callback (cb) to Func.
type Args []*vagrant_plugin_sdk.FuncSpec_Value

// appendValue appends an argmapper.Value to Args. The Value must
// be an *any.Any.
func appendValue(args Args, v argmapper.Value) Args {
	return append(args, &vagrant_plugin_sdk.FuncSpec_Value{
		Name:  v.Name,
		Type:  v.Subtype,
		Value: v.Value.Interface().(*any.Any),
	})
}
