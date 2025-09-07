package underattack

import (
	"sort"
	"strings"

	"github.com/admpub/once"
	"github.com/coscms/webcore/library/ipfilter"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
)

func NewConfig() *Config {
	return &Config{}
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
		for _, rule := range c.Rules {
			if ck {
				if !strings.HasPrefix(dp, rule.Path) && !strings.HasPrefix(ctx.Path(), rule.Path) {
					continue
				}
			} else if !strings.HasPrefix(dp, rule.Path) {
				continue
			}
			if !rule.IsAllowed(ctx) {
				return c.filter.IsAllowed(ctx.RealIP())
			}
			return true
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
