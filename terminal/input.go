// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package terminal

// Input is the configuration for an input.
type Input struct {
	// Prompt is a single-line prompt to give the user such as "Continue?"
	// The user will input their answer after this prompt.
	Prompt string

	// Style is the style to apply to the input. If this is blank,
	// the output won't be styled in any way.
	Style string

	// True if this input is a secret. The input will be masked.
	Secret bool

	// Color is the color to apply to the input. If this is blank,
	// the output will be the default color for the terminal.
	Color string
}
