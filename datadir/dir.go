package datadir

import (
	"os"

	"github.com/hashicorp/vagrant-plugin-sdk/helper/path"
)

// TODO(spox): Windows should use UNC based paths to side step MAX_PATH limitations

// TODO(mitchellh): tests! like any tests

// Dir is the interface implemented so that consumers can store data
// locally in a consistent way.
type Dir interface {
	// CacheDir returns the path to a folder that can be used for
	// cache data. This directory may not be empty if a previous run
	// stored data, but it may also be emptied at any time between runs.
	CacheDir() path.Path

	// DataDir returns the path to a folder that can be used for data
	// that is persisted between runs.
	DataDir() path.Path

	// RootDir returns the top level folder.
	RootDir() path.Path

	// TempDir returns the path to a folder that can be used for temporary
	// data storage. Directory is routinely cleaned.
	TempDir() path.Path
}

// basicDir implements Dir in the simplest possible way.
type basicDir struct {
	cacheDir path.Path
	dataDir  path.Path
	tempDir  path.Path
	rootDir  path.Path
}

// CacheDir impl Dir
func (d *basicDir) CacheDir() path.Path { return d.cacheDir }

// DataDir impl Dir
func (d *basicDir) DataDir() path.Path { return d.dataDir }

// RootDir impl Dir
func (d *basicDir) RootDir() path.Path { return d.rootDir }

// TempDir impl Dir
func (d *basicDir) TempDir() path.Path { return d.tempDir }

// newRootDir creates a basicDir for the root directory which puts
// data at <path>/cache, etc.
func newRootDir(rPath string) (Dir, error) {
	root := path.NewPath(rPath)
	if err := os.MkdirAll(root.String(), 0755); err != nil {
		return nil, err
	}

	cacheDir := root.Join("cache")
	dataDir := root.Join("data")
	tmpDir := root.Join("tmp")
	if err := os.MkdirAll(cacheDir.String(), 0755); err != nil {
		return nil, err
	}
	if err := os.MkdirAll(dataDir.String(), 0755); err != nil {
		return nil, err
	}
	if err := os.MkdirAll(tmpDir.String(), 0755); err != nil {
		return nil, err
	}

	return &basicDir{
		cacheDir: cacheDir,
		dataDir:  dataDir,
		rootDir:  root,
		tempDir:  tmpDir,
	}, nil
}

// NewBasicDir creates a Dir implementation with a manually specified
// set of directories.
func NewBasicDir(rootDir, cacheDir, dataDir, tempDir string) Dir {
	return &basicDir{
		cacheDir: path.NewPath(cacheDir),
		dataDir:  path.NewPath(dataDir),
		tempDir:  path.NewPath(tempDir),
		rootDir:  path.NewPath(rootDir),
	}
}

// NewScopedDir creates a ScopedDir for the given parent at the relative
// child path of path. The caller should take care that multiple scoped
// dirs with overlapping paths are not created, since they could still
// collide.
func NewScopedDir(parent Dir, rPath string) (Dir, error) {
	relPath := path.NewPath(rPath)
	cacheDir := parent.CacheDir().Join(relPath.String())
	dataDir := parent.DataDir().Join(relPath.String())
	if err := os.MkdirAll(cacheDir.String(), 0755); err != nil {
		return nil, err
	}
	if err := os.MkdirAll(dataDir.String(), 0755); err != nil {
		return nil, err
	}

	return &basicDir{
		cacheDir: cacheDir,
		dataDir:  dataDir,
		rootDir:  parent.RootDir(),
		tempDir:  parent.TempDir(),
	}, nil
}

var _ Dir = (*basicDir)(nil)
