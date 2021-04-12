package core

import "github.com/DavidGamba/go-getoptions/option"

type CommandInfo struct {
	Name     []string
	Help     string
	Synopsis string
	Flags    []*option.Option
}

type Command interface {
	Execute(name string) (int64, error)
	Subcommands() ([]Command, error)
	CommandInfo() (*CommandInfo, error)
}
