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
	"github.com/webx-top/echo/middleware"
	"github.com/webx-top/echo/middleware/render"
	"github.com/webx-top/echo/subdomains"

	captchaLib "github.com/coscms/webcore/library/captcha"
	"github.com/coscms/webcore/library/captcha/captchabiz"
	"github.com/coscms/webcore/library/captcha/driver/captcha_go"

	"github.com/coscms/webcore/cmd/bootconfig"
	"github.com/coscms/webcore/initialize/backend"
	backendLib "github.com/coscms/webcore/library/backend"
	"github.com/coscms/webcore/library/config"
	"github.com/coscms/webcore/library/httpserver"
	"github.com/coscms/webcore/library/notice"
	ngingMW "github.com/coscms/webcore/middleware"
	"github.com/coscms/webfront/dbschema"
	"github.com/coscms/webfront/library/frontend"
	"github.com/coscms/webfront/library/routepage"
	"github.com/coscms/webfront/library/xmetrics"
	xMW "github.com/coscms/webfront/middleware"
	"github.com/coscms/webfront/model/official"

	// 执行 init()
	_ "github.com/coscms/webfront/library/formbuilder" // 引入表单模板文件
	_ "github.com/coscms/webfront/registry/route"      //初始化前台服务路由
)

const (
	DefaultTemplateDir    = `./template/frontend`      // 前台模板路径默认值
	DefaultAssetsDir      = `./public/assets/frontend` // 前台素材路径默认值
	DefaultAssetsURLPath  = `/public/assets/frontend`  // 前台素材网址路径默认值
	RouteDefaultExtension = `.html`                    // 前台网页扩展名默认值
)

var TmplCustomParser func(tmpl string, content []byte) []byte
var Notify = notice.NewUserNotices(false, nil)

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
	httpserver.Frontend.RouteDefaultExtension = RouteDefaultExtension
	httpserver.Frontend.HostCheckerRegexpKey = `frontend.hostRuleRegexp`
	httpserver.Frontend.FuncSetters = []func(echo.Context) error{
		FrontendURLFunc,
		xMW.SetFunc,
	}
	bootconfig.OnStart(-1, start)
}

func start() {
	// 如果指定了后台域名则只能用该域名访问后台，此时将其它域名指向前台；如果前台域名和后台域名均未设置则代表使用指定路径访问后台，默认指向前台
	if len(config.FromCLI().BackendDomain) > 0 || len(config.FromCLI().FrontendDomain) == 0 {
		subdomains.Default.Default = httpserver.KindFrontend // 设置默认(没有匹配到域名的时候)访问的服务别名
	}
	// 初始化前台服务
	InitWebServer()
}

func SetPrefix(prefix string) {
	httpserver.Frontend.SetPrefix(prefix)
}

func getFrontendDomain() string {
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
			len(Prefix()) == 0 && len(backend.Prefix()) == 0 { // 前后台都没有指定域名和路径前缀的时候，给前台强制指定一个域名
			frontendDomain = `localhost:` + fmt.Sprintf(`%d`, config.FromCLI().Port)
		}
	} else {
		var domains []string
		for _, domain := range backend.MakeSubdomains(frontendDomain, []string{}) {
			if _, ok := subdomains.Default.Hosts.GetOk(domain); !ok { // 排除后台指定的域名
				domains = append(domains, domain)
			}
		}

		frontendDomain = strings.Join(domains, `,`)
	}
	return frontendDomain
}

func InitWebServer() {
	e := httpserver.Frontend.Router.Echo()
	config.FromFile().Sys.SetRealIPParams(e.RealIPConfig())
	e.SetRenderDataWrapper(xMW.DefaultRenderDataWrapper)

	// 子域名设置
	frontendDomain := getFrontendDomain()
	subdomains.Default.Add(httpserver.KindFrontend+`@`+frontendDomain, e)
	log.Infof(`Registered host: %s@%s`, httpserver.KindFrontend, frontendDomain)

	// 前台服务设置
	addMiddleware(e)
	e.Get(`/favicon.ico`, faviconHandler).SetMetaKV(httpserver.PermGuestKV())
	e.Use(xMW.SessionInfo)
	if config.IsInstalled() {
		routepage.Apply(e, frontend.GlobalFuncMap())
		applyRouteRewrite(e)
	}
	Apply()
}

func addMiddleware(e *echo.Echo) {
	if bootconfig.Bindata {
		e.SetDebug(false)
	} else {
		e.SetDebug(true)
	}
	// Prometheus
	xmetrics.Register(e)

	backendStaticMW := httpserver.Backend.GetStaticMW()
	if backendStaticMW != nil {
		middlewares := []interface{}{}
		if !config.FromFile().Sys.DisableHTTPLog {
			middlewares = append(middlewares, middleware.Log())
		}
		middlewares = append(middlewares, backendStaticMW)
		httpserver.Frontend.Middlewares = middlewares
	}

	httpserver.Frontend.GlobalFuncMap = frontend.GlobalFuncMap()
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
		httpserver.Frontend.Renderer().ClearCache()
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
	e.Use(xMW.UseTheme(httpserver.Frontend.Template.ThemeInfo))
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

	// Inviter
	e.Use(xMW.Inviter)

	e.Use(xMW.Middlewares...)

	captcha.New(``).Wrapper(e).SetMetaKV(httpserver.PermGuestKV())
	captchaGoG := e.Group(`/captchago`, captchabiz.CheckEnable(captchaLib.TypeGo)).SetMetaKV(httpserver.PermGuestKV())
	captcha_go.RegisterRoute(captchaGoG)

	e.Route("GET", `/qrcode`, backendLib.QrCode).SetMetaKV(httpserver.PermGuestKV())

	i18nG := e.Group(`/i18n`).SetMetaKV(httpserver.PermGuestKV())
	httpserver.Frontend.I18n().Handler(i18nG, `App.i18n`)
}

func URLFor(purl string) string {
	return xMW.URLFor(purl)
}

func Customer(c echo.Context) *dbschema.OfficialCustomer {
	return xMW.Customer(c)
}
