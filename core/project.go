package core

import (
	"io"

	"github.com/hashicorp/vagrant-plugin-sdk/datadir"
	"github.com/hashicorp/vagrant-plugin-sdk/helper/path"
	"github.com/hashicorp/vagrant-plugin-sdk/proto/vagrant_plugin_sdk"
	"github.com/hashicorp/vagrant-plugin-sdk/terminal"
)

type Project interface {
	ActiveTargets() (targets []Target, err error)
	Boxes() (boxes BoxCollection, err error)
	Config() (v *vagrant_plugin_sdk.Vagrantfile_Vagrantfile, err error)
	CWD() (path string, err error)
	DataDir() (dir *datadir.Project, err error)
	DefaultPrivateKey() (path string, err error)
	DefaultProvider() (name string, err error)
	Home() (path string, err error)
	Host() (h Host, err error)
	LocalData() (path string, err error)
	PrimaryTargetName() (name string, err error)
	ResourceId() (string, error)
	RootPath() (path string, err error)
	Target(name string) (t Target, err error)
	TargetIds() (ids []string, err error)
	TargetIndex() (index TargetIndex, err error)
	TargetNames() (names []string, err error)
	Tmp() (path string, err error)
	UI() (ui terminal.UI, err error)
	VagrantfileName() (name string, err error)
	VagrantfilePath() (p path.Path, err error)

	// Not entirely sure if these are needed yet
	// Lock(name string) (err error)
	// Unlock(name string) (err error)
	// Unload() (err error)

	io.Closer
}
