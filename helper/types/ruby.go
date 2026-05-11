// Copyright IBM Corp. 2020, 2026
// SPDX-License-Identifier: MPL-2.0

package types

type RawRubyValue struct {
	Source Class // Ruby source class
	Data   map[string]interface{}
}
