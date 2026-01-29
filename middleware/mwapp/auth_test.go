package mwapp

import (
	"net/url"
	"testing"
	"time"

	"github.com/admpub/log"
	"github.com/stretchr/testify/assert"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/defaults"
	"github.com/webx-top/echo/param"
)

func TestAuth(t *testing.T) {
	defer log.Close()
	cfg := DefaultAuthAppConfig
	cfg.SetSecretGetter(func(ctx echo.Context, appID string) (string, error) {
		return `secret`, nil
	})
	cfg.SetDefaults()
	ctx := defaults.NewMockContext()
	ctx.Request().Form().Set(cfg.FormAppIDKey, `test`)
	nowTs := param.AsString(time.Now().Unix())
	sign := cfg.SignMaker()(url.Values{
		cfg.FormAppIDKey: []string{`test`},
		cfg.FormTimeKey:  []string{nowTs},
	}, `secret`)
	ctx.Request().Form().Set(cfg.FormSignKey, sign)
	ctx.Request().Form().Set(cfg.FormTimeKey, nowTs)
	err := cfg.Verify(ctx)
	assert.NoError(t, err)
}
