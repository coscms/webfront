package xhtml

import (
	"net/http"
	"net/url"

	"github.com/coscms/webfront/middleware/sessdata"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/defaults"
)

func MakeWithSiteURL(siteURL, method string, path string, saveAs string, reqRewrite ...func(*http.Request)) error {
	u, err := url.Parse(siteURL)
	if err != nil {
		return err
	}
	langCode := GetLangCodeByPath(path)
	saveAs = BuildCacheKey(u.Hostname(), langCode, saveAs)
	_, err, _ = makerSingleflight.Do(u.Hostname()+`_`+langCode+`_`+method+`_`+path, func() (interface{}, error) {
		return nil, makeDo(siteURL, method, path, saveAs, reqRewrite...)
	})
	return err
}

func IsCachedDomain(ctx echo.Context, cacheKey string, urlWithQueryString ...bool) (bool, error) {
	if defaults.IsMockContext(ctx) {
		return false, nil
	}
	if ctx.Echo().Multilingual() {
		cacheKey = ctx.Lang().Normalize() + `/` + cacheKey
	}
	cacheKey = ctx.Domain() + `/` + cacheKey
	if customer := sessdata.Customer(ctx); customer != nil && customer.Uid > 0 {
		cached, err := controlCache(ctx, cacheKey, urlWithQueryString...)
		if err != nil || !cached {
			return cached, err
		}
	}
	err := ETagCallback(ctx, cacheKey)
	return err == nil, err
}
