package xhtml

import (
	"errors"
	"net/http"
	"net/url"

	"github.com/coscms/webcore/library/httpserver"
	"github.com/coscms/webfront/middleware/sessdata"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/defaults"
)

var ErrNoSetValidateDomain = errors.New(`the DomainValidator function has not been set`)

func validateDomain(domain string) error {
	detected, err := httpserver.Frontend.ValidateDomain(domain)
	if err != nil || detected {
		return err
	}
	return ErrNoSetValidateDomain
}

func MakeWithSiteURL(siteURL, method string, path string, saveAs string, reqRewrite ...func(*http.Request)) error {
	u, err := url.Parse(siteURL)
	if err != nil {
		return err
	}
	if err = validateDomain(u.Hostname()); err != nil {
		return err
	}
	langCode := GetLangCodeByPath(path)
	saveAs = BuildCacheKey(u.Hostname(), langCode, saveAs)
	_, err, _ = makerSingleflight.Do(u.Hostname()+`_`+method+`_`+path, func() (interface{}, error) {
		return nil, makeDo(siteURL, method, path, saveAs, reqRewrite...)
	})
	return err
}

func IsCachedDomain(ctx echo.Context, cacheKey string, urlWithQueryString ...bool) (bool, error) {
	if defaults.IsMockContext(ctx) {
		return false, nil
	}
	err := validateDomain(ctx.Domain())
	if err != nil {
		return false, err
	}
	langCode := ctx.Lang().Normalize()
	cacheKey = langCode + `/` + cacheKey
	cacheKey = ctx.Domain() + `/` + cacheKey
	if customer := sessdata.Customer(ctx); customer != nil && customer.Uid > 0 {
		cached, err := controlCache(ctx, langCode, cacheKey, urlWithQueryString...)
		if err != nil || !cached {
			return cached, err
		}
	}
	err = ETagCallback(ctx, cacheKey)
	return err == nil, err
}
