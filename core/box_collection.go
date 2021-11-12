package core

type BoxCollection interface {
	Add(path, name, version string, metadataURL string, providers []string) (box Box, err error)
	All() (list Box, err error)
	Clean(name string) (err error)
	Find(name string, providers []string, version string) (box Box, err error)
}
