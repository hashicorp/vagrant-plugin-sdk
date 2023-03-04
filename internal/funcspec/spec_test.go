// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package funcspec

import (
	"reflect"
	"testing"

	"github.com/hashicorp/go-argmapper"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/hashicorp/go-hclog"
)

func init() {
	hclog.L().SetLevel(hclog.Trace)
}

func TestSpec(t *testing.T) {
	t.Run("proto to proto", func(t *testing.T) {
		require := require.New(t)

		spec, err := Spec(func(*emptypb.Empty) *emptypb.Empty { return nil })
		require.NoError(err)
		require.NotNil(spec)
		require.Len(spec.Args, 1)
		require.Empty(spec.Args[0].Name)
		require.Equal("google.protobuf.Empty", spec.Args[0].Type)
		require.Len(spec.Result, 1)
		require.Empty(spec.Result[0].Name)
		require.Equal("google.protobuf.Empty", spec.Result[0].Type)
	})

	t.Run("converted args to proto", func(t *testing.T) {
		require := require.New(t)

		type Foo struct{}

		spec, err := Spec(func(*Foo) *emptypb.Empty { return nil },
			argmapper.Converter(func(*emptypb.Empty) *Foo { return nil }),
		)
		require.NoError(err)
		require.NotNil(spec)
		require.Len(spec.Args, 1)
		require.Empty(spec.Args[0].Name)
		require.Equal("google.protobuf.Empty", spec.Args[0].Type)
		require.Len(spec.Result, 1)
		require.Empty(spec.Result[0].Name)
		require.Equal("google.protobuf.Empty", spec.Result[0].Type)
	})

	t.Run("unsatisfied conversion", func(t *testing.T) {
		require := require.New(t)

		type Foo struct{}
		type Bar struct{}

		spec, err := Spec(func(*Foo) *emptypb.Empty { return nil },
			argmapper.Converter(func(*emptypb.Empty) *Bar { return nil }),
		)
		require.Error(err)
		require.Nil(spec)
	})

	t.Run("proto to int", func(t *testing.T) {
		require := require.New(t)

		spec, err := Spec(func(*emptypb.Empty) int { return 0 })
		require.NoError(err)
		require.NotNil(spec)
	})

	t.Run("WithOutput proto to interface, doesn't implement", func(t *testing.T) {
		require := require.New(t)

		spec, err := Spec(func(*emptypb.Empty) struct{} { return struct{}{} },
			argmapper.FilterOutput(argmapper.FilterType(reflect.TypeOf((*testSpecInterface)(nil)).Elem())),
		)
		require.Error(err)
		require.Nil(spec)
	})

	t.Run("args as extra values", func(t *testing.T) {
		require := require.New(t)

		type Foo struct{}
		type Bar struct{}

		spec, err := Spec(func(*Foo, *Bar) *emptypb.Empty { return nil },
			argmapper.Converter(func(*emptypb.Empty) *Foo { return nil }),
			argmapper.Typed(&Bar{}),
		)
		require.NoError(err)
		require.NotNil(spec)
		require.Len(spec.Args, 1)
		require.Empty(spec.Args[0].Name)
		require.Equal("google.protobuf.Empty", spec.Args[0].Type)
		require.Len(spec.Result, 1)
		require.Empty(spec.Result[0].Name)
		require.Equal("google.protobuf.Empty", spec.Result[0].Type)
	})
}

type testSpecInterface interface {
	hello()
}

type testSpecInterfaceImpl struct{}

func (testSpecInterfaceImpl) hello() {}
