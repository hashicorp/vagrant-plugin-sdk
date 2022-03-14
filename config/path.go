// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package config

import (
	"os"

	"github.com/hashicorp/vagrant-plugin-sdk/helper/path"
	"github.com/hashicorp/vagrant-plugin-sdk/helper/paths"
)

// RubyFilename is the default Ruby based Vagrantfile name
const RubyFilename = "Vagrantfile"

// HCLFilename is the default HCL based Vagrantfile name
const HCLFilename = "Vagrantfile.hcl"

// Returns the valid file names for the Vagrantfile. If
// the VAGRANT_VAGRANTFILE environment variable is set,
// the slice will only contain that value.
func GetVagrantfileName() []string {
	if f := os.Getenv("VAGRANT_VAGRANTFILE"); f != "" {
		return []string{f}
	}
	return []string{RubyFilename, HCLFilename}
}

// FindPath looks for our configuration file starting at "start" and
// traversing parent directories until it is found. If it is found, the
// path is returned. If it is not found, nil is returned. Error will be
// non-nil only if an error occurred.
//
// If start is nil, start will be the current working directory. If
// filenames are nil, it will default to the value(s) from GetVagrantfileName().
func FindPath(
	start path.Path, // starting directory. defaults to vagrant CWD if nil
	filenames []string, // list of valid Vagrantfile names
) (p path.Path, err error) {
	if start == nil {
		start, err = paths.VagrantCwd()
		if err != nil {
			return nil, err
		}
	}

	if filenames == nil {
		filenames = GetVagrantfileName()
	}

	for _, f := range filenames {
		p = start
		for {
			p = p.Join(f)
			if p.Exists() {
				return
			}
			root, err := p.Parent().IsRoot()
			if err != nil {
				return nil, err
			}
			if root {
				break
			}
			p = p.Parent().Parent()
		}
	}

	return nil, nil
}

// Detect existing path within directory
func ExistingPath(
	dir path.Path, // directory to check withing
	filenames []string, // list of valid file names
) (p path.Path, err error) {
	for _, f := range filenames {
		p = dir.Join(f)
		if _, err = os.Stat(p.String()); err == nil || !os.IsNotExist(err) {
			return
		}
	}

	return
}
