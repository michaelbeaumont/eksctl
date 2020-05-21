// Code generated by go-bindata. DO NOT EDIT.
// sources:
// assets/schema.json (20.357kB)

package v1alpha5

import (
	"bytes"
	"compress/gzip"
	"crypto/sha256"
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
		return nil, fmt.Errorf("read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("read %q: %v", name, err)
	}
	if clErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes  []byte
	info   os.FileInfo
	digest [sha256.Size]byte
}

type bindataFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

func (fi bindataFileInfo) Name() string {
	return fi.name
}
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}
func (fi bindataFileInfo) IsDir() bool {
	return false
}
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _schemaJson = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xec\x5c\x5f\x6f\xe3\xb8\x11\x7f\xcf\xa7\x10\x7c\x7d\x2a\xf6\x36\x2d\x70\xbd\x87\x7d\x73\x9c\xdc\xae\xb1\x9b\xc4\x88\x0f\x5b\xa0\xc5\x3d\x8c\xa8\xb1\xcc\x0b\x45\xaa\xe4\xc8\x89\xaf\xd8\xef\x5e\xc8\xb2\x12\xc9\x12\x65\x99\x96\x64\xa7\xb7\x0f\x0b\x2c\x64\x0d\xf9\x1b\xce\xff\xe1\x28\xff\xbd\xf0\xbc\xd1\x5f\x0c\x5b\x62\x04\xa3\x0f\xde\x68\x49\x14\x7f\xb8\xbc\xfc\xdd\x28\xf9\x63\xf6\xf4\xbd\xd2\xe1\x65\xa0\x61\x41\x3f\xfe\xed\xa7\xcb\xec\xd9\x0f\xa3\x77\x1b\x3a\x8d\x8b\x94\xe8\x87\xcb\x00\x17\x5c\x72\xe2\x4a\x9a\xcb\x89\x48\x0c\xa1\x9e\x28\xb9\xe0\x61\xf6\x62\xe1\xe7\xd1\x07\x2f\xdd\xd3\xf3\x46\xf9\x7b\x42\x25\xc1\x3f\x81\xd8\xf2\xe5\x27\xcf\x1b\xc5\x5a\xc5\xa8\x89\xa3\x29\x3c\xf5\xbc\x11\xcb\x88\xbe\xa8\x30\xe4\x32\x2c\xfd\xe6\xce\xc8\x2b\x7d\x23\x43\x2f\x40\xf3\xdd\x5f\x48\xbf\x6d\xff\xf7\x2d\x5f\x6d\x04\x41\xb0\xa1\x06\x31\x2b\x72\xb2\x00\x61\xf0\xe5\x25\x5a\xc7\x98\x6e\xa7\xfc\xdf\x91\x51\xb6\xdc\x76\x89\xea\xe9\x54\x59\xb6\x1e\x12\x4a\xf0\x05\xfe\xba\x8e\x77\x7e\xf0\xbc\x11\x27\x8c\x76\x1f\x16\x90\x18\xd2\x25\xc6\x0a\x2c\x95\xde\x03\xad\x61\x3d\x08\xff\x99\x16\xb5\x60\xfa\x91\xcb\x60\x97\xdb\x1d\xb6\x4a\xac\xe0\x33\x44\xb1\xd8\xac\xf2\xef\xf2\x61\x94\xb7\x2e\xfc\xf6\xdb\x45\xcd\xa9\x8c\x20\xe6\x5f\x51\x1b\xae\x64\x37\xdb\xe3\xa3\x61\x24\xde\x73\x75\xb9\xfa\x3b\x88\x78\x09\xff\xd8\x8f\x21\x42\x82\x00\x08\x06\x35\x88\x5b\x24\x18\xd5\xc2\xe1\x10\x0d\x8a\x64\x3a\xbe\xad\x07\xb2\x8a\xd9\xa0\x40\xbe\xce\x26\xf5\x40\xa4\x0a\xf0\xa3\x56\x49\xdc\xd6\x22\x8f\x84\x69\x05\x7a\x97\x03\x39\xdc\xc8\x8b\xfa\x06\x12\x42\x0c\xee\xce\x85\xab\xdb\x1d\x3c\x47\x31\xb7\x00\x1d\x02\xe1\x4c\xab\x05\x17\xad\x5d\x68\x5f\xac\xfd\x52\x42\x73\x14\x63\xb0\x02\x2e\xc0\xe7\x82\xd3\xfa\x5f\x4a\x0e\x16\x1d\xde\x15\xe3\x77\x4d\xb4\xef\xe2\xfc\xda\xc6\xee\x7a\x58\x06\x99\x46\x32\x37\x92\xe9\x75\x4c\x35\xae\xbc\x1f\x74\xf3\xca\xb6\xf5\xe8\x08\x28\xa9\x08\xab\xd7\x03\x9b\x67\x5b\xd6\xc2\x09\x39\x0d\x83\xe5\x23\xa7\x21\x72\x8c\x1b\x19\xc4\x8a\x4b\x32\x6d\xd2\x8c\x58\xf3\x15\x10\x8e\x19\x43\x53\x11\x49\xbe\x99\xaf\x94\x40\xb0\x48\x33\x4e\x7c\xc1\xd9\xa1\x0b\xf4\xc6\x7d\x1a\x3b\x5b\xf0\x6d\x50\xaf\x38\xc3\x07\x25\x70\xfc\x70\xb7\x27\xd3\xb1\xd8\xd8\xcb\x0a\x33\xd4\x11\x37\x69\xca\x64\xae\x54\x22\x03\xd0\x6b\x97\x15\x73\x47\xad\x82\x9b\x67\x64\x49\x7a\x1a\x47\xe0\xb3\xac\xd6\x11\xd6\x27\x4e\xcb\xfb\xe9\xf5\xc4\x49\x65\xb6\x47\x37\x66\x4c\x25\x65\x3d\xf5\x4e\x10\x91\x5e\x15\x67\x5e\xc2\x75\x46\x65\xc3\x74\x7c\xbb\xc9\x4e\x5b\x28\xb6\x84\x08\x5d\x04\x9a\xd2\x99\x18\x98\x13\xb1\x00\x1f\x45\x45\x8e\x31\x10\xa1\x96\xb3\x7a\xa4\x9b\x57\xde\xff\xb5\xf2\xac\x31\x40\xbf\x1e\xb2\x5d\x2a\xc5\x73\xdc\x05\x0a\x52\x2a\x82\x72\xe5\x7e\x66\x68\xfb\x54\xa2\x1d\x05\x6f\xa1\x4e\x27\xa9\xc4\x72\x75\xaf\x17\x21\x11\xb0\xe5\x4c\x09\xce\xd6\xe3\x87\xbb\x13\x24\x7d\x45\x04\x9d\x69\x91\x45\xca\xa4\x13\xec\x50\xff\xe3\x6e\xbc\xff\x09\x32\xb8\x8a\xf6\x36\xa5\x74\x04\xe1\x77\xfb\xce\x4f\xa8\x85\x95\xeb\x03\xd3\x8c\xde\x78\xd8\x8d\x72\x1a\xff\x93\x70\x8d\x41\xa9\xad\x94\xc5\xb8\x82\xb8\x35\x86\x85\x62\xe3\xb7\x77\x7d\xc5\xc8\xed\x3e\x0e\x94\xab\x56\x1d\xb5\xef\xba\x9c\xeb\xc1\xdd\xf8\xd7\x36\x7a\x9b\x26\xb8\x4f\xd0\xde\x85\xf5\x86\xb7\xbd\xa9\xe1\xb6\x3a\x73\x51\x05\x96\x2e\xb8\xe0\x2c\x2d\xd8\x12\x5a\x2a\xcd\x69\x7d\x5d\x13\x9c\x9b\x1a\xb5\x11\x06\x7c\x97\xc0\xf3\x46\x3e\x97\xa0\xd7\x37\x92\xa9\x20\xeb\xca\x8f\x7c\x30\xf8\xf3\x4f\xa5\x40\x59\x1f\x0d\xb5\x93\x5a\x1b\x02\xf6\x78\x77\x88\x21\xf6\x27\xbc\xc4\x97\x78\x50\xb9\xdc\x99\x3d\x1e\x5d\xd3\xd8\x1b\xa3\x48\x4f\x4a\x3f\x76\x99\x3c\x67\x95\x7e\x77\xbc\xf7\x8a\xbb\x37\x6d\xf9\x3a\x9b\xb4\xd1\x14\xbe\xef\xf6\xa6\xde\xc2\x79\xa0\xdd\x9a\x11\x2c\x49\xbd\x41\xd6\x2d\x76\x59\xa0\x62\x03\x5e\xef\x29\x5d\x6e\x77\xb5\x80\xf0\x99\x34\x4c\xa6\xd7\x0f\x27\x48\xf0\xcd\x12\x74\xd6\x7d\x9f\x1f\x7b\xae\x90\x90\x1a\x0b\xa1\x52\x9f\x3d\x9d\xad\x7e\x76\xea\x97\x48\x18\xa8\x43\x59\x88\xbe\xf5\xea\x69\xef\x30\xf6\x8f\xea\x75\xd7\x06\xef\x94\xf5\x21\x87\xd4\x9b\xee\x1d\xcd\xce\x15\x49\xbb\x74\xb8\xf7\xe4\x37\xee\xa6\x2d\x69\x50\x20\x23\xa5\x4f\xdd\xf5\x2b\x1f\xf2\x7c\x8b\xea\x38\xb7\x51\xef\x43\xfb\xf7\x57\x7f\xae\x12\xc1\x22\xb9\x56\x66\x92\x75\x38\x0f\xb2\x95\xff\xd3\xa6\x68\xf7\x72\xf9\xc8\x5b\x75\x14\x35\xc6\x6a\x98\xb8\xf1\x90\xee\x54\x2b\x9a\x14\x19\x94\x95\xa6\x47\x1c\xf7\xf9\x6e\xb5\x58\x7c\xa5\xc8\x90\x86\xb8\xea\xef\x7b\xc4\x54\xb9\x7f\xef\x41\x21\x2a\x53\x0c\x67\x12\xc9\x20\xe2\xbf\x40\xc4\x85\x53\xeb\x93\x4b\x43\x20\xd9\x66\x1c\xcd\x85\x3e\x40\x93\xf2\x3d\x81\x18\x18\x27\x2b\x04\x2e\x09\x43\xb4\xa8\x4c\xc4\xe5\x9c\xff\x61\xdd\xbe\x99\x16\x9e\x9d\x69\x57\x4a\x24\x11\x3a\x93\x9f\xc1\xbc\x86\x31\xd5\x41\x8d\xe6\xb9\xa2\xf9\xfc\xd3\x5b\x74\xef\xc5\xc4\x2d\xeb\x59\x6c\x0b\xeb\x9a\x31\xd3\x56\x35\xc8\x39\xe7\x18\x45\xf3\xac\x19\xd5\x6b\x16\x70\x79\xd8\xae\x07\x2f\x88\xa4\x39\x33\x13\x25\xd2\x44\xa5\xdc\x04\xb6\xb8\xc1\x50\x83\x4c\x04\xa4\x85\x67\x7b\x6f\x58\x24\x72\xf0\x4b\x51\x06\xf3\xcd\x16\x4d\x79\xdf\xe8\x3c\x5a\x33\x3d\xf0\x77\x86\x51\xf4\xcd\x05\xdf\x9c\xde\x5c\xf3\xf4\x35\x3f\x19\x6e\x36\xee\xd5\xdf\xd4\x62\x38\xd7\x80\xf9\x46\xdc\x7e\x37\x31\xae\xd4\x45\x1d\xa8\xb9\xf5\x9a\x68\x7c\xb4\x34\xb6\xde\x72\xc6\x08\x26\x6c\x8a\x7f\x5e\x83\x12\xf7\x36\x5b\x5d\xc1\x73\x8c\x7d\xa0\x6f\xee\x63\xe2\x11\xff\x03\xad\x51\xa5\x51\xe7\x8e\xcc\xa9\x33\x72\x57\x7f\x98\x51\x1f\x74\x0b\x57\xa1\xde\x4e\xfb\x1e\xc5\xfe\xe7\xc8\x7c\xc6\xf5\xf4\xda\x1d\xc5\xf4\x7e\x36\x77\xd5\xee\x99\x0a\xcc\x0c\x75\x6a\x89\x4e\x4b\xbc\x99\x4a\x80\xa0\xae\x69\x7f\x86\x40\x99\x00\x63\x38\xfb\xa2\x20\xb8\x02\x91\x46\x4b\x9d\x2a\xe9\x49\xe2\x9f\x0e\x91\x36\x0e\xfa\x34\xc3\x66\x75\x15\x6b\xcf\x81\xc8\x56\xf1\x0e\xf6\x19\x94\xa5\x36\x2b\x35\xcd\x88\x04\x6a\xc5\x1e\x71\xa0\xab\xb1\x17\x4c\x57\xc5\xad\x2d\xc9\x08\x5e\xe5\x4d\xbd\x89\x8a\x22\x90\xc1\x09\x14\x47\xad\x50\x6b\x1e\x54\xa0\x38\xd5\x3e\xd9\x0d\xdc\xf5\x9d\xd5\xc7\x36\x51\x3f\x26\x3e\x0a\xa4\x9b\xcd\x8d\xee\xee\x07\x90\xde\x51\x2e\xa8\xc7\xa9\xc9\x8b\x9d\xb7\x3b\x2c\xe1\xae\xea\xb5\x77\xcf\x07\xb0\xe3\x20\xe2\x72\xa2\x64\xea\xc5\xd1\x5a\x8a\xee\x49\x70\x89\xb8\xec\x30\xa3\x7f\x9b\xe7\xdf\xf2\x1b\x91\xd3\x4f\x1a\xe7\xf5\xea\xb6\x5b\xef\x78\xe5\x9a\xaf\x72\xc4\xad\x6d\x71\x09\xd7\x74\xb1\xb8\x46\x87\x5f\xa1\x8c\x83\x40\xc9\x8d\x90\xaa\x7a\x3b\x40\x7c\x2a\x6f\x3f\x94\xfa\xda\x98\xb6\xf7\xbc\x22\x08\xf1\x2a\xe1\x22\x70\x74\x1d\x90\x90\x9a\x33\x10\x8e\xe4\xf8\x9c\x3a\x18\x10\x0d\x21\xa4\x91\x9e\xa1\xa6\xec\x3e\xc9\x11\x7e\x1c\xdf\x62\x35\x95\x6a\x87\xdd\x77\xfb\x32\x6e\x61\x9e\xdd\xf6\x5b\xb8\xed\x07\xc2\x9f\xca\x50\xbb\x7e\xc8\xf7\xfc\x60\x9f\xa4\x6d\x96\x8d\xf5\x23\x58\x3b\x79\x9f\xe6\xb1\xa7\xc9\x67\xe9\x9a\x16\x7b\x8b\xa6\x7d\xfb\xb4\x4c\x36\x78\x84\x48\xcb\x57\xcd\xed\x03\x0a\x32\x89\x7c\x5b\xdd\xaa\xe4\x35\xa6\x19\xe1\x15\x18\x3c\xaa\xbb\x94\x2f\x34\x43\xcd\x50\x12\x84\x38\xf6\xd5\x0a\x8f\x5e\xd7\xc4\x8a\x72\x69\xce\x94\xaa\xd6\xd8\xad\x57\xd9\x8e\xe0\x71\x25\xe7\xa4\x81\x30\x3c\xe5\xcc\x78\xa9\xe5\xd7\x3a\x0f\x99\x5e\x9f\x40\xbd\xd2\x00\x3b\xdf\x4c\x43\x3a\x79\x86\x94\xfc\x8b\x62\x20\xce\xc2\x31\xa4\x85\x6d\x9b\xe3\x16\x42\x3d\x39\xb1\x9b\x8d\x20\x7e\xc6\xf5\x0c\xc8\xea\x0b\x1b\x67\xec\xf2\x05\x8e\x22\x76\x4d\xd0\x8c\x4a\x34\x2b\x8f\xbc\x4e\x07\xab\x5c\xbb\x97\xfb\x7d\x75\xc8\xc7\x2a\xf3\x4d\x0b\xcf\xe5\xcc\x8e\x9a\x10\x4b\xcd\xe3\x13\x8a\x4a\x57\x65\x48\xeb\xa8\x99\x32\xb5\x7e\x33\xbf\xd1\x0f\x17\x46\x35\xae\xb8\xeb\xf7\x4f\x2a\xa1\x38\xa1\x83\x2c\xaa\xfb\x53\x7a\x28\x0f\xad\x59\x8f\x28\xd1\x4e\x5a\xe4\x6b\x90\xf6\xec\xa9\xd1\xe8\x81\x96\x27\x88\x0b\x0b\x91\x3c\xbb\x3a\xb9\xc4\xd8\x73\xf8\x26\x3a\x8c\x80\x3b\x9d\xee\xf6\x72\x70\x3e\xff\x74\xa8\x6b\xee\x5e\x91\xaa\x7f\x22\x65\x6f\x62\xfa\x88\xeb\xb4\x6a\xdf\x3e\xd8\x9f\x91\x6e\xdf\x3f\x01\x8f\x17\xe9\xbf\x6f\x17\xff\x0b\x00\x00\xff\xff\xba\xc9\x1a\xb4\x85\x4f\x00\x00")

func schemaJsonBytes() ([]byte, error) {
	return bindataRead(
		_schemaJson,
		"schema.json",
	)
}

func schemaJson() (*asset, error) {
	bytes, err := schemaJsonBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "schema.json", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0x52, 0x82, 0xd9, 0x2a, 0x5, 0x96, 0x12, 0x1f, 0x6e, 0x52, 0x95, 0xc1, 0x7a, 0x88, 0xac, 0xb1, 0xb2, 0xa8, 0x45, 0xa7, 0xbe, 0xf8, 0xb8, 0x15, 0xb2, 0xb, 0x75, 0xd2, 0xd3, 0xf3, 0x68, 0x59}}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	canonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[canonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// AssetString returns the asset contents as a string (instead of a []byte).
func AssetString(name string) (string, error) {
	data, err := Asset(name)
	return string(data), err
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

// MustAssetString is like AssetString but panics when Asset would return an
// error. It simplifies safe initialization of global variables.
func MustAssetString(name string) string {
	return string(MustAsset(name))
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	canonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[canonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetDigest returns the digest of the file with the given name. It returns an
// error if the asset could not be found or the digest could not be loaded.
func AssetDigest(name string) ([sha256.Size]byte, error) {
	canonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[canonicalName]; ok {
		a, err := f()
		if err != nil {
			return [sha256.Size]byte{}, fmt.Errorf("AssetDigest %s can't read by error: %v", name, err)
		}
		return a.digest, nil
	}
	return [sha256.Size]byte{}, fmt.Errorf("AssetDigest %s not found", name)
}

// Digests returns a map of all known files and their checksums.
func Digests() (map[string][sha256.Size]byte, error) {
	mp := make(map[string][sha256.Size]byte, len(_bindata))
	for name := range _bindata {
		a, err := _bindata[name]()
		if err != nil {
			return nil, err
		}
		mp[name] = a.digest
	}
	return mp, nil
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
	"schema.json": schemaJson,
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
// then AssetDir("data") would return []string{"foo.txt", "img"},
// AssetDir("data/img") would return []string{"a.png", "b.png"},
// AssetDir("foo.txt") and AssetDir("notexist") would return an error, and
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		canonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(canonicalName, "/")
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
	"schema.json": &bintree{schemaJson, map[string]*bintree{}},
}}

// RestoreAsset restores an asset under the given directory.
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
	return os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
}

// RestoreAssets restores an asset under the given directory recursively.
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
	canonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(canonicalName, "/")...)...)
}
