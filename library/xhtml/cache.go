package xhtml

import (
	"context"
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
	case `4`, `rm`:
		err := Remove(cacheKey)
		return false, err
	}
	return true, nil
}

func IsCached(ctx echo.Context, cacheKey string, urlWithQueryString ...bool) (bool, error) {
	if defaults.IsMockContext(ctx) {
		return false, nil
	}

	if ctx.Echo().Multilingual() {
		cacheKey = ctx.Lang().Normalize() + `/` + cacheKey
	}
	if customer := sessdata.Customer(ctx); customer != nil && customer.Uid > 0 {
		cached, err := controlCache(ctx, cacheKey, urlWithQueryString...)
		if err != nil || !cached {
			return cached, err
		}
	}
	err := ETagCallback(ctx, cacheKey)
	return err == nil, err
}

func getHash(c context.Context, cacheKey string) string {
	var etag string
	cache.Get(c, cacheKey+`.hash`, &etag)
	return etag
}

func getContent(c context.Context, cacheKey string) (string, error) {
	var cachedHTML string
	err := cache.Get(c, cacheKey, &cachedHTML)
	return cachedHTML, err
}

func ETagCallback(ctx echo.Context, cacheKey string, weak ...bool) error {
	var _weak bool
	var tagValue string
	var tagGetted bool
	if len(weak) > 0 {
		_weak = weak[0]
	}
	if reqETag := ctx.Header(`If-None-Match`); len(reqETag) > 0 {
		if _weak {
			reqETag = strings.TrimPrefix(reqETag, `W/`)
		}
		tagValue = getHash(ctx, cacheKey)
		tagGetted = true
		if reqETag == `"`+tagValue+`"` {
			return ctx.NotModified()
		}
	}
	cachedHTML, err := getContent(ctx, cacheKey)
	if err != nil {
		return err
	}
	if !tagGetted {
		tagValue = getHash(ctx, cacheKey)
	}
	if len(tagValue) > 0 {
		eTag := `"` + tagValue + `"`
		if _weak { // 弱ETag是指在资源内容发生变化时，ETag值不一定会随之改变的标识符
			eTag = `W/` + eTag
		}
		ctx.Response().Header().Set(`ETag`, eTag)
	}
	return ctx.HTML(cachedHTML)
}
