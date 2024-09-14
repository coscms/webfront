//go:build bindata
// +build bindata

package minify

import (
	"net/http"
	"path"

	"github.com/coscms/webcore/initialize/backend"
	"github.com/coscms/webcore/library/bindata"
	"github.com/coscms/webfront/initialize/frontend"
)

func openfile(asset string, file string) (http.File, error) {
	if asset == `AssetsURL` {
		file = path.Join(backend.AssetsDir, file)
	} else {
		file = path.Join(frontend.AssetsDir, file)
	}
	return bindata.StaticAssetFS.Open(file)
}
