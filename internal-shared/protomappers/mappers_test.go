// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package protomappers

import (
	"testing"

	"github.com/hashicorp/go-argmapper"
	"github.com/hashicorp/go-hclog"
	"github.com/stretchr/testify/require"

	"github.com/hashicorp/vagrant-plugin-sdk/proto/vagrant_plugin_sdk"
)

func TestMappers(t *testing.T) {
	var cases = []struct {
		Name   string
		Mapper interface{}
		Input  []interface{}
		Output interface{}
		Error  string
	}{
		{
			"Logger",
			Logger,
			[]interface{}{&vagrant_plugin_sdk.Args_Logger{Name: "foo"}},
			hclog.L().ResetNamed("foo"),
			"",
		},

		{
			"LoggerProto",
			LoggerProto,
			[]interface{}{hclog.L().ResetNamed("foo")},
			&vagrant_plugin_sdk.Args_Logger{Name: "foo"},
			"",
		},
	}

	for _, tt := range cases {
		t.Run(tt.Name, func(t *testing.T) {
			require := require.New(t)

			f, err := argmapper.NewFunc(tt.Mapper)
			require.NoError(err)

			var args []argmapper.Arg
			for _, input := range tt.Input {
				args = append(args, argmapper.Typed(input))
			}

			result := f.Call(args...)
			if tt.Error != "" {
				require.Error(result.Err())
				require.Contains(result.Err().Error(), tt.Error)
				return
			}
			require.NoError(result.Err())
			require.Equal(tt.Output, result.Out(0))
		})
	}
}
