//go:build !bindata
// +build !bindata

package bindata

import (
	"path/filepath"
	"strings"

	"github.com/webx-top/echo"
	"github.com/webx-top/echo/middleware"
	"github.com/webx-top/echo/middleware/render/driver"

	"github.com/admpub/log"
	"github.com/coscms/webcore/cmd/bootconfig"
	"github.com/coscms/webcore/initialize/backend"
	"github.com/coscms/webcore/library/bindata"
	"github.com/coscms/webcore/library/httpserver"
	"github.com/coscms/webcore/library/modal"
	"github.com/coscms/webfront/initialize/frontend"
)

var (
	// StaticOptions 前台static中间件选项
	StaticOptions = &middleware.StaticOptions{
		Root:   "",
		Path:   "",
		MaxAge: bootconfig.HTTPCacheMaxAge,
	}
	NgingDir = `../nging`
	WebxDir  = `../webx`
)

func init() {
	httpserver.Frontend.StaticOptions = StaticOptions
}

// PrependBackendAssetsDir 往前插入后台素材文件夹
func PrependBackendAssetsDir(assetsDir string) {
	oldRoot := bindata.StaticOptions.Root
	bindata.StaticOptions.Root = assetsDir
	if len(oldRoot) == 0 {
		httpserver.Backend.AssetsDir = filepath.Join(NgingDir, backend.DefaultAssetsDir) //素材文件夹
		oldRoot = httpserver.Backend.AssetsDir
	}
	bindata.StaticOptions.AddFallback(oldRoot)
}

// AppendBackendAssetsDir 追加后台素材文件夹
func AppendBackendAssetsDir(assetsDir string) {
	bindata.StaticOptions.AddFallback(assetsDir)
}

// PrependFrontendAssetsDir 往前插入前台素材文件夹
func PrependFrontendAssetsDir(assetsDir string) {
	oldRoot := StaticOptions.Root
	StaticOptions.Root = assetsDir
	if len(oldRoot) == 0 {
		httpserver.Frontend.AssetsDir = filepath.Join(WebxDir, frontend.DefaultAssetsDir) //素材文件夹
		oldRoot = httpserver.Frontend.AssetsDir
	}
	StaticOptions.AddFallback(oldRoot)
}

// AppendFrontendAssetsDir 追加前台素材文件夹
func AppendFrontendAssetsDir(assetsDir string) {
	StaticOptions.AddFallback(assetsDir)
}

// Initialize 后台和前台模板等素材初始化配置
func Initialize(callbacks ...func()) {
	frontend.AutoBackendPrefix()
	httpserver.Backend.AssetsDir = filepath.Join(NgingDir, `public/assets/backend`)
	httpserver.Backend.TemplateDir = filepath.Join(NgingDir, `template/backend`)
	bindata.StaticOptions.AddFallback(filepath.Join(WebxDir, `public/assets/backend`))
	if len(callbacks) > 0 && callbacks[0] != nil {
		callbacks[0]()
	}
	bindata.Initialize() // 后台素材初始化配置
	backendTemplateDir, err := filepath.Abs(filepath.Join(WebxDir, `template/backend`))
	if err != nil {
		panic(err)
	}
	log.Debug(`[backend] `, `Template subfolder "official" is relocated to: `, backendTemplateDir)
	httpserver.Backend.Template.AddDir(`official`, backendTemplateDir)

	// 应用后台模块的文件别名分组到后台模板路径修正器
	httpserver.Backend.Template.ApplyAliases()

	httpserver.Backend.RendererDo = func(renderer driver.Driver) {
		httpserver.Backend.Template.Register(renderer, backendTemplateDir)
	}
	modal.PathFixer = func(c echo.Context, file string) string {
		file = strings.TrimPrefix(file, httpserver.Backend.TemplateDir+`/`)
		return httpserver.Backend.Template.Handle(c, ``, file)
	}
	httpserver.Frontend.TemplateDir = filepath.Join(WebxDir, frontend.DefaultTemplateDir) //模板文件夹
	httpserver.Frontend.AssetsDir = filepath.Join(WebxDir, frontend.DefaultAssetsDir)     //素材文件夹
	//注册前台静态资源
	frontendTemplateDir := filepath.Join(WebxDir, `template/frontend`)
	httpserver.Frontend.Template.PathAliases.AddAllSubdir(frontendTemplateDir)
	//httpserver.Frontend.Template.PathAliases.Add(`default`, frontendTemplateDir)

	// 应用前台模块的文件别名分组到前台模板路径修正器
	httpserver.Frontend.Template.ApplyAliases()

	httpserver.Frontend.RendererDo = func(renderer driver.Driver) {
		httpserver.Frontend.Template.SetEnableTheme(true).Register(renderer)
	}
}
