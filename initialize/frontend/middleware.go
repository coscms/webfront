package frontend

import (
	"github.com/coscms/webcore/library/ntemplate"
	"github.com/coscms/webfront/library/xtheme"
	"github.com/webx-top/echo"
)

func FrontendURLFuncMW() echo.MiddlewareFunc {
	return func(h echo.Handler) echo.Handler {
		return echo.HandlerFunc(func(c echo.Context) error {
			FrontendURLFunc(c)
			return h.Handle(c)
		})
	}
}

func FrontendURLFunc(c echo.Context) error {
	var themeInfo *ntemplate.ThemeInfo
	c.SetFunc(`ThemeInfo`, func() *ntemplate.ThemeInfo {
		if themeInfo != nil {
			return themeInfo
		}
		themeInfo = xtheme.GetThemeInfo(c)
		return themeInfo
	})
	return nil
}
