// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package core

type Downloader interface {
	// component.Configurable

	Download() error
}
