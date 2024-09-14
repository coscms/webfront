//go:build !bindata
// +build !bindata

package minify

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/coscms/webcore/initialize/backend"
	bindataBackend "github.com/coscms/webcore/library/bindata"
	"github.com/coscms/webfront/initialize/frontend"
	bindataFrontend "github.com/coscms/webfront/library/bindata"
)

func openfile(asset string, file string) (http.File, error) {
	if asset == `AssetsURL` {
		file = filepath.Join(backend.AssetsDir, file)
	} else {
		file = filepath.Join(frontend.AssetsDir, file)
	}
	f, err := os.Open(file)
	if err == nil {
		return f, err
	}
	if asset == `AssetsURL` {
		for _, fallback := range bindataBackend.StaticOptions.Fallback {
			file = filepath.Join(fallback, file)
			f, err = os.Open(file)
			if err == nil {
				return f, err
			}
		}
		return f, err
	}
	for _, fallback := range bindataFrontend.StaticOptions.Fallback {
		file = filepath.Join(fallback, file)
		f, err = os.Open(file)
		if err == nil {
			return f, err
		}
	}
	return f, err
}
