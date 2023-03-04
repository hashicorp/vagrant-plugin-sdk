// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package types

type RawRubyValue struct {
	Source Class // Ruby source class
	Data   map[string]interface{}
}
