package middleware

import (
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/defaults"

	"github.com/coscms/webcore/library/config"
	cfgIPFilter "github.com/coscms/webfront/library/ipfilter"
	"github.com/coscms/webfront/library/underattack"
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

func UnderAttack(maxAge int) echo.MiddlewareFunc {
	return underattack.Middleware(maxAge)
}
