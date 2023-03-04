// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package core

type StateBag interface {
	Get(string) interface{}
	GetOk(string) (interface{}, bool)
	Put(string, interface{})
	Remove(string)
}
