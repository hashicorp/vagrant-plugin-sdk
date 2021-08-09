// This is used for getting common Vagrant
// paths that are in use
package paths

import (
	"io/ioutil"
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
		v = xdg.RuntimeDir
		if _, err := os.Stat(v); err != nil {
			if v, err = ioutil.TempDir("", "vagrant-tmp"); err != nil {
				return nil, err
			}
		}
		val = path.NewPath(v).Join("vagrant")
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
