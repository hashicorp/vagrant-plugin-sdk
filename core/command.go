package core

import "github.com/DavidGamba/go-getoptions"

type Command interface {
	Help() (string, error)
	Synopsis() (string, error)
	Flags() (*getoptions.GetOpt, error)
	Execute(name string) (int64, error)
}
