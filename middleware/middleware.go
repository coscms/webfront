package middleware

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/webx-top/echo"
	stdCode "github.com/webx-top/echo/code"
	"github.com/webx-top/echo/subdomains"

	"github.com/admpub/log"
	"github.com/coscms/webcore/library/common"
	"github.com/coscms/webcore/library/config"
	"github.com/coscms/webcore/library/httpserver"
	"github.com/coscms/webcore/library/license"
	nav "github.com/coscms/webcore/library/navigate"
	"github.com/coscms/webcore/library/nerrors"
	uploadLibrary "github.com/coscms/webcore/library/upload"
	uploadClient "github.com/webx-top/client/upload"

	"github.com/coscms/webfront/dbschema"
	"github.com/coscms/webfront/initialize/frontend/usernav"
	"github.com/coscms/webfront/library/cache"
	"github.com/coscms/webfront/library/oauth2client"
	"github.com/coscms/webfront/library/xconst"
	"github.com/coscms/webfront/library/xrole"
	"github.com/coscms/webfront/library/xrole/xroleutils"
	"github.com/coscms/webfront/middleware/sessdata"
	modelCustomer "github.com/coscms/webfront/model/official/customer"
)

func SessionInfo(h echo.Handler) echo.HandlerFunc {
	return func(c echo.Context) error {
		ppath := c.Request().URL().Path()
		switch ppath {
		case c.Echo().Prefix() + `/favicon.ico`:
			return h.Handle(c)
		default:
			if strings.HasPrefix(ppath, c.Echo().Prefix()+uploadLibrary.UploadURLPath) {
				return h.Handle(c)
			}
		}
		baseCfg := config.Setting(`base`)
		siteClose := baseCfg.Uint(`siteClose`)
		if siteClose == xconst.SiteClose {
			return close(c, baseCfg)
		}
		var customer *dbschema.OfficialCustomer
		enabledJWT := c.Internal().Bool(`enabledJWT`)
		if enabledJWT {
			customerM := modelCustomer.NewCustomer(c)
			var err error
			customer, err = customerM.GetByJWT()
			if err != nil {
				log.Debug(err.Error())
			}
		}
		if customer == nil {
			customer, _ = c.Session().Get(`customer`).(*dbschema.OfficialCustomer)
		}
		if customer != nil {
			if siteClose == xconst.SiteOnlyAdmin && customer.Uid < 1 {
				return close(c, baseCfg)
			}
			c.Internal().Set(`customer`, customer)
		} else {
			if siteClose == xconst.SiteOnlyMember || siteClose == xconst.SiteOnlyAdmin {
				handlerPermission := c.Route().String(httpserver.MetaKeyPermission)
				if handlerPermission != httpserver.PermissionGuest {
					return goToSignIn(c)
				}
			}
			ouser, exists, err := oauth2client.GetSession(c)
			//echo.Dump(echo.H{`err`: err, `ouser`: ouser})
			if err == nil {
				c.Set(`oauth`, ouser) // 表单数据
			} else {
				if exists {
					log.Debug(err.Error())
				}
			}
		}
		return h.Handle(c)
	}
}

func close(c echo.Context, baseCfg echo.H) error {
	showAnnouncement := baseCfg.Bool(`showAnnouncement`)
	if !showAnnouncement {
		return c.Render(`error/under-maintenance`, nil, http.StatusServiceUnavailable)
		//return c.Render(`error/not-found`, nil, http.StatusNotFound)
	}
	siteAnnouncement := baseCfg.String(`siteAnnouncement`)
	siteAnnouncement = strings.TrimSpace(siteAnnouncement)
	data := echo.H{
		`title`:         `Coming Soon`,
		`content`:       siteAnnouncement,
		`enabledNotify`: false, //是否支持访客接收邮件通知(未实现)
	}
	if strings.HasPrefix(siteAnnouncement, `<h1>`) {
		siteAnnouncement = strings.TrimPrefix(siteAnnouncement, `<h1>`)
		splited := strings.SplitN(siteAnnouncement, `</h1>`, 2)
		switch len(splited) {
		case 2:
			data[`title`] = splited[0]
			data[`content`] = splited[1]
		}
	}
	return c.Render(`error/coming-soon`, data)
}

func userCenter(c echo.Context, customer *dbschema.OfficialCustomer) error {
	m := modelCustomer.NewCustomer(c)
	err := m.VerifySession(customer)
	if err != nil {
		if nerrors.IsUserNotLoggedIn(err) {
			common.SendErr(c, err)
			return goToSignIn(c)
		}
		return err
	}
	return permCheck(c, customer)
}

func permCheck(c echo.Context, customer *dbschema.OfficialCustomer) error {
	permission := xroleutils.CustomerPermission(c, customer)
	//echo.Dump(permission)
	customerID := fmt.Sprint(customer.Id)
	cacheTTL := xroleutils.CustomerPermTTL(c)
	ttlOpt := cache.GetTTLByNumber(cacheTTL, nil)
	c.SetFunc(`LeftNavigate`, func() nav.List {
		list := &nav.List{}
		cache.XFunc(c, sessdata.LeftNavigateCacheKey+customerID, list, func() error {
			*list = permission.FilterNavigate(c, usernav.LeftNavigate)
			return nil
		}, ttlOpt)
		return *list
	})
	c.SetFunc(`TopNavigate`, func() nav.List {
		list := &nav.List{}
		cache.XFunc(c, sessdata.TopNavigateCacheKey+customerID, list, func() error {
			*list = permission.FilterNavigate(c, usernav.TopNavigate)
			return nil
		}, ttlOpt)
		return *list
	})
	if !c.Internal().Bool(`skipCurrentURLPermCheck`) {
		rpath := c.Path()
		if len(c.Echo().Prefix()) > 0 {
			rpath = strings.TrimPrefix(rpath, c.Echo().Prefix())
		}
		if err := checkPermission(c, customer, permission, rpath); err != nil {
			return err
		}
	}
	c.SetFunc(`CheckPerm`, func(routePath string) error {
		return checkPermission(c, customer, permission, routePath)
	})
	return nil
}

func checkPermission(ctx echo.Context, customer *dbschema.OfficialCustomer, permission *xrole.RolePermission, routePath string) error {
	handlerPermission := ctx.Route().String(httpserver.MetaKeyPermission)
	if handlerPermission == httpserver.PermissionPublic {
		return nil
	}
	return xrole.CheckPermissionByRoutePath(ctx, customer, permission, routePath)
}

func SkipCurrentURLPermCheck(h echo.Handler) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Internal().Set(`skipCurrentURLPermCheck`, true)
		return h.Handle(c)
	}
}

func AuthCheck(h echo.Handler) echo.HandlerFunc {
	return func(c echo.Context) error {
		//检查是否已安装
		if !config.IsInstalled() {
			c.Data().SetError(c.NewError(stdCode.SystemNotInstalled, c.T(`请先安装`)))
			return c.Redirect(subdomains.Default.URL(`/setup`, `backend`))
		}

		//验证授权文件
		if !license.Ok(c) {
			c.Data().SetError(c.NewError(stdCode.SystemUnauthorized, c.T(`请先获取本系统授权`)))
			return c.Redirect(subdomains.Default.URL(`/license`, `backend`))
		}

		customer := Customer(c)
		if customer != nil {
			if err := userCenter(c, customer); err != nil {
				return err
			}
			return h.Handle(c)
		}
		return goToSignIn(c)
	}
}

func TrimPathSuffix(ignorePrefixes ...string) echo.MiddlewareFuncd {
	return func(h echo.Handler) echo.HandlerFunc {
		return func(c echo.Context) error {
			upath := c.Request().URL().Path()
			for _, ignorePrefix := range ignorePrefixes {
				if strings.HasPrefix(upath, ignorePrefix) {
					return h.Handle(c)
				}
			}
			cleanedPath := strings.TrimSuffix(upath, c.DefaultExtension())
			c.Request().URL().SetPath(cleanedPath)
			return h.Handle(c)
		}
	}
}

func goToSignIn(c echo.Context) error {
	var queryString string
	var next string
	if c.IsGet() {
		next = c.Request().URI()
	} else if c.IsPost() {
		if c.Format() == echo.ContentTypeJSON {
			client := c.Form(`client`)
			if len(client) > 0 {
				cli := uploadClient.Get(client)
				cli.Init(c, &uploadClient.Result{})
				cli.SetError(c.NewError(stdCode.Unauthenticated, c.T(`请先登录`)))
				return cli.Response()
			}
		}
		next = c.Referer()
	}
	if len(next) > 0 && !strings.Contains(next, `/sign_in`) {
		queryString = `?next=` + url.QueryEscape(next)
	}
	c.Data().SetError(c.NewError(stdCode.Unauthenticated, c.T(`请先登录`)))
	return c.Redirect(URLFor(`/sign_in`) + queryString)
}

var (
	Customer = sessdata.Customer
	URLFor   = sessdata.URLFor
)

var Middlewares []interface{}

func Use(m ...interface{}) {
	Middlewares = append(Middlewares, m...)
}
