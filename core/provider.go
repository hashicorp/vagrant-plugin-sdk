package core

type Provider interface {
	Usable() bool
	Installed() bool
	Init(Machine) interface{}
	Name() string
}
