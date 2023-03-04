// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

// Package core has the interfaces for all core types
// used by Vagrant. These are implemented within Vagrant
// and provided to plugins as needed/required. This is
// currently a high level mapping of the interface
// provided by Vagrant in its current state. By doing
// a best effort match (and mapping where things are
// different) it should be easier to integrated with
// Vagrant to keep existing plugins working as well as
// making porting plugins less time consuming.

package core

// This Symbol type represents a Symbol type in Ruby.
// It is required for interoperability between legacy
// Vagrant and Go Vagrant. It's primary function is to
// allow config maps from Ruby that contain Symbols to
// be interpreted in Go while retaining type information
type Symbol string
