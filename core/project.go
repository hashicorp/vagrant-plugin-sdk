// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package core

import (
	"io"

	//	"github.com/hashicorp/vagrant-plugin-sdk/config"
	"github.com/hashicorp/vagrant-plugin-sdk/datadir"
	"github.com/hashicorp/vagrant-plugin-sdk/helper/path"
	"github.com/hashicorp/vagrant-plugin-sdk/terminal"
)

type DefaultProviderOptions struct {
	CheckUsable  bool
	Exclude      []string
	ForceDefault bool
	MachineName  string
}

func (d *DefaultProviderOptions) IsExcluded(provider string) bool {
	for _, e := range d.Exclude {
		if e == provider {
			return true
		}
	}
	return false
}

type Project interface {
	ActiveTargets() (targets []Target, err error)
	Boxes() (boxes BoxCollection, err error)
	// accessors
	CWD() (path path.Path, err error)
	Config() (v Vagrantfile, err error)
	DataDir() (dir *datadir.Project, err error)
	DefaultPrivateKey() (path path.Path, err error)
	DefaultProvider(opts *DefaultProviderOptions) (name string, err error)
	Home() (path path.Path, err error)
	Host() (h Host, err error)
	LocalData() (path path.Path, err error)
	PrimaryTargetName() (name string, err error)
	ResourceId() (string, error)
	RootPath() (path path.Path, err error)

	// Target loads a target within this project with the given name. The
	// provider parameter is optional and is used to specify which provider
	// should be used to load the machine. This second parameter can be left
	// blank when fetching an existing target, but can be specified during
	// machine up to indicate a user flag that's been provided.
	Target(name string, provider string) (t Target, err error)
	TargetIds() (ids []string, err error)
	TargetIndex() (index TargetIndex, err error)
	TargetNames() (names []string, err error)
	Tmp() (path path.Path, err error)
	UI() (ui terminal.UI, err error)
	Vagrantfile() (Vagrantfile, error)
	VagrantfileName() (name string, err error)
	VagrantfilePath() (p path.Path, err error)

	// Not entirely sure if these are needed yet
	// Lock(name string) (err error)
	// Unlock(name string) (err error)
	// Unload() (err error)

	io.Closer
}
