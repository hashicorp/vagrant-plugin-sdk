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

var _localizerLocalesEnJson = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x5c\x52\xbd\x6e\xdb\x30\x10\xde\xfd\x14\x07\x2f\x5e\x02\x19\x59\xb3\xb5\x43\x81\x0c\x29\x32\xb8\x01\x0a\x04\x50\x4e\xe4\x49\x24\x2a\xf1\xd4\xe3\x51\xb2\x61\xf8\xdd\x0b\x92\xaa\x1b\x74\x93\xc8\xef\xbe\xbf\xe3\x75\x07\xb0\x47\x6b\x7d\x18\xda\xe5\xb1\xed\xf8\xbc\x7f\x82\xfd\x7b\xf8\x52\x8e\xe0\xed\x11\xde\x70\x10\x0c\x0a\x1d\x9f\x1f\xe0\xe3\x7a\x6d\xbe\xf2\xf9\x3b\x4e\x74\xbb\x7d\x34\xf7\xcb\xd5\x8f\x23\x60\x52\x9e\x50\xbd\xc1\x71\xbc\x40\x9a\x07\x41\x4b\xa0\xce\xc7\x3c\xdc\xc0\xeb\x48\x18\x09\x0c\x87\xe8\x2d\xc9\x86\xc8\x32\xea\x08\x22\x27\x31\x54\x91\x4e\x75\x8e\x4f\xc7\xe3\xba\xae\xcd\x52\x25\xd2\xdc\x18\x9e\x8e\x96\x4d\x3c\x76\x7c\xa6\x78\xec\x59\x26\xd4\xf7\xb0\x7f\xc8\x19\x3a\x3e\xb7\x68\x6d\x76\xbf\x79\xcf\x59\xca\x95\xa5\x59\xc8\xa0\x92\x6d\xfb\x11\x87\x0c\x39\x39\x82\x59\x78\xf1\x96\x2c\xe4\x43\x38\x5c\xaf\xcd\xb7\x11\x87\xdb\xed\x00\x3e\xc2\xbf\x99\x06\x9e\x03\xf4\x49\x93\x10\x08\x95\x08\x11\xb8\xbf\x47\xf7\x5b\xfa\xc0\x0a\x1d\x01\x2e\xe8\x47\xec\x46\x6a\xaa\x38\x89\xb0\xb4\x81\x83\x0f\x4a\x82\x46\xfd\x42\x6d\xf2\xd9\xc4\x9d\x21\x02\xaa\xd2\x34\x6b\xe9\x82\xa1\x40\x7b\x34\x04\xab\x57\x57\xda\xf9\xf1\x0c\x3e\x00\xc2\x8a\x17\x50\x87\x0a\x42\xbf\x93\x17\x8a\x80\x70\x3a\xfd\x6c\xe0\x85\xa3\x42\xa6\xe7\x10\x33\xf4\x2f\xf9\x67\x70\xc5\x82\xc3\xa5\x6c\xa1\xf7\x43\x12\xcc\x13\x10\x57\xaf\xc6\x51\xcc\xea\xd6\xc7\xec\xbf\x2e\x6e\x9b\x9c\x28\xe8\x7d\x81\x96\x2b\x2b\x0b\x48\x0a\x9f\xde\x80\xba\xe2\xa5\x06\xdf\xea\xcd\xd9\xb5\x4d\x85\xf2\xbf\xe6\xa5\x94\xfe\xba\xfd\xe4\xe2\x0b\xed\x8a\x55\x96\xa2\x92\xcd\x86\x3a\x34\xbf\x4a\x09\x13\x1a\xe7\x03\x95\xb1\x97\xfa\xbd\xad\x4b\x68\x66\xa9\xf5\x65\x0a\x9f\x4b\x0d\x07\x85\xaa\x0b\x1c\x6a\x9a\x78\x89\x4a\xd3\x7e\x77\xdb\xfd\x09\x00\x00\xff\xff\xcd\xd3\x02\x40\xf8\x02\x00\x00")

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

	info := bindataFileInfo{name: "localizer/locales/en.json", size: 760, mode: os.FileMode(420), modTime: time.Unix(1662500611, 0)}
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
