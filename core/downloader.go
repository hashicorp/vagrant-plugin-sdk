// Copyright IBM Corp. 2020, 2023
// SPDX-License-Identifier: MPL-2.0

package core

type Downloader interface {
	// component.Configurable

	Download() error
}
