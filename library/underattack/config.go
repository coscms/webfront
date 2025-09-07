package underattack

import (
	"net/url"
	"regexp"
	"sort"
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

type Rule struct {
	Path        string
	UAWhitelist string
	Headers     map[string]string
	QueryString string // a=b&b=c
	regexpUA    *regexp.Regexp
	kvURIQuery  map[string][]string
}

func (r *Rule) init() {
	if len(r.UAWhitelist) > 0 {
		rows := param.Unique(com.TrimSpaceForRows(r.UAWhitelist))
		if len(rows) > 0 {
			r.regexpUA = regexp.MustCompile(strings.Join(rows, `|`))
		}
	}
	r.kvURIQuery = map[string][]string{}
	if len(r.QueryString) > 0 {
		r.kvURIQuery, _ = url.ParseQuery(r.QueryString)
	}
}

func (r Rule) Validate(ctx echo.Context) error {
	var err error
	if len(r.UAWhitelist) > 0 {
		rows := com.TrimSpaceForRows(r.UAWhitelist)
		for _, row := range rows {
			_, err = regexp.Compile(row)
			if err != nil {
				return ctx.NewError(code.InvalidParameter, `正则表达式语法错误: %s`, row)
			}
		}
	}
	return err
}

func (r Rule) IsAllowed(ctx echo.Context) bool {
	var dflt bool
	if len(r.Headers) > 0 {
		dflt = true
		for k, v := range r.Headers {
			if ctx.Header(k) != v {
				return false
			}
		}
	}
	if r.kvURIQuery != nil {
		dflt = true
		for key, values := range r.kvURIQuery {
			inputs := ctx.FormValues(key)
			for _, value := range values {
				if !com.InSlice(value, inputs) {
					return false
				}
			}
		}
	}
	if r.regexpUA != nil {
		dflt = true
		if !r.regexpUA.MatchString(ctx.Request().UserAgent()) {
			return false
		}
	}
	return dflt
}

type Config struct {
	On          bool
	IPWhitelist string
	Rules       []*Rule // Rules[0][Path]
	filter      *ipfilter.IPFilter
	once        once.Once
}

func (c *Config) FromStore(r echo.H) *Config {
	c.On = r.Bool(`On`)
	c.IPWhitelist = strings.TrimSpace(r.String(`IPWhitelist`))
	m := r.GetStore(`Rules`)
	pathList := getSliceString(m.Get(`Path`))
	uaList := getSliceString(m.Get(`UAWhitelist`))
	uaLength := len(uaList)
	hdList := getSliceString(m.Get(`Headers`))
	hdLength := len(hdList)
	qsList := getSliceString(m.Get(`QueryString`))
	qsLength := len(qsList)
	for index, ppath := range pathList {
		if uaLength > index && hdLength > index && qsLength > index {
			ppath := strings.TrimSpace(ppath)
			if len(ppath) == 0 {
				continue
			}
			if !strings.HasPrefix(ppath, `/`) {
				ppath = `/` + ppath
			}
			c.Rules = append(c.Rules, &Rule{
				Path:        ppath,
				UAWhitelist: strings.TrimSpace(uaList[index]),
				Headers:     com.SplitKVRows(hdList[index]),
				QueryString: strings.TrimSpace(qsList[index]),
			})
			continue
		}
		break
	}
	sort.Slice(c.Rules, func(i, j int) bool {
		return len(c.Rules[i].Path) > len(c.Rules[j].Path)
	})
	return c
}

func getSliceString(v interface{}) []string {
	value, ok := v.([]interface{})
	if !ok {
		return nil
	}
	var ss []string
	for _, val := range value {
		if ss, ok = val.([]string); ok {
			return ss
		}
	}
	return ss
}

func (c *Config) Validate(ctx echo.Context) error {
	err := ipfilter.ValidateRows(ctx, c.IPWhitelist)
	if err != nil {
		return err
	}
	for _, rule := range c.Rules {
		if err = rule.Validate(ctx); err != nil {
			return err
		}
	}
	return err
}

func (c *Config) IsAllowed(ctx echo.Context) bool {
	c.once.Do(c.initFilter)
	if len(c.Rules) > 0 {
		dp := ctx.DispatchPath()
		ck := dp != ctx.Path()
		var passed bool
		for _, rule := range c.Rules {
			if ck {
				if !strings.HasPrefix(dp, rule.Path) && !strings.HasPrefix(ctx.Path(), rule.Path) {
					continue
				}
			} else if !strings.HasPrefix(dp, rule.Path) {
				continue
			}
			if !rule.IsAllowed(ctx) {
				return false
			}
			passed = true
		}
		if passed {
			return passed
		}
	}
	return c.filter.IsAllowed(ctx.RealIP())
}

func (c *Config) initFilter() {
	c.filter = ipfilter.NewWithIP(``, c.IPWhitelist).SetDisallow(true)
	for _, rule := range c.Rules {
		rule.init()
	}
}
