// Code generated for package localizer by go-bindata DO NOT EDIT. (@generated)
// sources:
// localizer/locales/en.json
package localizer

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func bindataRead(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}
	if clErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes []byte
	info  os.FileInfo
}

type bindataFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

// Name return file name
func (fi bindataFileInfo) Name() string {
	return fi.name
}

// Size return file size
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}

// Mode return file mode
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}

// Mode return file modify time
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}

// IsDir return file whether a directory
func (fi bindataFileInfo) IsDir() bool {
	return fi.mode&os.ModeDir != 0
}

// Sys return file is sys mode
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _localizerLocalesEnJson = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x4c\x90\xb1\x6a\xc3\x40\x10\x44\x7b\x7f\xc5\xa0\xc6\x4d\xd0\x07\xa4\x4b\x13\x70\x11\x48\xe1\x04\x52\x89\x95\x6e\xa5\x5b\x38\xdf\x39\x7b\x2b\xd9\xc1\xf8\xdf\xc3\x5a\x38\x4e\x2b\xbd\x99\x37\x7b\x97\x0d\xd0\xf4\xe5\xdc\x51\x08\xcd\x33\x9a\x97\x10\x24\x4f\xe8\xcb\xb9\x79\xf2\x5f\xac\x5a\xb4\xcb\x25\x4b\x36\x56\x1a\x4c\x16\xee\x66\x71\xf4\x93\x26\xa5\x6c\x90\x0a\x32\xe3\xc3\xd1\x3c\x69\x05\x37\x74\xa4\x81\x71\x12\x8b\xb0\xc8\xf8\xd8\x41\x32\x08\x27\xfa\x81\x45\x32\x28\x7f\xcf\xa2\x5c\x41\xd8\xef\xbf\x5a\xbc\x95\x6a\xf0\xfa\x92\xab\xa3\xf7\xf2\xff\xf0\xca\x22\xd2\xc2\x18\x4a\x1e\x65\x9a\x95\x3c\x81\x7a\x12\x1b\x22\x57\xb7\x07\xa9\xd4\x27\x86\x45\xa9\xf7\xe4\x81\xb3\xb5\x78\x4f\x4c\x95\x11\xca\xda\x5a\x14\x3a\x3f\x4c\xb7\xad\xbe\x65\x3d\x3c\xf0\x51\x79\x20\xe3\xd0\x8d\x89\x26\x3f\x78\x1f\x19\x47\x2d\x8b\x04\x0e\xf0\x8f\xd8\x5e\x2e\xed\x6b\xa2\xe9\x7a\xdd\xfa\x33\x3c\x32\x2d\x76\x19\xe3\x6c\xb3\x32\x94\x6f\xde\x8a\x32\xfe\xc9\xc4\x7d\x29\x21\x17\x43\xcf\xa0\x85\x24\xf9\xea\xb6\xd9\x5c\x37\xbf\x01\x00\x00\xff\xff\x86\xe1\xa1\x8a\x94\x01\x00\x00")

func localizerLocalesEnJsonBytes() ([]byte, error) {
	return bindataRead(
		_localizerLocalesEnJson,
		"localizer/locales/en.json",
	)
}

func localizerLocalesEnJson() (*asset, error) {
	bytes, err := localizerLocalesEnJsonBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "localizer/locales/en.json", size: 404, mode: os.FileMode(420), modTime: time.Unix(1653669692, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if err != nil {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() (*asset, error){
	"localizer/locales/en.json": localizerLocalesEnJson,
}

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for childName := range node.Children {
		rv = append(rv, childName)
	}
	return rv, nil
}

type bintree struct {
	Func     func() (*asset, error)
	Children map[string]*bintree
}

var _bintree = &bintree{nil, map[string]*bintree{
	"localizer": &bintree{nil, map[string]*bintree{
		"locales": &bintree{nil, map[string]*bintree{
			"en.json": &bintree{localizerLocalesEnJson, map[string]*bintree{}},
		}},
	}},
}}

// RestoreAsset restores an asset under the given directory
func RestoreAsset(dir, name string) error {
	data, err := Asset(name)
	if err != nil {
		return err
	}
	info, err := AssetInfo(name)
	if err != nil {
		return err
	}
	err = os.MkdirAll(_filePath(dir, filepath.Dir(name)), os.FileMode(0755))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
	if err != nil {
		return err
	}
	err = os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
	if err != nil {
		return err
	}
	return nil
}

// RestoreAssets restores an asset under the given directory recursively
func RestoreAssets(dir, name string) error {
	children, err := AssetDir(name)
	// File
	if err != nil {
		return RestoreAsset(dir, name)
	}
	// Dir
	for _, child := range children {
		err = RestoreAssets(dir, filepath.Join(name, child))
		if err != nil {
			return err
		}
	}
	return nil
}

func _filePath(dir, name string) string {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}
