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
		file = path.Join(backend.AssetsDir, `backend`, file)
	} else {
		file = path.Join(frontend.AssetsDir, `frontend`, file)
	}
	return bindata.StaticAssetFS.Open(file)
}
