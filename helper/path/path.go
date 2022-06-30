// This package provides a simple way to represent a path
// value and perform common operations on it. All the functions
// expected with `path/filepath` are provided. The reason
// for using this over a simple string is that it allows
// Vagrant to modify the format of the path on demand. For
// example, if we wanted to reference all paths on Windows
// via UNC, or short names, or msys2, we simply have to
// enable it here.
package path

import (
	"errors"
	"os"
	"os/user"
	"path/filepath"
)

type Path interface {
	Abs() (Path, error)
	String() string
	Base() Path
	Clean() Path
	Dir() Path
	EvalSymLinks() (Path, error)
	Exists() bool
	Ext() string
	FromSlash() Path
	HasPrefix(string) bool
	IsAbs() bool
	IsRoot() (bool, error)
	Join(...string) Path
	Parent() Path
	SameFile(Path) (bool, error)
	Split() (Path, string)
	SplitList() []string
	ToSlash() string
	VolumeName() string
	Walk(filepath.WalkFunc) error
}

type path struct {
	path string
}

func NewPath(p string) Path {
	if len(p) > 0 && p[0:1] == "~" {
		u, err := user.Current()
		if err != nil {
			// TODO(spox): remove the panic and just retain when error
			// is encountered. panic for now to see if we run into it.
			panic("cannot determine current user")
		}
		p = filepath.Join(u.HomeDir, p[1:])
	}
	return &path{path: p}
}

func (p *path) String() string {
	return p.path
}

func (p *path) Abs() (newP Path, err error) {
	result, err := filepath.Abs(p.path)
	if err == nil {
		newP = &path{path: result}
	}
	return
}

func (p *path) Base() (newP Path) {
	return &path{path: filepath.Base(p.path)}
}

func (p *path) Clean() (newP Path) {
	return &path{path: filepath.Clean(p.path)}
}

func (p *path) Dir() (newP Path) {
	return &path{path: filepath.Dir(p.path)}
}

func (p *path) EvalSymLinks() (newP Path, err error) {
	result, err := filepath.EvalSymlinks(p.path)
	if err == nil {
		newP = &path{path: result}
	}
	return
}

func (p *path) Exists() bool {
	_, err := os.Stat(p.String())
	return !errors.Is(err, os.ErrNotExist)
}

func (p *path) Ext() string {
	return filepath.Ext(p.path)
}

func (p *path) FromSlash() (newP Path) {
	return &path{path: filepath.FromSlash(p.path)}
}

func (p *path) HasPrefix(prefix string) bool {
	return filepath.HasPrefix(p.path, prefix)
}

func (p *path) IsAbs() bool {
	return filepath.IsAbs(p.path)
}

func (p *path) IsRoot() (bool, error) {
	return p.SameFile(p.Parent())
}

func (p *path) Join(elm ...string) Path {
	return &path{path: filepath.Join(append([]string{p.path}, elm...)...)}
}

func (p *path) Parent() Path {
	parent, _ := p.Split()
	return parent.Dir()
}

func (p *path) SameFile(other Path) (bool, error) {
	myInfo, err := os.Stat(p.String())
	if err != nil {
		return false, err
	}

	otherInfo, err := os.Stat(other.String())
	if err != nil {
		return false, err
	}

	return os.SameFile(myInfo, otherInfo), nil
}

func (p *path) Split() (dir Path, file string) {
	d, file := filepath.Split(p.path)
	dir = &path{path: d}
	return
}

func (p *path) SplitList() []string {
	return filepath.SplitList(p.path)
}

func (p *path) ToSlash() string {
	return filepath.ToSlash(p.path)
}

func (p *path) VolumeName() string {
	return filepath.VolumeName(p.path)
}

func (p *path) Walk(walkFn filepath.WalkFunc) error {
	return filepath.Walk(p.path, walkFn)
}

func Match(pattern string, p Path) (bool, error) {
	return filepath.Match(pattern, p.String())
}

func Glob(pattern string) (matches []Path, err error) {
	m, err := filepath.Glob(pattern)
	if err != nil {
		return
	}
	for _, p := range m {
		matches = append(matches, &path{path: p})
	}
	return
}

func Rel(basepath, targetpath Path) (relP Path, err error) {
	result, err := filepath.Rel(basepath.String(), targetpath.String())
	if err == nil {
		relP = &path{path: result}
	}
	return
}

var _ Path = (*path)(nil)
