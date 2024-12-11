package qrsignin

import (
	"time"

	cached "github.com/admpub/cache"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"

	"github.com/coscms/webfront/library/cache"
)

func newCacheQRSignIn() cacheQRSignIn {
	c := cacheQRSignIn{}
	c.SetDefaults()
	return c
}

type cacheQRSignIn struct {
	Prefix    string
	KeyPrefix string
}

func (c *cacheQRSignIn) SetDefaults() {
	if len(c.Prefix) == 0 {
		c.Prefix = `QRSignIn_`
	}
	if len(c.KeyPrefix) == 0 {
		c.KeyPrefix = `QRSIKey_`
	}
}

func (c cacheQRSignIn) Encode(ctx echo.Context, signInData QRSignIn) (string, error) {
	qrkeysKey := c.KeyPrefix + ctx.Session().MustID()
	{
		var oldKey string
		if err := cache.Get(ctx, qrkeysKey, &oldKey); err == nil && len(oldKey) > 0 {
			cache.Delete(ctx, c.Prefix+oldKey)
		}
	}
	key := GenerateUniqueKey(ctx.RealIP(), ctx.Request().UserAgent())
	timeout := signInData.Expires - time.Now().Unix()
	err := cache.Put(ctx, c.Prefix+key, signInData, timeout)
	if err == nil {
		err = cache.Put(ctx, qrkeysKey, key, timeout)
	}
	return key, err
}

func (c cacheQRSignIn) Decode(ctx echo.Context, key string) (QRSignIn, error) {
	signInData := QRSignIn{}
	if !com.StrIsAlphaNumeric(key) {
		return signInData, ctx.NewError(code.InvalidParameter, `二维码包含无效字符`).SetZone(`data`)
	}
	err := cache.Get(ctx, c.Prefix+key, &signInData)
	if err == nil {
		if err := cache.Delete(ctx, c.Prefix+key); err == nil {
			cache.Delete(ctx, c.KeyPrefix+ctx.Session().MustID())
		}
	} else if err == cached.ErrNotFound {
		return signInData, ctx.NewError(code.DataHasExpired, `二维码已经失效`).SetZone(`data`)
	} else if err == cached.ErrExpired {
		return signInData, ctx.NewError(code.DataHasExpired, `二维码已经过期`).SetZone(`data`)
	}
	return signInData, err
}
