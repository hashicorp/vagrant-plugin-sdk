package action

import (
	"github.com/hashicorp/vagrant-plugin-sdk/multistep"
	"github.com/hashicorp/vagrant-plugin-sdk/terminal"
)

// TODO(spox): based on the packer implementation the runner has some
// customized implementations for different behaviors (debugging) as
// well as wrapping each steps in a custom step for providing extra
// UI output when aborted, etc. not sure if we want to care about
// extras like that, so leave the UI here for now, but remove if
// decided it's not required.
func NewRunner(steps []multistep.Step, ui terminal.UI) multistep.Runner {
	return &multistep.BasicRunner{Steps: steps}
}
