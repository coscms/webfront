package oauthutils

import (
	"github.com/admpub/goth"
	"github.com/admpub/log"
	"github.com/webx-top/com"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/handler/oauth2"
	"github.com/webx-top/echo/subdomains"

	"github.com/coscms/oauth2s/client/goth/providers"
	dbschemaNging "github.com/coscms/webcore/dbschema"
	"github.com/coscms/webcore/library/config"
	"github.com/coscms/webcore/library/httpserver"
	oauthLibrary "github.com/coscms/webcore/library/oauth"
	"github.com/coscms/webcore/registry/settings"
	"github.com/coscms/webfront/library/xcommon"
)

var (
	defaultOAuth      *oauth2.OAuth
	SuccessHandler    interface{}
	BeginAuthHandler  echo.Handler
	AfterLoginSuccess []func(ctx echo.Context, ouser *goth.User) (end bool, err error)
)

func OnAfterLoginSuccess(hooks ...func(ctx echo.Context, ouser *goth.User) (end bool, err error)) {
	AfterLoginSuccess = append(AfterLoginSuccess, hooks...)
}

func FireAfterLoginSuccess(ctx echo.Context, ouser *goth.User) (end bool, err error) {
	for _, hook := range AfterLoginSuccess {
		end, err = hook(ctx, ouser)
		if err != nil || end {
			return
		}
	}
	return
}

func Default() *oauth2.OAuth {
	return defaultOAuth
}

// InitOauth 第三方登录
func InitOauth(e *echo.Echo, middlewares ...interface{}) {
	if config.IsInstalled() {
		settings.Init(nil)
	}
	host := xcommon.SiteURL(nil)
	if len(host) == 0 {
		host = subdomains.Default.URL(``, httpserver.KindFrontend)
	}
	oauth2Config := oauth2.NewConfig()
	RegisterProvider(oauth2Config)

	if config.IsInstalled() {
		if oauthAccounts, err := FindAccounts(); err != nil {
			log.Error(err)
		} else {
			oauth2Config.AddAccount(oauthAccounts...)
		}
	}

	defaultOAuth = oauth2.New(host, oauth2Config)
	defaultOAuth.SetSuccessHandler(SuccessHandler)
	defaultOAuth.SetBeginAuthHandler(BeginAuthHandler)
	defaultOAuth.Wrapper(e, middlewares...)
}

func Accounts() []oauth2.Account {
	var accounts []oauth2.Account
	if Default() == nil {
		return accounts
	}
	Default().Config.RangeAccounts(func(a *oauth2.Account) bool {
		account := *a
		accounts = append(accounts, account)
		return true
	})
	return accounts
}

// FindAccounts 第三方登录平台账号
func FindAccounts() ([]*oauth2.Account, error) {
	m := &dbschemaNging.NgingConfig{}
	_, err := m.ListByOffset(nil, nil, 0, -1, db.Cond{`group`: `oauth`})
	var result []*oauth2.Account
	isProduction := config.FromFile().Sys.IsEnv(`prod`)
	for _, row := range m.Objects() {
		if len(row.Value) == 0 {
			continue
		}
		cfg := &oauthLibrary.Config{Name: row.Key, On: row.Disabled != `Y`}
		err = com.JSONDecode([]byte(row.Value), cfg)
		if err != nil {
			return result, err
		}
		value := cfg.ToAccount(row.Key)
		var provider func(account *oauth2.Account) goth.Provider
		if !isProduction {
			provider = providers.Get(value.Name + `_dev`)
		}
		if provider == nil {
			provider = providers.Get(value.Name)
		}
		value.SetConstructor(provider)
		if value.On {
			result = append(result, value)
		}
	}
	return result, err
}

// UpdateAccount 第三方登录平台账号
func UpdateAccount() error {
	if Default() == nil {
		return nil
	}
	accounts, err := FindAccounts()
	if err != nil {
		return err
	}
	Default().Config.ClearAccounts()
	Default().Config.AddAccount(accounts...)
	Default().Config.GenerateProviders()
	log.Debug(`update oauth configuration information`)
	return nil
}
