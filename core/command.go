package core

import "github.com/DavidGamba/go-getoptions/option"

type CommandInfo struct {
	Name        string
	Help        string
	Synopsis    string
	Flags       []*option.Option
	Subcommands []*CommandInfo
}

type Command interface {
	Execute([]string) (int64, error)
	CommandInfo([]string) (*CommandInfo, error)
}
