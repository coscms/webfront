package underattack

import (
	"net/http"
	"strings"
	"time"

	"github.com/webx-top/com"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"
	"github.com/webx-top/echo/defaults"
	"github.com/webx-top/echo/param"

	"github.com/coscms/webcore/library/captcha/captchabiz"
	"github.com/coscms/webcore/library/config"
	"github.com/coscms/webcore/library/httpserver"
	"github.com/coscms/webcore/library/nerrors"
	"github.com/coscms/webfront/dbschema"
)

func underAttackSkipper(c echo.Context) (*Config, bool) {
	if defaults.IsMockContext(c) || c.Route().Bool(`noAttack`) {
		return nil, true
	}
	customer, ok := c.Session().Get(`customer`).(*dbschema.OfficialCustomer)
	if ok && customer != nil {
		return nil, true
	}
	cfg := config.Setting(`frequency`).Get(`underAttack`)
	switch v := cfg.(type) {
	case *Config:
		return v, !v.On || v.IsAllowed(c)
	case string:
		return nil, v != `1`
	default:
		return nil, true
	}
}

func Middleware(maxAge int) echo.MiddlewareFunc {
	partsLen := 4
	return func(h echo.Handler) echo.Handler {
		return echo.HandlerFunc(func(c echo.Context) error {
			cfg, skip := underAttackSkipper(c)
			if skip {
				return h.Handle(c)
			}
			if cookieValue := c.Cookie().DecryptGet(`CaptVerified`); len(cookieValue) > 0 {
				parts := strings.SplitN(cookieValue, `|`, partsLen)
				if len(parts) == partsLen {
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
					err := c.NewError(code.InvalidParameter, param.AsString(data.GetInfo())).SetZone(`code`)
					if c.Format() == echo.ContentTypeJSON {
						return c.JSON(data)
					}
					return err
				}
				duration := time.Second * time.Duration(maxAge)
				cookieValue := c.RealIP() + `|` + com.Md5(c.Request().UserAgent()) + `|` + com.String(time.Now().Add(duration).Unix()) + `|` + com.RandomAlphanumeric(6)
				c.Cookie().EncryptSet(`CaptVerified`, cookieValue, duration)
				return c.Redirect(c.FullRequestURI())
			}
			_, captchaType := captchabiz.GetCaptchaType()
			c.Set(`captchaType`, captchaType)
			if cfg != nil {
				c.Set(`message`, cfg.Message)
			}
			return c.Render(`under_attack`, nil, http.StatusForbidden) //http.StatusRequestTimeout
		})
	}
}
