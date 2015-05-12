package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"strings"
	"os"
	"time"
	"io/ioutil"
	"path"
	"path/filepath"
)

func bindata_read(data []byte, name string) ([]byte, error) {
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

type bindata_file_info struct {
	name string
	size int64
	mode os.FileMode
	modTime time.Time
}

func (fi bindata_file_info) Name() string {
	return fi.name
}
func (fi bindata_file_info) Size() int64 {
	return fi.size
}
func (fi bindata_file_info) Mode() os.FileMode {
	return fi.mode
}
func (fi bindata_file_info) ModTime() time.Time {
	return fi.modTime
}
func (fi bindata_file_info) IsDir() bool {
	return false
}
func (fi bindata_file_info) Sys() interface{} {
	return nil
}

var _templates_mappings_md = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\x7c\x50\xc1\x6a\x85\x30\x10\xbc\xe7\x2b\x02\xf1\x24\x98\xde\x85\x9e\x5a\x0a\x05\xed\xa9\xf4\x1e\x71\x95\x40\x12\xc5\xe4\x16\xf3\xef\x5d\x93\xa8\xef\x81\xbe\x5c\x66\xd9\x9d\xcc\xcc\x2e\xa3\xde\xf3\x46\x74\xa0\x42\x20\xa4\xfc\xfe\xac\xb7\xc6\x8f\xd0\x10\xc2\x1b\x56\x7f\xb0\x58\x39\x99\x10\x4a\x42\x18\xa3\xbf\xa2\x53\x60\x89\xf7\x8b\x30\x23\xd0\xc2\xd1\xfa\x9d\xf2\xd4\xe5\x8d\xb4\x0e\x55\x18\xdb\x44\x0b\x97\x55\xc8\xc9\x1e\x36\x36\x0e\xbe\x24\xa8\xfe\xe0\x7b\x2f\x07\x9c\xf1\x56\xcc\xb3\x34\xa3\x4d\x12\x51\x63\xd8\x35\xda\xa9\x07\x45\x57\x9a\xe3\x60\x15\x3d\x11\xa3\x16\xe2\xc7\xa4\x35\x18\x47\xaa\xf8\xd6\x6a\x7f\xeb\x0d\x56\x67\x2c\x1d\x63\x3d\x06\x40\x6b\x9d\x52\xa6\xdd\x78\xf4\xcf\x61\xd0\xec\x66\x7e\x1c\xeb\x8a\x72\xf5\xf9\xb9\x97\x57\x88\x27\x01\xd3\xbf\xc0\xff\x00\x00\x00\xff\xff\x53\xd1\xb3\x27\xb6\x01\x00\x00")

func templates_mappings_md_bytes() ([]byte, error) {
	return bindata_read(
		_templates_mappings_md,
		"templates/mappings.md",
	)
}

func templates_mappings_md() (*asset, error) {
	bytes, err := templates_mappings_md_bytes()
	if err != nil {
		return nil, err
	}

	info := bindata_file_info{name: "templates/mappings.md", size: 438, mode: os.FileMode(420), modTime: time.Unix(1431396413, 0)}
	a := &asset{bytes: bytes, info:  info}
	return a, nil
}

var _templates_model_md = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\x64\x8e\x3d\x8b\x83\x40\x10\x86\xfb\xf9\x15\x03\x6b\x71\x27\xe8\xf5\xc2\x75\x12\x08\x48\x2a\x49\x13\x52\xac\x71\x0c\x0b\xc6\x88\xbb\x4d\x58\xe7\xbf\x67\xd6\x8f\x28\x64\x9b\x77\x67\x60\x9e\xe7\x55\xe8\x7d\x5a\xe8\x8a\x5a\x66\x80\xf8\x98\x67\x61\x71\xd2\x0f\x62\xfe\x93\xdf\x99\x06\x6b\x9e\x1d\x73\x0c\xa0\x14\x96\xba\x6a\xc9\x82\xf7\x83\xee\xee\x84\x91\xc3\xec\x1f\xd3\x79\x9b\x16\xc6\x3a\xe6\x04\x2f\xde\x47\x6e\x61\x5c\x7f\xd4\x6e\xfa\x95\x4b\xea\x6a\x51\x6d\x84\x2f\x80\x78\x42\xab\xcf\x11\x40\x48\x1c\xb1\x7c\xf5\x21\x72\xb2\xb7\xc1\xf4\x4e\x6a\x41\x12\xde\x98\xec\x63\x7e\x9b\xa0\x11\x81\xb0\x0e\x86\xda\x7a\x55\x08\xbd\x59\xe8\x02\x9c\xa6\x00\xdf\xa6\x9d\x63\x2a\xbb\x96\x9e\xf3\x1d\x00\x00\xff\xff\x1f\x66\x97\xf6\x35\x01\x00\x00")

func templates_model_md_bytes() ([]byte, error) {
	return bindata_read(
		_templates_model_md,
		"templates/model.md",
	)
}

func templates_model_md() (*asset, error) {
	bytes, err := templates_model_md_bytes()
	if err != nil {
		return nil, err
	}

	info := bindata_file_info{name: "templates/model.md", size: 309, mode: os.FileMode(420), modTime: time.Unix(1431394174, 0)}
	a := &asset{bytes: bytes, info:  info}
	return a, nil
}

var _templates_schema_md = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\x7c\x90\x4d\x4b\xc4\x30\x10\x86\xef\xf3\x2b\x06\xb2\x07\x5d\x68\xbd\x2f\x78\x5b\x04\xa1\x88\xe0\xe2\x45\x3c\x64\x77\xa7\x35\x10\xa3\xb4\xf1\x20\xe9\xfc\x77\x67\xfa\x95\x9c\xec\xe5\xfd\x28\xc9\xfb\x10\x83\x29\xd5\x8d\x3d\x93\x67\x06\xd8\x3f\x1e\x0f\x5a\x3c\xd9\x4f\x62\xbe\x13\xf7\x4a\xfd\xe0\xbe\x02\xf3\x1e\xc0\x18\x3c\xd9\xb3\xa7\x01\x52\xea\x6d\xe8\x08\x77\x11\x0f\xf7\x58\xcf\x6d\xdd\xb8\x21\x32\x57\xf8\x96\xd2\x2e\x2e\x77\xbc\xdf\x98\x22\xdd\xca\x49\x0a\x57\x99\xfa\xef\x06\x19\x52\xac\xed\x14\x80\x2a\x8e\x78\xfa\xfd\x56\x69\x28\x74\xf1\x43\xcc\x73\x4f\x17\xa7\x78\xe2\x5f\x2e\xd6\xeb\xcf\x23\xb5\xf6\xc7\x47\xa8\xf4\x1b\xab\x52\x0a\x53\xf8\xad\xcb\x4c\xad\x32\xc9\xfc\x83\x23\x7f\x5d\xa9\x04\xa8\x5d\x80\x64\x66\x4a\xca\x93\xd3\x8c\x95\xf3\x46\x97\xab\x09\x32\xc7\x85\x75\x7a\x8e\xf5\x59\x66\xfd\x0b\x00\x00\xff\xff\x3e\xeb\x12\xcc\x97\x01\x00\x00")

func templates_schema_md_bytes() ([]byte, error) {
	return bindata_read(
		_templates_schema_md,
		"templates/schema.md",
	)
}

func templates_schema_md() (*asset, error) {
	bytes, err := templates_schema_md_bytes()
	if err != nil {
		return nil, err
	}

	info := bindata_file_info{name: "templates/schema.md", size: 407, mode: os.FileMode(420), modTime: time.Unix(1431394204, 0)}
	a := &asset{bytes: bytes, info:  info}
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
	if (err != nil) {
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
	"templates/mappings.md": templates_mappings_md,
	"templates/model.md": templates_model_md,
	"templates/schema.md": templates_schema_md,
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
	for name := range node.Children {
		rv = append(rv, name)
	}
	return rv, nil
}

type _bintree_t struct {
	Func func() (*asset, error)
	Children map[string]*_bintree_t
}
var _bintree = &_bintree_t{nil, map[string]*_bintree_t{
	"templates": &_bintree_t{nil, map[string]*_bintree_t{
		"mappings.md": &_bintree_t{templates_mappings_md, map[string]*_bintree_t{
		}},
		"model.md": &_bintree_t{templates_model_md, map[string]*_bintree_t{
		}},
		"schema.md": &_bintree_t{templates_schema_md, map[string]*_bintree_t{
		}},
	}},
}}

// Restore an asset under the given directory
func RestoreAsset(dir, name string) error {
        data, err := Asset(name)
        if err != nil {
                return err
        }
        info, err := AssetInfo(name)
        if err != nil {
                return err
        }
        err = os.MkdirAll(_filePath(dir, path.Dir(name)), os.FileMode(0755))
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

// Restore assets under the given directory recursively
func RestoreAssets(dir, name string) error {
        children, err := AssetDir(name)
        if err != nil { // File
                return RestoreAsset(dir, name)
        } else { // Dir
                for _, child := range children {
                        err = RestoreAssets(dir, path.Join(name, child))
                        if err != nil {
                                return err
                        }
                }
        }
        return nil
}

func _filePath(dir, name string) string {
        cannonicalName := strings.Replace(name, "\\", "/", -1)
        return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}
