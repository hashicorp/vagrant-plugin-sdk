package core

import "github.com/hashicorp/vagrant-plugin-sdk/helper/path"

type BoxCollection interface {
	Add(path path.Path, name, version, metadataURL string, force bool, providers ...string) (box Box, err error)
	All() (boxes []Box, err error)
	Clean(name string) (err error)
	Find(name string, version string, providers ...string) (box Box, err error)
}
