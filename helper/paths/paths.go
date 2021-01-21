// This is used for getting common Vagrant
// paths that are in use
package paths

import (
	"github.com/hashicorp/vagrant-plugin-sdk/helper/path"
	// "github.com/adrg/xdg" // TODO(spox): this is the lib we'll use for defaults
)

func VagrantHome() (path.Path, error) {
	return path.NewPath("~/.vagrant.d").Abs()
}
