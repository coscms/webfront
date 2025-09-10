//go:build bindata
// +build bindata

package minify

import (
	"net/http"
	"path"

	"github.com/coscms/webcore/library/bindata"
	"github.com/coscms/webcore/library/httpserver"
	"github.com/webx-top/echo"
)

func openfile(asset string, file string) (http.File, error) {
	file = echo.CleanPath(file)
	if asset == `AssetsURL` {
		file = path.Join(httpserver.Backend.AssetsDir, file)
	} else {
		file = path.Join(httpserver.Frontend.AssetsDir, file)
	}
	return bindata.StaticAssetFS.Open(file)
}
