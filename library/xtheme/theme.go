package xtheme

import (
	"github.com/coscms/webcore/library/httpserver"
	"github.com/coscms/webcore/library/ntemplate"
	"github.com/webx-top/echo"
)

func GetThemeInfo(ctx echo.Context) *ntemplate.ThemeInfo {
	themeInfo, _ := ctx.Internal().Get(ntemplate.InternalKeyThemeConfig).(*ntemplate.ThemeInfo)
	if themeInfo == nil {
		themeInfo = httpserver.Frontend.Template.ThemeInfo(ctx)
		ctx.Internal().Set(ntemplate.InternalKeyThemeConfig, themeInfo)
	}
	return themeInfo
}
