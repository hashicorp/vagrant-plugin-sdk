package core

type Named interface {
	SetPluginName(string) error
	PluginName() (name string, err error)
}
