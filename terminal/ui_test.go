// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package terminal

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNamedValues(t *testing.T) {
	require := require.New(t)

	var buf bytes.Buffer
	var ui basicUI
	ui.NamedValues([]NamedValue{
		{"hello", "a"},
		{"this", "is"},
		{"a", "test"},
		{"of", "foo"},
		{"the_key_value", "style"},
	},
		WithWriter(&buf),
	)

	expected := `
          hello: a
           this: is
              a: test
             of: foo
  the_key_value: style

`

	require.Equal(strings.TrimLeft(expected, "\n"), buf.String())
}

func TestNamedValues_server(t *testing.T) {
	require := require.New(t)

	var buf bytes.Buffer
	var ui basicUI
	ui.Output("Server configuration:", WithHeaderStyle(), WithWriter(&buf))
	ui.NamedValues([]NamedValue{
		{"DB Path", "data.db"},
		{"gRPC Address", "127.0.0.1:1234"},
		{"HTTP Address", "127.0.0.1:1235"},
		{"URL Service", "api.alpha.vagrant.run:443 (account: token)"},
	},
		WithWriter(&buf),
	)

	expected := `
==> Server configuration:
       DB Path: data.db
  gRPC Address: 127.0.0.1:1234
  HTTP Address: 127.0.0.1:1235
   URL Service: api.alpha.vagrant.run:443 (account: token)

`

	require.Equal(expected, buf.String())
}

func TestStatusStyle(t *testing.T) {
	require := require.New(t)

	var buf bytes.Buffer
	var ui basicUI
	ui.Output(strings.TrimSpace(`
one
two
  three`),
		WithWriter(&buf),
		WithInfoStyle(),
	)

	expected := `  one
  two
    three
`

	require.Equal(expected, buf.String())
}

func TestInterpret(t *testing.T) {
	require := require.New(t)

	type interpretTest struct {
		interpretOptions        []interface{}
		expectedStyle           string
		expectedDisabledNewLine bool
		expectedColor           string
	}

	tests := []*interpretTest{
		{
			interpretOptions:        []interface{}{WithSuccessBoldStyle(), WithoutNewLine()},
			expectedStyle:           SuccessBoldStyle,
			expectedDisabledNewLine: true,
			expectedColor:           "green",
		},
		{
			interpretOptions:        []interface{}{WithSuccessStyle(), WithoutNewLine(), WithColor("magenta")},
			expectedStyle:           SuccessStyle,
			expectedDisabledNewLine: true,
			expectedColor:           "magenta",
		},
		{
			interpretOptions:        []interface{}{WithColor("magenta"), WithSuccessStyle(), WithoutNewLine()},
			expectedStyle:           SuccessStyle,
			expectedDisabledNewLine: true,
			expectedColor:           "magenta",
		},
		{
			interpretOptions:        []interface{}{WithStyle("mystyle"), WithNewLine(), WithColor("lightRed")},
			expectedStyle:           "mystyle",
			expectedDisabledNewLine: false,
			expectedColor:           "lightRed",
		},
	}

	testMessage := "This is a message"
	testWriter := os.Stdout
	for _, test := range tests {
		test.interpretOptions = append(test.interpretOptions, WithWriter(testWriter))
		msg, style, disableNewLine, writer, color := Interpret(
			testMessage, test.interpretOptions...,
		)

		require.Contains(msg, testMessage)
		require.Equal(writer, os.Stdout)

		require.Equal(style, test.expectedStyle)
		require.Equal(disableNewLine, test.expectedDisabledNewLine)
		require.Equal(color, test.expectedColor)
	}
}
