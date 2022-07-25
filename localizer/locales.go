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

var _localizerLocalesEnJson = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x4c\x52\x3f\x6f\x9d\x30\x10\xdf\xdf\xa7\xf8\x89\x25\x4b\x04\xca\x9a\xad\x1d\x2a\x65\x68\xd5\xe1\x35\x52\xa5\x48\xe4\xc0\x07\x58\x32\x3e\x6a\x9f\x81\x08\xf1\xdd\x2b\xc3\xeb\x6b\x56\xdf\xef\xef\x9d\xb7\x0b\x50\x90\x31\xd6\xf7\xf5\xfc\x54\x37\xb2\x16\xcf\x28\xde\xfc\x97\xe3\x09\xaf\x4f\x78\xa5\x3e\x90\x57\x34\xb2\x3e\xe2\x7d\xdb\xca\xaf\xb2\xfe\xa0\x91\xf7\xfd\xbd\xbc\x0f\x17\xeb\x1c\x28\xa9\x8c\xa4\xb6\x25\xe7\x3e\x90\xa6\x3e\x90\x61\xe8\x60\x63\x26\x97\xf8\xe9\x98\x22\xa3\x15\x1f\xad\xe1\x70\x43\x64\x1b\x1d\x18\x51\x52\x68\xf9\x44\x0e\xaa\x53\x7c\xae\xaa\x65\x59\xca\xf9\xb4\x48\x53\xd9\xca\x58\x19\x69\x63\xd5\xc8\xca\xb1\xea\x24\x8c\xa4\x6f\xbe\x78\xcc\x1d\x1a\x59\x6b\x32\x26\xa7\xbf\x65\xcf\x5d\x8e\x91\xe1\x29\x70\x4b\xca\xa6\xee\x1c\xf5\x19\x72\x1d\x18\x53\x90\xd9\x1a\x36\xc8\x8f\x78\xd8\xb6\xf2\x9b\xa3\x7e\xdf\x1f\x60\x23\xfe\x73\x4a\xbc\x78\x74\x49\x53\x60\x04\x3e\x2a\x44\x48\x77\xaf\x6e\x6f\xed\xbd\x28\x1a\x06\xcd\x64\x1d\x35\x8e\xcb\xd3\x9c\x43\x90\x50\x7b\xf1\xd6\x2b\x07\x6a\xd5\xce\x5c\x27\x9b\x43\xdc\x15\x22\x48\x95\xc7\x49\x8f\x5d\x08\x0e\x68\x47\x2d\x63\xb1\x3a\x1c\xdb\xf9\xf5\x02\xeb\x41\x58\xe8\x03\x3a\x90\x22\xf0\x9f\x64\x03\x47\x10\xae\xd7\xdf\x25\xbe\x4b\x54\x64\x79\xf1\x31\x43\xff\x89\x7f\x06\x9f\x58\x0c\x34\x1f\x57\xe8\x6c\x9f\x02\x65\x06\xe2\x62\xb5\x1d\x38\x66\x77\x63\x63\xce\x7f\x1e\xee\xc6\x1c\xd9\xeb\xfd\x80\x46\x4e\x55\x09\x08\xc9\x7f\xfa\x03\x3a\x1c\x59\x8a\xcb\x7e\xf9\x1b\x00\x00\xff\xff\x5f\xc7\xa1\x7d\x58\x02\x00\x00")

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

	info := bindataFileInfo{name: "localizer/locales/en.json", size: 600, mode: os.FileMode(420), modTime: time.Unix(1658769759, 0)}
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
