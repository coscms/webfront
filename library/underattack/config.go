package underattack

import (
	"regexp"
	"strings"
	"sync"

	"github.com/coscms/webcore/library/ipfilter"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"
)

func NewConfig() *Config {
	return &Config{}
}

type Config struct {
	On          bool
	IPWhitelist string
	UAWhitelist string
	HeaderName  string
	HeaderValue string
	filter      *ipfilter.IPFilter
	regexp      *regexp.Regexp
	sg          sync.Once
}

func (c *Config) FromStore(r echo.H) *Config {
	c.On = r.Bool(`On`)
	c.IPWhitelist = strings.TrimSpace(r.String(`IPWhitelist`))
	c.UAWhitelist = strings.TrimSpace(r.String(`UAWhitelist`))
	c.HeaderName = strings.TrimSpace(r.String(`HeaderName`))
	c.HeaderValue = strings.TrimSpace(r.String(`HeaderValue`))
	return c
}

func (c *Config) Validate(ctx echo.Context) error {
	err := ipfilter.ValidateRows(ctx, c.IPWhitelist)
	if err != nil {
		return err
	}
	if len(c.UAWhitelist) > 0 {
		rows := com.TrimSpaceForRows(c.UAWhitelist)
		for _, row := range rows {
			_, err = regexp.Compile(row)
			if err != nil {
				return ctx.NewError(code.InvalidParameter, `正则表达式语法错误: %s`, row)
			}
		}
	}
	return err
}

func (c *Config) IsAllowed(ctx echo.Context) bool {
	if len(c.HeaderName) > 0 {
		if len(c.HeaderValue) > 0 {
			if ctx.Header(c.HeaderName) == c.HeaderValue {
				return true
			}
		} else {
			if len(ctx.Request().Header().Values(c.HeaderName)) > 0 {
				return true
			}
		}
	}
	c.sg.Do(c.initFilter)
	if c.regexp != nil {
		if c.regexp.MatchString(ctx.Request().UserAgent()) {
			return true
		}
	}
	return c.filter.IsAllowed(ctx.RealIP())
}

func (c *Config) initFilter() {
	c.filter = ipfilter.NewWithIP(``, c.IPWhitelist)
	if len(c.UAWhitelist) > 0 {
		rows := com.TrimSpaceForRows(c.UAWhitelist)
		c.regexp = regexp.MustCompile(strings.Join(rows, `|`))
	}
}
