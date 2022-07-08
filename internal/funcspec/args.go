package funcspec

import (
	"github.com/hashicorp/go-argmapper"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/hashicorp/vagrant-plugin-sdk/proto/vagrant_plugin_sdk"
)

// Args is a type that will be populated with all the expected args of
// the FuncSpec. This can be used in the callback (cb) to Func.
type Args []*vagrant_plugin_sdk.FuncSpec_Value

// appendValue appends an argmapper.Value to Args. The Value must
// be an *anypb.Any.
func appendValue(args Args, v argmapper.Value) Args {
	return append(args, &vagrant_plugin_sdk.FuncSpec_Value{
		Name:  v.Name,
		Type:  v.Subtype,
		Value: v.Value.Interface().(*anypb.Any),
	})
}
