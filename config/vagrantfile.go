package config

import (
	"github.com/hashicorp/hcl/v2"
)

type Vagrantfile struct {
	Name string `hcl:"name,label"`

	Body   hcl.Body `hcl:",body" json:"-"`
	Remain hcl.Body `hcl:",remain" json:"-"`
}
