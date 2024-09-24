package frontend

import (
	"github.com/coscms/webcore/library/httpserver"
	"github.com/coscms/webcore/library/ntemplate"
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
	c.SetFunc(`ThemeInfo`, func() *ntemplate.ThemeInfo {
		return httpserver.Frontend.Template.ThemeInfo(c)
	})
	return nil
}
