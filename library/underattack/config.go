package underattack

import (
	"regexp"
	"strings"

	"github.com/admpub/once"
	"github.com/coscms/webcore/library/ipfilter"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"
	"github.com/webx-top/echo/param"
)

func NewConfig() *Config {
	return &Config{}
}

type Config struct {
	On                bool
	IPWhitelist       string
	UAWhitelist       string
	HeaderName        string
	HeaderValue       string
	URIPathWhitelist  []string
	URIQueryWhitelist []string
	filter            *ipfilter.IPFilter
	regexpUA          *regexp.Regexp
	kvURIQuery        map[string][]string
	once              once.Once
}

func (c *Config) FromStore(r echo.H) *Config {
	c.On = r.Bool(`On`)
	c.IPWhitelist = strings.TrimSpace(r.String(`IPWhitelist`))
	c.UAWhitelist = strings.TrimSpace(r.String(`UAWhitelist`))
	c.HeaderName = strings.TrimSpace(r.String(`HeaderName`))
	c.HeaderValue = strings.TrimSpace(r.String(`HeaderValue`))
	c.URIPathWhitelist = param.Unique(com.TrimSpaceForRows(strings.TrimSpace(r.String(`URIPathWhitelist`))))
	c.URIQueryWhitelist = param.Unique(com.TrimSpaceForRows(strings.TrimSpace(r.String(`URIQueryWhitelist`))))
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
	if len(c.URIPathWhitelist) > 0 {
		if com.InSlice(ctx.Path(), c.URIPathWhitelist) {
			return true
		}
		if ctx.Path() != ctx.DispatchPath() && com.InSlice(ctx.DispatchPath(), c.URIPathWhitelist) {
			return true
		}
	}
	c.once.Do(c.initFilter)
	if c.kvURIQuery != nil {
		for key, values := range c.kvURIQuery {
			inputs := ctx.FormValues(key)
			for _, value := range values {
				if com.InSlice(value, inputs) {
					return true
				}
			}
		}
	}
	if c.regexpUA != nil {
		if c.regexpUA.MatchString(ctx.Request().UserAgent()) {
			return true
		}
	}
	return c.filter.IsAllowed(ctx.RealIP())
}

func (c *Config) initFilter() {
	c.filter = ipfilter.NewWithIP(``, c.IPWhitelist).SetDisallow(true)
	if len(c.UAWhitelist) > 0 {
		rows := param.Unique(com.TrimSpaceForRows(c.UAWhitelist))
		if len(rows) > 0 {
			c.regexpUA = regexp.MustCompile(strings.Join(rows, `|`))
		}
	}
	c.kvURIQuery = map[string][]string{}
	for _, row := range c.URIQueryWhitelist {
		parts := strings.SplitN(row, `=`, 2)
		for k, v := range parts {
			parts[k] = strings.TrimSpace(v)
		}
		if len(parts[0]) == 0 {
			continue
		}
		if _, ok := c.kvURIQuery[parts[0]]; !ok {
			c.kvURIQuery[parts[0]] = []string{}
		}
		c.kvURIQuery[parts[0]] = append(c.kvURIQuery[parts[0]], parts[1])
	}
}
