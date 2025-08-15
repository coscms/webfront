package frontend

import (
	"sync"

	"github.com/coscms/webfront/dbschema"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/defaults"
	"github.com/webx-top/echo/middleware"
)

var rewriteConfig = &rewriteWithLock{
	cfg: &middleware.RewriteConfig{
		Rules: map[string]string{},
	},
}

type rewriteWithLock struct {
	cfg *middleware.RewriteConfig
	mu  sync.RWMutex
}

func (r *rewriteWithLock) SetConfig(cfg *middleware.RewriteConfig) {
	r.mu.Lock()
	r.cfg = cfg
	r.mu.Unlock()
}

func (r *rewriteWithLock) Config() *middleware.RewriteConfig {
	r.mu.RLock()
	cfg := r.cfg
	r.mu.RUnlock()
	return cfg
}

func (r *rewriteWithLock) Rewrite(urlPath string) string {
	return r.Config().Rewrite(urlPath)
}

func (r *rewriteWithLock) Reverse(urlPath string) string {
	return r.Config().Reverse(urlPath)
}

func (r *rewriteWithLock) UnrewriteMiddleware() echo.MiddlewareFuncd {
	return func(next echo.Handler) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			if r.Config().Skipper(c) {
				return next.Handle(c)
			}

			req := c.Request()
			rew := r.Config().Reverse(req.URL().Path())
			if rew != req.URL().Path() {
				req.URL().SetPath(rew)
			}
			return next.Handle(c)
		}
	}
}

func applyRouteRewrite(e *echo.Echo) error {
	err := SetRewriter(e)
	e.Pre(rewriteConfig.UnrewriteMiddleware)
	return err
}

func SetRewriter(e *echo.Echo) error {
	cfg, err := MakeRouteRewriter()
	if err != nil {
		return err
	}
	rewriteConfig.SetConfig(&cfg)
	e.SetRewriter(rewriteConfig)
	return err
}

func MakeRouteRewriter() (cfg middleware.RewriteConfig, err error) {
	cond := db.NewCompounds()
	cond.AddKV(`disabled`, `N`)
	f := dbschema.NewOfficialCommonRouteRewrite(defaults.NewMockContext())
	_, err = f.ListByOffset(nil, nil, 0, -1, cond.And())
	if err != nil {
		return
	}
	cfg.Skipper = echo.DefaultSkipper
	cfg.Rules = map[string]string{}
	for _, v := range f.Objects() {
		cfg.Rules[v.Route] = v.RewriteTo
	}
	cfg.Init()
	return
}

func ResetRouteRewrite() error {
	cfg, err := MakeRouteRewriter()
	if err != nil {
		return err
	}
	rewriteConfig.SetConfig(&cfg)
	return err
}
