// Copyright IBM Corp. 2020, 2023
// SPDX-License-Identifier: MPL-2.0

package dynamic

import (
	"os"

	"github.com/hashicorp/go-hclog"
)

func init() {
	level := hclog.Error
	if os.Getenv("VAGRANT_LOG_ARGMAPPER") != "" {
		level = hclog.Trace
	}
	Logger = hclog.New(&hclog.LoggerOptions{
		Name:       "vagrant.plugin.argmapper",
		Level:      level,
		Output:     os.Stderr,
		Color:      hclog.AutoColor,
		JSONFormat: false,
	})
}

var Logger hclog.Logger
