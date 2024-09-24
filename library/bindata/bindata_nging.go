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

	"github.com/coscms/webcore/initialize/backend"
	"github.com/coscms/webcore/library/bindata"
	"github.com/coscms/webcore/library/httpserver"
	"github.com/coscms/webcore/library/modal"
	"github.com/coscms/webcore/library/ntemplate"
	"github.com/coscms/webfront/initialize/frontend"
)

func Initialize(callbacks ...func()) {
	frontend.AutoBackendPrefix()
	bindata.BackendTmplAssetPrefix = "template"
	bindata.FrontendTmplAssetPrefix = "template"
	if len(callbacks) > 0 && callbacks[0] != nil {
		callbacks[0]()
	}
	bindata.Initialize()

	httpserver.Backend.AssetsDir = backend.DefaultAssetsDir
	httpserver.Backend.TemplateDir = backend.DefaultTemplateDir
	httpserver.Backend.RendererDo = func(renderer driver.Driver) {
		httpserver.Backend.Template.SetCustomFS(bindata.BackendTmplAssetFS).Register(renderer)
	}
	httpserver.Frontend.RendererDo = func(renderer driver.Driver) {
		httpserver.Frontend.Template.SetEnableTheme(true).SetCustomFS(bindata.FrontendTmplAssetFS).Register(renderer)
	}

	if echo.String(`LABEL`) != `dev` { // 在开发环境下不启用，避免无法测试 bindata 真实效果
		// 在 bindata 模式，支持优先读取本地的静态资源文件和模板文件，在没有找到的情况下才读取 bindata 内的文件

		// StaticMW

		fileSystems := ntemplate.NewFileSystems()
		fileSystems.Register(ntemplate.NewStaticDir(filepath.Join(echo.Wd(), "public/assets"), "/public/assets")) // 注册本地文件系统内的文件
		fileSystems.Register(ntemplate.NewFileSystemTrimPrefix(frontend.Prefix(), bindata.StaticAssetFS))         // 注册 bindata 打包的文件
		httpserver.Frontend.StaticMW = mwBindata.Static(frontend.Prefix()+"/public/assets", fileSystems)

		// Template file manager

		// 后台
		backendManagers := []driver.Manager{
			manager.New(),              // 本地文件系统内的模板文件
			httpserver.Backend.TmplMgr, // bindata 打包的模板文件
		}
		backendMultiManager := ntemplate.NewMultiManager(httpserver.Backend.TemplateDir, backendManagers...)
		log.Debugf(`%s Enabled MultiManager (num: %d)`, color.GreenString(`[backend.renderer]`), len(backendManagers))
		httpserver.Backend.TmplMgr = backendMultiManager

		// 前台
		frontendManagers := []driver.Manager{
			manager.New(),               // 本地文件系统内的模板文件
			httpserver.Frontend.TmplMgr, // bindata 打包的模板文件
		}
		frontendMultiManager := ntemplate.NewMultiManager(httpserver.Frontend.TemplateDir, frontendManagers...)
		log.Debugf(`%s Enabled MultiManager (num: %d)`, color.GreenString(`[frontend.renderer]`), len(frontendManagers))
		httpserver.Frontend.TmplMgr = frontendMultiManager
	}

	modal.PathFixer = func(c echo.Context, file string) string {
		file = strings.TrimPrefix(file, httpserver.Backend.TemplateDir+`/`)
		return httpserver.Backend.Template.Handle(c, ``, file)
	}
}
