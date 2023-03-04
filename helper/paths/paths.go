// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

// This is used for getting common Vagrant
// paths that are in use
package paths

import (
	"fmt"
	"os"

	"github.com/adrg/xdg"
	"github.com/hashicorp/vagrant-plugin-sdk/helper/path"
)

func VagrantConfig() (path.Path, error) {
	var val path.Path
	v, ok := os.LookupEnv("VAGRANT_CONFIG")
	if ok {
		val = path.NewPath(v)
	} else {
		val = path.NewPath(xdg.ConfigHome).Join("vagrant")
	}

	return setupPath(val)
}

func VagrantCache() (path.Path, error) {
	var val path.Path
	v, ok := os.LookupEnv("VAGRANT_CACHE")
	if ok {
		val = path.NewPath(v)
	} else {
		val = path.NewPath(xdg.CacheHome).Join("vagrant")
	}

	return setupPath(val)
}

func VagrantCwd() (path.Path, error) {
	var val path.Path
	v, ok := os.LookupEnv("VAGRANT_CWD")
	if ok {
		if _, err := os.Stat(v); os.IsNotExist(err) {
			return nil, fmt.Errorf("VAGRANT_CWD set to path (%s) that does not exist", v)
		}
		val = path.NewPath(v)
	} else {
		cwd, err := os.Getwd()
		if err != nil {
			return nil, err
		}
		val = path.NewPath(cwd)
	}

	return val, nil
}

func VagrantData() (path.Path, error) {
	var val path.Path
	v, ok := os.LookupEnv("VAGRANT_DATA")
	if ok {
		val = path.NewPath(v)
	} else {
		val = path.NewPath(xdg.DataHome).Join("vagrant")
	}

	return setupPath(val)
}

func VagrantTmp() (path.Path, error) {
	var val path.Path
	v, ok := os.LookupEnv("VAGRANT_TMP")
	if ok {
		val = path.NewPath(v)
	} else {
		var err error
		if val, err = VagrantCache(); err != nil {
			return nil, err
		}
		val = val.Join("tmp")
	}

	return setupPath(val)
}

func NamedVagrantConfig(n string) (path.Path, error) {
	c, err := VagrantConfig()
	if err != nil {
		return nil, err
	}
	c = c.Join(n)
	return setupPath(c)
}

func NamedVagrantCache(n string) (path.Path, error) {
	c, err := VagrantCache()
	if err != nil {
		return nil, err
	}
	c = c.Join(n)
	return setupPath(c)
}

func NamedVagrantData(n string) (path.Path, error) {
	c, err := VagrantData()
	if err != nil {
		return nil, err
	}
	c = c.Join(n)
	return setupPath(c)
}

func NamedVagrantTmp(n string) (path.Path, error) {
	c, err := VagrantTmp()
	if err != nil {
		return nil, err
	}
	c = c.Join(n)
	return setupPath(c)
}

func setupPath(val path.Path) (path.Path, error) {
	p, err := val.Abs()
	if err != nil {
		return nil, err
	}

	if err := os.MkdirAll(p.String(), 0755); err != nil {
		return nil, err
	}
	return p, nil
}
