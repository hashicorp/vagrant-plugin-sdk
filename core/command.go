package core

import "github.com/DavidGamba/go-getoptions/option"

type Command interface {
	Name() ([]string, error)
	Help() (string, error)
	Synopsis() (string, error)
	Flags() ([]*option.Option, error)
	Execute(name string) (int64, error)
	Subcommands() ([]Command, error)
}
