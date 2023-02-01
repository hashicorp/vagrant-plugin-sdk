// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package core

import "github.com/hashicorp/vagrant-plugin-sdk/component"

type Config interface {
	Register() (*component.ConfigRegistration, error)
	Struct() (interface{}, error)
	Merge(base, toMerge *component.ConfigData) (merged *component.ConfigData, err error)
	Finalize(*component.ConfigData) (*component.ConfigData, error)
}
