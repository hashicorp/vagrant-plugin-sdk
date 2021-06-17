package plugincomponent

import (
	"os"

	"github.com/hashicorp/go-hclog"
)

var ArgmapperLogger hclog.Logger = hclog.New(&hclog.LoggerOptions{
	Name:   "vagrant.plugin.argmapper",
	Level:  hclog.Error,
	Output: os.Stderr,
	Color:  hclog.AutoColor,

	// Critical that this is JSON-formatted. Since we're a plugin this
	// will enable the host to parse our logs and output them in a
	// structured way.
	JSONFormat: true,
})
