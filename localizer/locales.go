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

var _localizerLocalesEnJson = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x4c\x8e\xb1\x4e\x2b\x41\x0c\x45\xfb\xfd\x8a\xab\xa9\x9f\xf2\x01\xaf\xa3\xa4\x40\xa2\x08\x48\x54\x2b\xef\x8e\xb3\x63\x29\xb1\xc1\xf6\x66\x83\x10\xff\x8e\x26\x51\x04\xb5\xcf\x39\xbe\x5f\x03\x50\x26\xbb\x8c\x54\x6b\xf9\x8f\xf2\x50\xab\xe8\x82\xc9\x2e\xe5\x5f\x3f\xb1\xbb\xf9\xa8\xa6\xa2\xc9\x4e\x73\xca\x99\xc7\x55\x3a\xfa\x4a\x8b\x93\x26\x24\x40\x99\x7c\x7a\xcf\x6e\xa6\xe1\x8a\x1e\x68\x66\x6c\x92\x0d\xd9\x18\x2f\x8f\x10\x05\x61\xa3\x4f\x64\xa3\x84\xf3\xc7\x2a\xce\x01\xc2\x7e\xff\xb6\xc3\x93\x45\xa2\xe7\x4d\xa3\xa3\xf7\xf8\x5f\xf8\xc6\xa2\xd1\x99\x31\x9b\x1e\x64\x59\x9d\xba\x81\xd8\x24\xe7\xc6\xd1\xbf\x57\x09\x9a\x8e\x8c\x6c\x12\x77\xf3\xc4\x9a\x3b\x3c\x1f\x99\x82\x51\xed\x56\x35\x87\xaf\xbf\x9f\xae\x5b\xfb\x96\x32\x7c\x0f\x3f\x01\x00\x00\xff\xff\xd4\x48\xc4\x32\x17\x01\x00\x00")

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

	info := bindataFileInfo{name: "localizer/locales/en.json", size: 279, mode: os.FileMode(420), modTime: time.Unix(1652718867, 0)}
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
