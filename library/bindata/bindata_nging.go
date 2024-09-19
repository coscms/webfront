//go:build bindata
// +build bindata

package bindata

import (
	"path/filepath"
	"strings"

	"github.com/admpub/color"
	"github.com/admpub/log"
	"github.com/webx-top/echo"
	mwBindata "github.com/webx-top/echo/middleware/bindata"
	"github.com/webx-top/echo/middleware/render/driver"
	"github.com/webx-top/echo/middleware/render/manager"

	"github.com/coscms/webcore/cmd/bootconfig"
	"github.com/coscms/webcore/initialize/backend"
	"github.com/coscms/webcore/library/bindata"
	"github.com/coscms/webcore/library/modal"
	selfBackend "github.com/coscms/webfront/initialize/backend"
	"github.com/coscms/webfront/initialize/frontend"
	"github.com/coscms/webfront/library/xtemplate"
)

func Initialize(callbacks ...func()) {
	bindata.BackendTmplAssetPrefix = "template"
	bindata.FrontendTmplAssetPrefix = "template"
	if len(callbacks) > 0 && callbacks[0] != nil {
		callbacks[0]()
	}
	bindata.Initialize()

	backend.AssetsDir = backend.DefaultAssetsDir
	backend.TemplateDir = backend.DefaultTemplateDir
	backend.RendererDo = func(renderer driver.Driver) {
		selfBackend.TmplPathFixers.SetCustomFS(bindata.BackendTmplAssetFS).Register(renderer)
	}
	frontend.RendererDo = func(renderer driver.Driver) {
		frontend.TmplPathFixers.SetEnableTheme(true).SetCustomFS(bindata.FrontendTmplAssetFS).Register(renderer)
	}

	if echo.String(`LABEL`) != `dev` { // 在开发环境下不启用，避免无法测试 bindata 真实效果
		// 在 bindata 模式，支持优先读取本地的静态资源文件和模板文件，在没有找到的情况下才读取 bindata 内的文件

		// StaticMW

		fileSystems := xtemplate.NewFileSystems()
		fileSystems.Register(xtemplate.NewStaticDir(filepath.Join(echo.Wd(), "public/assets"), "/public/assets")) // 注册本地文件系统内的文件
		fileSystems.Register(xtemplate.NewFileSystemTrimPrefix(frontend.Prefix(), bindata.StaticAssetFS))         // 注册 bindata 打包的文件
		bootconfig.StaticMW = mwBindata.Static(frontend.Prefix()+"/public/assets", fileSystems)

		// Template file manager

		// 后台
		backendManagers := []driver.Manager{
			manager.New(),             // 本地文件系统内的模板文件
			bootconfig.BackendTmplMgr, // bindata 打包的模板文件
		}
		backendMultiManager := xtemplate.NewMultiManager(backend.TemplateDir, backendManagers...)
		log.Debugf(`%s Enabled MultiManager (num: %d)`, color.GreenString(`[backend.renderer]`), len(backendManagers))
		bootconfig.BackendTmplMgr = backendMultiManager

		// 前台
		frontendManagers := []driver.Manager{
			manager.New(),              // 本地文件系统内的模板文件
			bootconfig.FrontendTmplMgr, // bindata 打包的模板文件
		}
		frontendMultiManager := xtemplate.NewMultiManager(frontend.TemplateDir, frontendManagers...)
		log.Debugf(`%s Enabled MultiManager (num: %d)`, color.GreenString(`[frontend.renderer]`), len(frontendManagers))
		bootconfig.FrontendTmplMgr = frontendMultiManager
	}

	frontend.StaticMW = nil

	modal.PathFixer = func(c echo.Context, file string) string {
		file = strings.TrimPrefix(file, backend.TemplateDir+`/`)
		return selfBackend.TmplPathFixers.Handle(c, ``, file)
	}
}
