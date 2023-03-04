// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package core

type Named interface {
	SetPluginName(string) error
	PluginName() (name string, err error)
}
