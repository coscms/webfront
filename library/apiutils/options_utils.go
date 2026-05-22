package apiutils

import (
	"strings"

	"github.com/coscms/webcore/library/config"
	webxdbschema "github.com/coscms/webfront/dbschema"
	"github.com/coscms/webfront/library/xcommon"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"
)

func (o *Options) getAccount(cond db.Compound) (err error) {
	accountM := webxdbschema.NewOfficialCommonApiAccount(o.ctx)
	err = accountM.Get(nil, cond)
	if err != nil {
		if err == db.ErrNoMoreRows {
			err = o.ctx.NewError(code.DataNotFound, `没找到API账号数据：%s`, echo.Dump(cond, false))
		}
		return
	}
	gt := AppInfoFromAccount{
		Account: accountM,
	}
	if config.FromFile().Sys.IsEnv(`prod`) {
		gt.URLPrefix = accountM.Url
	} else {
		gt.URLPrefix = accountM.UrlDev
	}
	gt.URLPrefix = strings.TrimSuffix(gt.URLPrefix, `/`)
	o.Options.SetAppInfoGetter(gt)
	return
}

var AppInfoDefaultGetter = func(ctx echo.Context, cond db.Compound) (appInfo AppInfo, err error) {
	err = ctx.NewError(code.Unsupported, `尚未设置App信息获取方式`)
	return
}

func (o *Options) getApp(cond db.Cond) (err error) {
	var appInfo AppInfo
	appInfo, err = o.onlyGetApp(cond)
	if err != nil {
		return
	}
	gt := &AppInfoFromOpenApp{
		App:       appInfo,
		URLPrefix: xcommon.SiteURL(o.ctx),
	}
	o.Options.SetAppInfoGetter(gt)
	return
}

func (o *Options) onlyGetApp(cond db.Cond) (appInfo AppInfo, err error) {
	if o._appInfoGetter == nil {
		appInfo, err = AppInfoDefaultGetter(o.ctx, cond)
	} else {
		appInfo, err = o._appInfoGetter(o.ctx, cond)
	}
	return
}
