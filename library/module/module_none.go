//go:build !bindata
// +build !bindata

package module

import (
	"github.com/coscms/webcore/library/module"
	"github.com/coscms/webfront/initialize/frontend"
	"github.com/coscms/webfront/library/bindata"
)

func SetFrontendTemplate(key string, templatePath string) {
	module.SetTemplate(frontend.TmplPathFixers.PathAliases, key, templatePath)
}

func SetFrontendAssets(assetsPath string) {
	module.SetAssets(bindata.StaticOptions, assetsPath)
}
