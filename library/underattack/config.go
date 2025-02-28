package underattack

import (
	"sync"

	"github.com/coscms/webcore/library/ipfilter"
	"github.com/webx-top/echo"
)

type Config struct {
	On          bool
	IPWhitelist string
	filter      *ipfilter.IPFilter
	sg          sync.Once
}

func (c *Config) IsAllowed(ctx echo.Context) bool {
	c.sg.Do(c.initFilter)
	return c.filter.IsAllowed(ctx.RealIP())
}

func (c *Config) initFilter() {
	c.filter = ipfilter.NewWithIP(``, c.IPWhitelist)
}
