package static

import (
	"net/http"

	"github.com/elazarl/go-bindata-assetfs"
)

//go:generate go-bindata-assetfs -prefix "ui/" -ignore "\\.go" -pkg static -o bindata.go ./ui/...

// FileSystem  builds a binary file system struct
func FileSystem() http.FileSystem {
	fs := &assetfs.AssetFS{
		Asset:     Asset,
		AssetDir:  AssetDir,
		AssetInfo: AssetInfo,
		Prefix:    "",
	}
	return &binaryFileSystem{
		fs,
	}
}

type binaryFileSystem struct {
	fs http.FileSystem
}

func (b *binaryFileSystem) Open(name string) (http.File, error) {
	return b.fs.Open(name[1:])
}
