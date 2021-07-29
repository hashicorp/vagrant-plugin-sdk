package datadir

import (
	"os"

	"github.com/hashicorp/vagrant-plugin-sdk/helper/path"
	"github.com/hashicorp/vagrant-plugin-sdk/helper/paths"
)

// Dir is the interface implemented so that consumers can store data
// locally in a consistent way.
type Dir interface {
	// CacheDir returns the path to a folder that can be used for
	// cache data. This directory may not be empty if a previous run
	// stored data, but it may also be emptied at any time between runs.
	CacheDir() path.Path

	// DataDir returns the path to a folder that can be used for data
	// storage
	DataDir() path.Path

	// ConfigDir is the path to a folder that can be used for storing
	// configuration type information.
	ConfigDir() path.Path

	// TempDir returns the path to a folder that can be used for temporary
	// data storage. Directory is routinely cleaned.
	TempDir() path.Path
}

// basicDir implements Dir in the simplest possible way.
type basicDir struct {
	cacheDir  path.Path
	configDir path.Path
	dataDir   path.Path
	tempDir   path.Path
}

// CacheDir impl Dir
func (d *basicDir) CacheDir() path.Path { return d.cacheDir }

// DataDir impl Dir
func (d *basicDir) DataDir() path.Path { return d.dataDir }

// ConfigDir impl Dir
func (d *basicDir) ConfigDir() path.Path { return d.configDir }

// TempDir impl Dir
func (d *basicDir) TempDir() path.Path { return d.tempDir }

func newDir(ident string) (d Dir, err error) {
	var cfg, csh, dat, tmp path.Path
	if cfg, err = paths.NamedVagrantConfig(ident); err != nil {
		return
	}
	if csh, err = paths.NamedVagrantCache(ident); err != nil {
		return
	}
	if dat, err = paths.NamedVagrantData(ident); err != nil {
		return
	}
	if tmp, err = paths.NamedVagrantTmp(ident); err != nil {
		return
	}

	return &basicDir{
		cacheDir:  csh,
		configDir: cfg,
		dataDir:   dat,
		tempDir:   tmp,
	}, nil
}

// NewBasicDir creates a Dir implementation with a manually specified
// set of directories.
func NewBasicDir(configDir, cacheDir, dataDir, tempDir string) Dir {
	return &basicDir{
		cacheDir:  path.NewPath(cacheDir),
		dataDir:   path.NewPath(dataDir),
		tempDir:   path.NewPath(tempDir),
		configDir: path.NewPath(configDir),
	}
}

// NewScopedDir creates a ScopedDir for the given parent at the relative
// child path of path. The caller should take care that multiple scoped
// dirs with overlapping paths are not created, since they could still
// collide.
func NewScopedDir(parent Dir, ident string) (Dir, error) {
	csh := parent.CacheDir().Join(ident)
	cfg := parent.ConfigDir().Join(ident)
	dat := parent.DataDir().Join(ident)
	tmp := parent.TempDir().Join(ident)

	for _, p := range []path.Path{csh, cfg, dat, tmp} {
		if err := os.MkdirAll(p.String(), 0755); err != nil {
			return nil, err
		}
	}

	return &basicDir{
		cacheDir:  csh,
		configDir: cfg,
		dataDir:   dat,
		tempDir:   tmp,
	}, nil
}

var _ Dir = (*basicDir)(nil)
