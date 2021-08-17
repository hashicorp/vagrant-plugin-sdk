package core

type Provider interface {
	Init(Machine) interface{}
	Installed() bool
	Name() string
	Usable() bool
}
