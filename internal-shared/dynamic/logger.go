package dynamic

import (
	"os"

	"github.com/hashicorp/go-hclog"
)

var Logger hclog.Logger = hclog.New(&hclog.LoggerOptions{
	Name:       "vagrant.plugin.argmapper",
	Level:      hclog.Error,
	Output:     os.Stderr,
	Color:      hclog.AutoColor,
	JSONFormat: false,
})
