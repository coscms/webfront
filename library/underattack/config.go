package underattack

import (
	"strings"
	"sync"

	"github.com/coscms/webcore/library/ipfilter"
	"github.com/webx-top/echo"
)

func NewConfig() *Config {
	return &Config{}
}

type Config struct {
	On          bool
	IPWhitelist string
	filter      *ipfilter.IPFilter
	sg          sync.Once
}

func (c *Config) FromStore(r echo.H) *Config {
	c.On = r.Bool(`On`)
	c.IPWhitelist = strings.TrimSpace(r.String(`IPWhitelist`))
	return c
}

func (c *Config) IsAllowed(ctx echo.Context) bool {
	c.sg.Do(c.initFilter)
	return c.filter.IsAllowed(ctx.RealIP())
}

func (c *Config) initFilter() {
	c.filter = ipfilter.NewWithIP(``, c.IPWhitelist)
}
