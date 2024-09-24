package frontend

import (
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/admpub/log"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/engine/mock"
	"github.com/webx-top/echo/handler/captcha"
	"github.com/webx-top/echo/middleware/render"
	"github.com/webx-top/echo/subdomains"

	"github.com/coscms/webcore/cmd/bootconfig"
	"github.com/coscms/webcore/initialize/backend"
	backendLib "github.com/coscms/webcore/library/backend"
	"github.com/coscms/webcore/library/config"
	"github.com/coscms/webcore/library/httpserver"
	ngingMW "github.com/coscms/webcore/middleware"
	"github.com/coscms/webfront/dbschema"
	"github.com/coscms/webfront/library/frontend"
	"github.com/coscms/webfront/library/routepage"
	"github.com/coscms/webfront/library/xmetrics"
	xMW "github.com/coscms/webfront/middleware"
	"github.com/coscms/webfront/model/official"

	_ "github.com/coscms/webfront/library/formbuilder"
)

const (
	Name                  = `frontend`
	DefaultTemplateDir    = `./template/frontend`      // 前台模板路径默认值
	DefaultAssetsDir      = `./public/assets/frontend` // 前台素材路径默认值
	DefaultAssetsURLPath  = `/public/assets/frontend`  // 前台素材网址路径默认值
	RouteDefaultExtension = `.html`                    // 前台网页扩展名默认值
)

var TmplCustomParser func(tmpl string, content []byte) []byte

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
	httpserver.Frontend.SetRouter(IRegister())
	httpserver.Frontend.RouteDefaultExtension = RouteDefaultExtension
	httpserver.Frontend.HostCheckerRegexpKey = `frontend.hostRuleRegexp`
	httpserver.Frontend.FuncSetters = []func(echo.Context) error{
		FrontendURLFunc,
		xMW.SetFunc,
	}
	bootconfig.OnStart(0, start)
}

func start() {
	httpserver.Frontend.GlobalFuncMap = frontend.GlobalFuncMap()
	InitWebServer()
}

func SetPrefix(prefix string) {
	httpserver.Frontend.SetPrefix(prefix)
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
	e := httpserver.Frontend.Router.Echo()
	config.FromFile().Sys.SetRealIPParams(IRegister().Echo().RealIPConfig())
	e.SetRenderDataWrapper(xMW.DefaultRenderDataWrapper)
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
		if len(config.FromCLI().BackendDomain) == 0 &&
			len(Prefix()) == 0 && len(backend.Prefix()) == 0 {
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
	// Prometheus
	xmetrics.Register(e)

	// 注册模板引擎
	if httpserver.Frontend.Renderer() != nil {
		httpserver.Frontend.Renderer().Close()
	}
	if echo.String(`LABEL`) != `dev` && config.FromFile().Extend.GetStore(`minify`).Bool(`on`) {
		httpserver.Frontend.SetTmplCustomParser(TmplCustomParser)
	} else {
		httpserver.Frontend.SetTmplCustomParser(nil)
	}
	keepExtensionPrefixes := []string{httpserver.Frontend.StaticRootURLPath}
	if config.IsInstalled() {
		ctx := echo.NewContext(mock.NewRequest(), mock.NewResponse(), e)
		routeM := official.NewRoutePage(ctx)
		routes, _ := routeM.ListWithExtensionRoutes(RouteDefaultExtension)
		keepExtensionPrefixes = append(keepExtensionPrefixes, routes...)
	}
	httpserver.Frontend.SetKeepExtensionPrefixes(keepExtensionPrefixes)
	httpserver.Frontend.Apply()
	echo.OnCallback(`nging.renderer.cache.clear`, func(_ echo.Event) error {
		log.Debug(`clear: Frontend Template Object Cache`)
		renderOptions.Renderer().ClearCache()
		return nil
	}, `clear-frontend-template-object-cache`)
	echo.OnCallback(`webx.renderer.cache.clear`, func(_ echo.Event) error {
		log.Debug(`clear: Frontend Template Object Cache`)
		httpserver.Frontend.Renderer().ClearCache()
		return nil
	}, `clear-frontend-template-object-cache`)
	echo.OnCallback(`webx.frontend.close`, func(_ echo.Event) error {
		if httpserver.Frontend.Renderer() != nil {
			log.Debug(`close: Frontend Template Manager`)
			httpserver.Frontend.Renderer().Close()
		}
		return nil
	}, `close-frontend-template-manager`)
	e.Use(xMW.UseTheme(httpserver.Frontend.TmplPathFixers.ThemeInfo))
	e.Use(FrontendURLFuncMW())
	e.Use(ngingMW.FuncMap())
	e.Use(xMW.FuncMap())
	e.Use(render.Auto())
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

	httpserver.Frontend.I18n().Handler(e, `App.i18n`)
}

func URLFor(purl string) string {
	return xMW.URLFor(purl)
}

func Customer(c echo.Context) *dbschema.OfficialCustomer {
	return xMW.Customer(c)
}
