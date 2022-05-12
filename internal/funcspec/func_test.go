package funcspec

import (
	"reflect"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/hashicorp/go-argmapper"
	"github.com/stretchr/testify/require"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vagrant-plugin-sdk/proto/vagrant_plugin_sdk"
)

func init() {
	hclog.L().SetLevel(hclog.Trace)
}

func TestFunc(t *testing.T) {
	t.Run("single any result", func(t *testing.T) {
		require := require.New(t)

		spec, err := Spec(func(*empty.Empty) *empty.Empty { return &empty.Empty{} })
		require.NoError(err)
		require.NotNil(spec)

		f := Func(spec, func(args Args) (*any.Any, error) {
			require.Len(args, 1)
			require.NotNil(args[0])

			// At this point we'd normally RPC out.
			return ptypes.MarshalAny(&empty.Empty{})
		})

		msg, err := ptypes.MarshalAny(&empty.Empty{})
		require.NoError(err)

		result := f.Func.Call(argmapper.TypedSubtype(msg, proto.MessageName(&empty.Empty{})))
		require.NoError(result.Err())
		require.Equal(reflect.Struct, reflect.ValueOf(result.Out(0)).Kind())
	})

	t.Run("single missing requirement", func(t *testing.T) {
		require := require.New(t)

		spec, err := Spec(func(*empty.Empty) *empty.Empty { return &empty.Empty{} })
		require.NoError(err)
		require.NotNil(spec)

		f := Func(spec, func(args Args) (*any.Any, error) {
			require.Len(args, 1)
			require.NotNil(args[0])

			// At this point we'd normally RPC out.
			return ptypes.MarshalAny(&empty.Empty{})
		})

		// Create an argument with the wrong type
		msg, err := ptypes.MarshalAny(&vagrant_plugin_sdk.FuncSpec{})
		require.NoError(err)
		result := f.Func.Call(argmapper.TypedSubtype(msg, proto.MessageName(&vagrant_plugin_sdk.FuncSpec{})))

		// We should have an error
		require.Error(result.Err())
		require.Contains(result.Err().Error(), "could not be satisfied")
	})

	t.Run("match callback output if no results", func(t *testing.T) {
		require := require.New(t)

		spec, err := Spec(func(*empty.Empty) *empty.Empty { return &empty.Empty{} })
		require.NoError(err)
		require.NotNil(spec)

		// No results
		spec.Result = nil

		// Build our func to return a primitive
		f := Func(spec, func(args Args) int {
			require.Len(args, 1)
			require.NotNil(args[0])
			return 42
		})

		// Call the function with the proto type we expect
		msg, err := ptypes.MarshalAny(&empty.Empty{})
		require.NoError(err)
		result := f.Func.Call(argmapper.TypedSubtype(msg, proto.MessageName(&empty.Empty{})))

		// Should succeed and give us our primitive
		require.NoError(result.Err())
		require.Equal(42, result.Out(0))
	})

	t.Run("provide input arguments", func(t *testing.T) {
		require := require.New(t)

		spec, err := Spec(func(*empty.Empty) *empty.Empty { return &empty.Empty{} })
		require.NoError(err)
		require.NotNil(spec)

		f := Func(spec, func(args Args, v int) (*any.Any, error) {
			require.Len(args, 1)
			require.NotNil(args[0])
			require.Equal(42, v)

			// At this point we'd normally RPC out.
			return ptypes.MarshalAny(&empty.Empty{})
		}, argmapper.Typed(int(42)))

		msg, err := ptypes.MarshalAny(&empty.Empty{})
		require.NoError(err)

		result := f.Func.Call(argmapper.TypedSubtype(msg, proto.MessageName(&empty.Empty{})))
		require.NoError(result.Err())
		require.Equal(reflect.Struct, reflect.ValueOf(result.Out(0)).Kind())
	})
}
