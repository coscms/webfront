package apiutils

import (
	"github.com/coscms/sdk/sdk_options"
	"github.com/coscms/webfront/dbschema"
)

var _ sdk_options.AppInfoGetter = (*AppInfoFromAccount)(nil)
var _ sdk_options.AppInfoGetter = (*AppInfoFromOpenApp)(nil)

type AppInfoFromAccount struct {
	Account   *dbschema.OfficialCommonApiAccount
	URLPrefix string
}

func (a AppInfoFromAccount) GetAppSecret() string {
	return a.Account.AppSecret
}

func (a AppInfoFromAccount) GetAppID() string {
	return a.Account.AppId
}

func (a AppInfoFromAccount) GetApiEndpoint() string {
	return a.URLPrefix
}

type AppInfoFromOpenApp struct {
	App       AppInfo
	URLPrefix string
}

func (a AppInfoFromOpenApp) GetAppSecret() string {
	return a.App.GetAppSecret()
}

func (a AppInfoFromOpenApp) GetAppID() string {
	return a.App.GetAppId()
}

func (a AppInfoFromOpenApp) GetApiEndpoint() string {
	return a.URLPrefix
}

type AppInfo interface {
	GetAppSecret() string
	GetAppId() string
	IsOfficial() bool
}
