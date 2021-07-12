package core

import (
	"io"

	"github.com/hashicorp/vagrant-plugin-sdk/component"
)

type Command interface {
	Execute([]string) (int64, error)
	CommandInfo() (*component.CommandInfo, error)

	io.Closer
}
