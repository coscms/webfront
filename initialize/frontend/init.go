package frontend

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/admpub/log"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/engine/mock"
	"github.com/webx-top/echo/handler/captcha"
	"github.com/webx-top/echo/middleware"
	"github.com/webx-top/echo/middleware/language"
	"github.com/webx-top/echo/middleware/render"
	"github.com/webx-top/echo/middleware/render/driver"
	"github.com/webx-top/echo/middleware/session"
	"github.com/webx-top/echo/subdomains"
	"github.com/webx-top/validator"

	"github.com/coscms/webcore/cmd/bootconfig"
	"github.com/coscms/webcore/initialize/backend"
	backendLib "github.com/coscms/webcore/library/backend"
	"github.com/coscms/webcore/library/common"
	"github.com/coscms/webcore/library/config"
	ngingMW "github.com/coscms/webcore/middleware"
	"github.com/coscms/webfront/dbschema"
	"github.com/coscms/webfront/library/frontend"
	"github.com/coscms/webfront/library/routepage"
	"github.com/coscms/webfront/library/xmetrics"
	xMW "github.com/coscms/webfront/middleware"
	"github.com/coscms/webfront/model/official"
)

const (
	Name                  = `frontend`
	DefaultTemplateDir    = `./template/` + Name
	DefaultAssetsDir      = `./public/assets/frontend`
	DefaultAssetsURLPath  = `/public/assets/frontend`
	RouteDefaultExtension = `.html`
)

var (
	StaticMW           interface{}
	TemplateDir        = DefaultTemplateDir //模板文件夹
	AssetsDir          = DefaultAssetsDir   //素材文件夹
	AssetsURLPath      = DefaultAssetsURLPath
	StaticRootURLPath  = `/public/`
	RendererDo         = func(driver.Driver) {}
	TmplCustomParser   func(tmpl string, content []byte) []byte
	DefaultMiddlewares = []interface{}{}
)

func init() {
	config.AddConfigInitor(func(c *config.Config) {
		c.AddReloader(func(newConfig *config.Config) {
			c.Sys.ReloadRealIPConfig(&newConfig.Sys, IRegister().Echo().RealIPConfig())
		})
	})
	prefix := os.Getenv(`NGING_FRONTEND_URL_PREFIX`)
	if len(prefix) > 0 {
		SetPrefix(prefix)
	}
	bootconfig.OnStart(1, InitWebServer)
}

func SetPrefix(prefix string) {
	IRegister().SetPrefix(prefix)
	AssetsURLPath = prefix + DefaultAssetsURLPath
	StaticRootURLPath = prefix + `/public/`
	frontend.AssetsURLPath = AssetsURLPath
}

func InitWebServer() {
	var frontendDomain string
	siteURL := config.Setting(`base`).String(`siteURL`)
	if len(siteURL) > 0 {
		info, err := url.Parse(siteURL)
		if err != nil {
			log.Error(siteURL, `: `, err)
		} else {
			frontendDomain = info.Scheme + `://` + info.Host
		}
	}
	e := IRegister().Echo()
	config.FromFile().Sys.SetRealIPParams(IRegister().Echo().RealIPConfig())
	e.SetRenderDataWrapper(xMW.DefaultRenderDataWrapper)
	e.SetDefaultExtension(RouteDefaultExtension)
	if len(config.FromCLI().BackendDomain) > 0 {
		// 如果指定了后台域名则只能用该域名访问后台。此时将其它域名指向前台
		subdomains.Default.Default = Name // 设置默认(没有匹配到域名的时候)访问的域名别名
	}
	if len(frontendDomain) == 0 {
		frontendDomain = config.FromCLI().FrontendDomain
	} else if len(config.FromCLI().FrontendDomain) > 0 {
		info := strings.SplitN(config.FromCLI().FrontendDomain, `,`, 2)
		if len(info) > 1 {
			frontendDomain += `,` + info[1]
		}
	}
	if len(frontendDomain) == 0 {
		if len(config.FromCLI().BackendDomain) == 0 {
			// 前后台都没有指定域名的时候，给前台强制指定一个域名
			frontendDomain = `localhost:` + fmt.Sprintf(`%d`, config.FromCLI().Port)
		}
	} else {
		var domains []string
		for _, domain := range backend.MakeSubdomains(frontendDomain, []string{}) {
			if _, ok := subdomains.Default.Hosts[domain]; !ok { // 排除后台指定的域名
				domains = append(domains, domain)
			}
		}

		frontendDomain = strings.Join(domains, `,`)
	}
	subdomains.Default.Add(Name+`@`+frontendDomain, e)
	addMiddleware(e)
	log.Infof(`Registered host: %s@%s`, Name, frontendDomain)
	e.Get(`/favicon.ico`, faviconHandler)
	e.Use(xMW.SessionInfo)
	if config.IsInstalled() {
		routepage.Apply(e, frontend.GlobalFuncMap())
		applyRouteRewrite(e)
	}
	Apply()
}

var renderOptions *render.Config

func addMiddleware(e *echo.Echo) {
	if bootconfig.Bindata {
		e.SetDebug(false)
	} else {
		e.SetDebug(true)
	}
	e.Use(middleware.Recover())
	e.Use(xMW.HostChecker())
	e.Use(ngingMW.MaxRequestBodySize)
	if len(DefaultMiddlewares) == 0 {
		if !config.FromFile().Sys.DisableHTTPLog {
			e.Use(middleware.Log())
		}
	} else {
		e.Use(DefaultMiddlewares...)
	}
	if StaticMW != nil {
		e.Use(StaticMW)
	}
	e.Use(bootconfig.StaticMW) //后台静态资源(在bindata模式下也包含了前台静态资源)

	// Prometheus
	xmetrics.Register(e)

	// 启用session
	e.Use(session.Middleware(config.SessionOptions))

	// 启用多语言支持
	i18n := language.New(&config.FromFile().Language)
	e.Use(i18n.Middleware())

	// 启用Validation
	e.Use(validator.Middleware())

	// 事物支持
	e.Use(ngingMW.Transaction())

	// 注册模板引擎
	if renderOptions != nil && renderOptions.Renderer() != nil {
		renderOptions.Renderer().Close()
	}
	renderOptions = &render.Config{
		TmplDir: TemplateDir,
		Engine:  `standard`,
		ParseStrings: map[string]string{
			`__TMPL__`: TemplateDir,
		},
		DefaultHTTPErrorCode: http.StatusOK,
		Reload:               true,
		ErrorPages:           config.FromFile().Sys.ErrorPages,
		ErrorProcessors:      common.ErrorProcessors,
		FuncMapGlobal:        frontend.GlobalFuncMap(),
	}
	if echo.String(`LABEL`) != `dev` && config.FromFile().Extend.GetStore(`minify`).Bool(`on`) {
		renderOptions.CustomParser = TmplCustomParser
	}
	if RendererDo != nil {
		renderOptions.AddRendererDo(RendererDo)
	}
	renderOptions.AddFuncSetter(FrontendURLFunc, ngingMW.ErrorPageFunc, xMW.SetFunc)
	renderOptions.ApplyTo(e, bootconfig.FrontendTmplMgr)
	echo.OnCallback(`nging.renderer.cache.clear`, func(_ echo.Event) error {
		log.Debug(`clear: Frontend Template Object Cache`)
		renderOptions.Renderer().ClearCache()
		return nil
	}, `clear-frontend-template-object-cache`)
	echo.OnCallback(`webx.renderer.cache.clear`, func(_ echo.Event) error {
		log.Debug(`clear: Frontend Template Object Cache`)
		renderOptions.Renderer().ClearCache()
		return nil
	}, `clear-frontend-template-object-cache`)
	echo.OnCallback(`webx.frontend.close`, func(_ echo.Event) error {
		log.Debug(`close: Frontend Template Manager`)
		renderOptions.Renderer().Close()
		renderOptions = nil
		return nil
	}, `close-frontend-template-manager`)
	e.Use(xMW.UseTheme(TmplPathFixers.ThemeInfo))
	e.Use(FrontendURLFuncMW())
	e.Use(ngingMW.FuncMap())
	e.Use(xMW.FuncMap())
	e.Use(render.Auto())

	keepExtensionPrefixes := []string{StaticRootURLPath}
	if config.IsInstalled() {
		ctx := echo.NewContext(mock.NewRequest(), mock.NewResponse(), e)
		routeM := official.NewRoutePage(ctx)
		routes, _ := routeM.ListWithExtensionRoutes(RouteDefaultExtension)
		keepExtensionPrefixes = append(keepExtensionPrefixes, routes...)
	}
	e.Pre(xMW.TrimPathSuffix(keepExtensionPrefixes...))

	// - verifier or guard -

	// RateLimiter
	e.Use(xMW.RateLimiter())

	// IPFilter
	e.Use(xMW.IPFilter())

	// 启用JWT
	e.Use(xMW.JWT())

	e.Use(xMW.Middlewares...)

	captcha.New(``).Wrapper(e)
	e.Route("GET", `/qrcode`, backendLib.QrCode)

	i18n.Handler(e, `App.i18n`)
}

func URLFor(purl string) string {
	return xMW.URLFor(purl)
}

func Customer(c echo.Context) *dbschema.OfficialCustomer {
	return xMW.Customer(c)
}
