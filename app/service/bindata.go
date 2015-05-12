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

var _templates_definition_md = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\x74\x50\x4d\x4b\x03\x31\x10\xbd\xcf\xaf\x18\xd8\x1e\x6c\x61\xe3\xbd\xe0\xad\x08\x42\xf1\x20\xc5\x4b\xf1\x90\xb6\x13\x09\xac\xb1\x24\xb9\x48\x3a\xff\xdd\xc9\xc7\xba\xbb\xa8\x7b\x99\x7d\xf3\xf1\xde\xcb\xeb\x30\x25\xb5\xd7\x27\x1a\x98\x01\x36\x4f\xbb\x6d\x6e\x3c\xeb\x0f\x62\xbe\x97\xbf\x57\xf2\xc1\x7e\x3a\xe6\x0d\x40\xd7\xe1\x41\x9f\x06\x0a\x90\x92\xd7\xee\x9d\x70\x15\x71\xfb\x80\xaa\x76\xd5\xde\x86\xc8\xdc\xe3\x31\xa5\x55\x6c\x1c\x6f\x77\xdd\x0c\xad\xe5\x92\xdc\x45\xa4\x26\x86\x5f\x04\xa2\x93\x5d\xfd\x1c\x01\x14\xb0\xa3\x70\xf6\xf6\x1a\x8b\x1b\x80\x3c\xc3\x1b\x1e\xbe\xae\xb9\xbc\x90\x21\x4f\xee\x4c\x41\xc0\x6c\x13\xa1\xcf\xdf\xad\x9f\x97\x3f\xc0\xe4\xc7\x88\x1f\x51\x7b\xb4\x34\x5c\x46\x47\xa2\x6f\x9a\x19\xa1\x2f\x28\xeb\x36\x64\x8d\x5c\x29\x71\x50\x6e\x98\x8f\x65\x61\xc4\xf5\x6d\x8b\x34\xfe\x19\xae\x71\x39\xab\xdd\x16\xd8\xa8\xbb\x4c\x61\x4a\xb3\xd6\xef\x00\x00\x00\xff\xff\x47\xd7\xd0\xcb\xce\x01\x00\x00")

func templates_definition_md_bytes() ([]byte, error) {
	return bindata_read(
		_templates_definition_md,
		"templates/definition.md",
	)
}

func templates_definition_md() (*asset, error) {
	bytes, err := templates_definition_md_bytes()
	if err != nil {
		return nil, err
	}

	info := bindata_file_info{name: "templates/definition.md", size: 462, mode: os.FileMode(420), modTime: time.Unix(1431429677, 0)}
	a := &asset{bytes: bytes, info:  info}
	return a, nil
}

var _templates_full_md = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\x74\x52\x3d\x6b\xc3\x30\x10\xdd\xf5\x2b\x0e\x9c\xa1\x09\xd8\xd9\x03\x9d\x1a\x0a\x85\xa4\x94\x36\x74\x09\x1d\x14\xe7\x9c\x1a\x6c\xc5\xd8\x5a\x8a\xac\xff\xde\xd3\x49\xb6\x95\xd4\xd5\xa2\xfb\x7e\xef\x9e\x94\x80\x31\xd9\x4e\x9e\xb0\xb2\x56\x88\x14\x5e\xb6\x1b\x17\x79\x95\x35\x5a\xbb\x26\xeb\x13\xdb\xae\xbc\x2a\xca\xa6\x70\x90\xa7\x0a\x3b\x57\x50\xa1\x82\xcc\xbb\xae\x2f\x49\x42\x4e\x18\xd3\x4a\x75\x41\x58\x68\xd8\x3c\x0e\x25\xd9\xae\xec\xb4\xb5\x29\x1c\x8d\x59\xe8\x30\xfc\xeb\x21\x89\xbc\x25\x75\xa2\x3a\xd3\xb0\x69\xc2\x9f\x01\x84\xe3\xf8\x8e\x4d\x42\xb0\xb3\xc5\x2e\x6f\xcb\x46\x7b\x9a\xc4\xf3\xb9\xc4\xea\x3c\xf2\xa4\x0a\x1f\xf0\x0d\x61\x7a\xe1\xf8\x8d\xa9\x08\x80\x11\x8a\x01\xc1\x98\xb2\xa0\xe2\xec\x1d\x0b\xae\xb4\x76\x45\x26\xb6\xa8\x72\xa7\xc4\x91\x6b\x87\xa4\x67\x7b\xb3\xdf\x3f\xc9\x25\xac\xe1\x36\xeb\xe3\xab\x41\x05\x71\xf8\x69\x10\x7a\xd8\xa1\xba\xe8\x6f\x32\xde\x5a\xcc\x4b\xf7\x12\x64\x7f\xe4\xb2\x72\xc9\x68\x71\x91\xba\xd3\xa7\xe1\x8c\x46\x64\xc7\xb1\x54\x30\xba\x03\xb1\x96\x26\xb1\xe7\xb1\x26\x7f\x84\x9c\x42\x8c\x3c\xb9\x77\xca\x0f\x62\xed\x65\xd3\x94\xea\x42\x82\x27\x2c\xe8\x56\x6a\x09\xf5\xf5\x8c\x15\xd4\x21\x25\xf6\xec\xf6\x10\xfe\x17\x59\xac\x0f\xdd\xac\x07\xdd\x4f\xd7\xba\x46\xa5\xc5\x3d\xfb\xfb\x7d\xfa\x68\xa7\xf0\xb8\x35\x3f\x6e\x4c\x84\xe8\xd6\x59\xfc\x0e\x8c\x1f\x54\xf7\xeb\xcc\xe5\xc7\xdf\x3f\x57\x32\xd7\x7c\x1b\x0b\x2b\xf0\x3f\x1a\x7e\xf7\xfc\xfd\x1b\x00\x00\xff\xff\x6e\x27\xe4\xd8\x88\x03\x00\x00")

func templates_full_md_bytes() ([]byte, error) {
	return bindata_read(
		_templates_full_md,
		"templates/full.md",
	)
}

func templates_full_md() (*asset, error) {
	bytes, err := templates_full_md_bytes()
	if err != nil {
		return nil, err
	}

	info := bindata_file_info{name: "templates/full.md", size: 904, mode: os.FileMode(420), modTime: time.Unix(1431432560, 0)}
	a := &asset{bytes: bytes, info:  info}
	return a, nil
}

var _templates_mappings_md = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\x7c\x50\x41\x8b\x86\x20\x14\xbc\xfb\x2b\x04\x3b\x05\xb9\xf7\x60\x4f\xbb\x2c\x2c\xd4\x9e\x96\xbd\x1b\xbd\x42\x50\x8b\xf4\x66\xfe\xf7\x7d\xa9\xd5\xf7\x41\x7d\x5e\x46\xc6\x71\x66\xde\x63\xd4\x7b\xde\x88\x0e\x54\x08\x84\x94\xdf\x9f\xf5\x46\xfc\x08\x0d\x21\xbc\xe1\xed\x0f\x16\x2b\x27\x13\x42\x49\x08\x63\xf4\x57\x74\x0a\x2c\xf1\x7e\x11\x66\x04\x5a\x38\x5a\xbf\x53\x9e\x58\xde\x48\xeb\xd0\x85\xb1\xcd\xb4\x70\xd9\xe5\x14\x0f\x9b\x18\xf9\x2f\x09\xaa\x3f\xe4\xde\xcb\x01\xdf\x78\x2b\xe6\x59\x9a\xd1\x26\x87\x68\x31\xec\x16\xed\xd4\x83\xa2\x2b\xcd\x6d\xf0\x16\x23\x11\xa3\x17\xe2\xc7\xa4\x35\x18\x47\xaa\x78\xd6\x6a\x3f\xeb\x0d\x56\x67\x2d\x1d\x6b\x3d\x16\xc0\x68\x9d\x5a\xa6\xd1\x78\xcc\xcf\x65\x30\xec\xe6\xfd\xd8\xd5\x95\xe4\xea\xf3\x33\x97\x47\x88\x2b\x01\xd3\xbf\xc0\xff\x00\x00\x00\xff\xff\xad\x77\x2b\x0e\xb5\x01\x00\x00")

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

	info := bindata_file_info{name: "templates/mappings.md", size: 437, mode: os.FileMode(420), modTime: time.Unix(1431427248, 0)}
	a := &asset{bytes: bytes, info:  info}
	return a, nil
}

var _templates_models_md = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\x52\x56\xf0\xcd\x4f\x49\xcd\x29\xe6\xaa\xae\x2e\x4a\xcc\x4b\x4f\x55\x50\xc9\x55\xb0\xb2\x55\xd0\xf3\xc9\x2c\x2e\xa9\xad\xd5\x55\xa8\xae\x56\xc9\xd5\xf3\x49\x4c\x4a\xcd\xa9\xad\x05\xaa\x49\xcd\x4b\x01\xd2\x80\x00\x00\x00\xff\xff\x2e\xd9\x1a\x4a\x35\x00\x00\x00")

func templates_models_md_bytes() ([]byte, error) {
	return bindata_read(
		_templates_models_md,
		"templates/models.md",
	)
}

func templates_models_md() (*asset, error) {
	bytes, err := templates_models_md_bytes()
	if err != nil {
		return nil, err
	}

	info := bindata_file_info{name: "templates/models.md", size: 53, mode: os.FileMode(420), modTime: time.Unix(1431434112, 0)}
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
	"templates/definition.md": templates_definition_md,
	"templates/full.md": templates_full_md,
	"templates/mappings.md": templates_mappings_md,
	"templates/models.md": templates_models_md,
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
		"definition.md": &_bintree_t{templates_definition_md, map[string]*_bintree_t{
		}},
		"full.md": &_bintree_t{templates_full_md, map[string]*_bintree_t{
		}},
		"mappings.md": &_bintree_t{templates_mappings_md, map[string]*_bintree_t{
		}},
		"models.md": &_bintree_t{templates_models_md, map[string]*_bintree_t{
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
