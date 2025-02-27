package middleware

import (
	"strings"
	"time"

	"github.com/webx-top/com"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"
	"github.com/webx-top/echo/defaults"

	"github.com/coscms/webcore/library/captcha/captchabiz"
	"github.com/coscms/webcore/library/config"
	"github.com/coscms/webcore/library/httpserver"
	"github.com/coscms/webcore/library/nerrors"
	"github.com/coscms/webfront/dbschema"
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
			if defaults.IsMockContext(c) {
				return true
			}
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
	if defaults.IsMockContext(c) || c.Route().Bool(`noAttack`) {
		return true
	}

	underAttack := config.Setting(`frequency`).String(`underAttack`, `0`)
	if underAttack != `1` {
		return true
	}
	customer, ok := c.Session().Get(`customer`).(*dbschema.OfficialCustomer)
	if ok && customer != nil {
		return true
	}
	return false
}

func UnderAttack(maxAge int) echo.MiddlewareFunc {
	return func(h echo.Handler) echo.Handler {
		return echo.HandlerFunc(func(c echo.Context) error {
			if underAttackSkipper(c) {
				return h.Handle(c)
			}
			if cookieValue := c.Cookie().DecryptGet(`CaptVerified`); len(cookieValue) > 0 {
				parts := strings.SplitN(cookieValue, `|`, 3)
				if len(parts) == 3 {
					unixtime := com.Int64(parts[2])
					passed := unixtime >= time.Now().Unix() && parts[0] == c.RealIP() && parts[1] == com.Md5(c.Request().UserAgent())
					if passed {
						return h.Handle(c)
					}
				}
			}
			if c.IsPost() {
				data := captchabiz.VerifyCaptcha(c, httpserver.KindFrontend, `code`)
				if nerrors.IsFailureCode(data.GetCode()) {
					err := c.NewError(code.InvalidParameter, `验证码不正确`).SetZone(`code`)
					if c.Format() == echo.ContentTypeJSON {
						return c.JSON(data)
					}
					return err
				}
				duration := time.Second * time.Duration(maxAge)
				cookieValue := c.RealIP() + `|` + com.Md5(c.Request().UserAgent()) + `|` + com.String(time.Now().Add(duration).Unix())
				c.Cookie().EncryptSet(`CaptVerified`, cookieValue, duration)
				return c.Redirect(c.FullRequestURI())
			}
			_, captchaType := captchabiz.GetCaptchaType()
			c.Set(`captchaType`, captchaType)
			return c.Render(`under_attack`, nil)
		})
	}
}
