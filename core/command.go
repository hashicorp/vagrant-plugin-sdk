package core

import "github.com/DavidGamba/go-getoptions/option"

type Command interface {
	Help() (string, error)
	Synopsis() (string, error)
	Flags() ([]*option.Option, error)
	Execute(name string) (int64, error)
}
