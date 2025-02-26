package middleware

import (
	"strings"

	"github.com/webx-top/com"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"

	"github.com/coscms/webcore/library/captcha/captchabiz"
	"github.com/coscms/webcore/library/config"
	"github.com/coscms/webcore/library/httpserver"
	"github.com/coscms/webcore/library/nerrors"
	cfgIPFilter "github.com/coscms/webfront/library/ipfilter"
	mwRateLimiter "github.com/webx-top/echo/middleware/ratelimiter"
)

func getRateLimiterConfig() (*cfgIPFilter.RateLimiterConfig, bool) {
	opts, ok := config.Setting(`frequency`).Get(`rateLimiter`).(*cfgIPFilter.RateLimiterConfig)
	return opts, ok
}

func RateLimiter() echo.MiddlewareFunc {
	rateLimiterConfig := &mwRateLimiter.RateLimiterConfig{
		Skipper: func(c echo.Context) bool {
			opts, ok := getRateLimiterConfig()
			if !ok || !opts.On {
				return true
			}
			return false
		},
	}
	if opts, ok := getRateLimiterConfig(); ok {
		opts.Apply(rateLimiterConfig)
	}
	return mwRateLimiter.RateLimiterWithConfig(*rateLimiterConfig)
}

func underAttackSkipper(c echo.Context) bool {
	switch c.Path() {
	case `/captcha/*`, `/captchago/:driver/:type`:
		return true
	}
	underAttack, ok := config.Setting(`frequency`).Get(`underAttack`).(string)
	if !ok {
		return true
	}
	return underAttack != `1`
}

func UnderAttack() echo.MiddlewareFunc {
	return func(h echo.Handler) echo.Handler {
		return echo.HandlerFunc(func(c echo.Context) error {
			if underAttackSkipper(c) {
				return h.Handle(c)
			}
			cookieValue := c.Cookie().DecryptGet(`CaptVerified`)
			if len(cookieValue) > 0 {
				parts := strings.SplitN(cookieValue, `|`, 2)
				if len(parts) != 2 {
					cookieValue = ``
				} else if parts[0] != c.RealIP() && parts[1] != com.Md5(c.Request().UserAgent()) {
					cookieValue = ``
				}
			}
			if len(cookieValue) > 0 {
				return h.Handle(c)
			}
			if c.IsPost() {
				data := captchabiz.VerifyCaptcha(c, httpserver.KindFrontend, `code`)
				if nerrors.IsFailureCode(data.GetCode()) {
					err := c.NewError(code.InvalidParameter, `验证码不正确`).SetZone(`code`)
					if c.Format() == `json` {
						return c.JSON(data)
					}
					return err
				}
				cookieValue = c.RealIP() + `|` + com.Md5(c.Request().UserAgent())
				c.Cookie().EncryptSet(`CaptVerified`, cookieValue, 86400)
				return c.Redirect(c.FullRequestURI())
			}
			return c.Render(`under_attack`, nil)
		})
	}
}
