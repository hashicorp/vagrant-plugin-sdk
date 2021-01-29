package core

type Command interface {
	Help() (string, error)
	Synopsis() (string, error)
	Flags() (string, error)
	Execute(name string) (int64, error)
}
