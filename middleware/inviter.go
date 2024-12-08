package middleware

import (
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
)

func Inviter(h echo.Handler) echo.HandlerFunc {
	return func(c echo.Context) error {
		inviterID := c.Formx(`inviter`).Uint64()
		if inviterID > 0 {
			c.SetCookie(`inviter`, com.String(inviterID))
		}
		return h.Handle(c)
	}
}
