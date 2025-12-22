// Copyright IBM Corp. 2020, 2023
// SPDX-License-Identifier: MPL-2.0

package core

type Push interface {
	Push() (err error)
}
