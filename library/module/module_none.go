//go:build !bindata
// +build !bindata

package module

import (
	bindataBackend "github.com/coscms/webcore/library/bindata"
	"github.com/coscms/webcore/library/httpserver"
	"github.com/coscms/webcore/library/module"
	"github.com/coscms/webfront/library/bindata"
)

func SetFrontendTemplate(key string, templatePath string) {
	module.SetTemplate(bindataBackend.PathAliases, httpserver.Frontend.Template.PathAliases, key, templatePath)
}

func SetFrontendAssets(assetsPath string) {
	module.SetAssets(bindataBackend.StaticOptions, bindata.StaticOptions, assetsPath)
}
