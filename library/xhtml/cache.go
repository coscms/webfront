package xhtml

import (
	"context"
	"database/sql"
	"net/http"
	"strings"

	"github.com/webx-top/echo"
	"github.com/webx-top/echo/defaults"

	"github.com/coscms/webfront/library/cache"
	"github.com/coscms/webfront/middleware/sessdata"
)

func controlCache(ctx echo.Context, cacheKey string, urlWithQueryString ...bool) (bool, error) {
	nocacheStr := ctx.Form(`nocache`)
	if len(nocacheStr) == 0 {
		return true, nil
	}
	switch nocacheStr {
	case `1`: // 禁用HTML缓存和数据缓存
		return false, nil
	case `2`, `make`: // 强制缓存新HTML
		fallthrough
	case `3`, `mkall`:
		isMakeAll := nocacheStr == `mkall`
		reqURL := ctx.Request().URL().Path()
		if len(urlWithQueryString) > 0 && urlWithQueryString[0] {
			if isMakeAll {
				query := ctx.Request().URL().Query()
				query.Set(`nocache`, `2`)
				reqURL += `?` + query.Encode()
			} else {
				if query := ctx.Request().URL().RawQuery(); len(query) > 0 {
					reqURL += `?` + query
				}
			}
		} else if isMakeAll {
			reqURL += `?nocache=2`
		}
		err := Make(http.MethodGet, reqURL, cacheKey)
		return true, err
	}
	return true, nil
}

func IsCached(ctx echo.Context, cacheKey string, urlWithQueryString ...bool) (bool, error) {
	if defaults.IsMockContext(ctx) {
		return false, nil
	}

	if customer := sessdata.Customer(ctx); customer != nil && customer.Uid > 0 {
		cached, err := controlCache(ctx, cacheKey, urlWithQueryString...)
		if err != nil || !cached {
			return cached, err
		}
	}
	var cachedETag sql.NullString
	getHash := func() string {
		if !cachedETag.Valid {
			cache.Get(context.Background(), cacheKey+`.hash`, &cachedETag.String)
			cachedETag.Valid = true
		}
		return cachedETag.String
	}
	getContent := func() (string, error) {
		var cachedHTML string
		err := cache.Get(context.Background(), cacheKey, &cachedHTML)
		return cachedHTML, err
	}
	err := ETagCallback(ctx, getHash, getContent)
	return err == nil, err
}

func ETagCallback(ctx echo.Context, contentEtag func() string, callback func() (string, error), weak ...bool) error {
	var _weak bool
	if len(weak) > 0 {
		_weak = weak[0]
	}
	if reqETag := ctx.Header(`If-None-Match`); len(reqETag) > 0 {
		if _weak {
			reqETag = strings.TrimPrefix(reqETag, `W/`)
		}
		if reqETag == `"`+contentEtag()+`"` {
			return ctx.NotModified()
		}
	}
	cachedHTML, err := callback()
	if err != nil {
		return err
	}
	if len(contentEtag()) > 0 {
		eTag := `"` + contentEtag() + `"`
		if _weak { // 弱ETag是指在资源内容发生变化时，ETag值不一定会随之改变的标识符
			eTag = `W/` + eTag
		}
		ctx.Response().Header().Set(`ETag`, eTag)
	}
	return ctx.HTML(cachedHTML)
}
