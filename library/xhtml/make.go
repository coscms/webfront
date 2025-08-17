package xhtml

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/admpub/log"
	"github.com/webx-top/com"

	"github.com/coscms/webcore/library/config"
	"github.com/coscms/webfront/library/cache"
	"github.com/coscms/webfront/library/xcommon"
	"github.com/coscms/webfront/registry/route"
	test "github.com/webx-top/echo/testing"
	"golang.org/x/sync/singleflight"
)

var ErrGenerateHTML = errors.New(`failed to generate html`)
var makerSingleflight = singleflight.Group{}

func BuildCacheKey(domain string, langCode string, cacheKey string) string {
	if len(langCode) > 0 {
		cacheKey = langCode + `/` + cacheKey
	}
	if len(domain) > 0 {
		cacheKey = domain + `/` + cacheKey
	}
	return cacheKey
}

func Make(method string, path string, saveAs string, reqRewrite ...func(*http.Request)) error {
	_, err, _ := makerSingleflight.Do(method+`_`+path, func() (interface{}, error) {
		return nil, makeDo(``, method, path, saveAs, reqRewrite...)
	})
	return err
}

func MakeWithSiteURL(siteURL, method string, path string, saveAs string, reqRewrite ...func(*http.Request)) error {
	u, err := url.Parse(siteURL)
	if err != nil {
		return err
	}
	langCode := config.FromFile().Language.Default
	lcInPath := strings.SplitN(strings.TrimPrefix(path, `/`), `/`, 2)[0]
	if langCode != lcInPath &&
		com.InSlice(lcInPath, config.FromFile().Language.AllList) {
		langCode = lcInPath
	}
	saveAs = BuildCacheKey(u.Hostname(), langCode, saveAs)
	_, err, _ = makerSingleflight.Do(u.Hostname()+langCode+`_`+`_`+method+`_`+path, func() (interface{}, error) {
		return nil, makeDo(siteURL, method, path, saveAs, reqRewrite...)
	})
	return err
}

func makeDo(siteURL, method string, path string, saveAs string, reqRewrite ...func(*http.Request)) error {
	if len(siteURL) == 0 {
		siteURL = xcommon.SiteURL(nil)
		if len(siteURL) == 0 {
			return fmt.Errorf(`%w: frontend URL cannot be empty`, ErrGenerateHTML)
		}
	}
	if strings.HasPrefix(path, `/`) {
		siteURL = strings.TrimSuffix(siteURL, `/`)
	} else if !strings.HasSuffix(siteURL, `/`) {
		siteURL = siteURL + `/`
	}
	rec := test.Request(method, siteURL+path, route.IRegister().Echo(), reqRewrite...)
	if rec.Code != http.StatusOK {
		err := fmt.Errorf(`%w: [%d] %v`, ErrGenerateHTML, rec.Code, rec.Body.String())
		log.Error(err)
		return err
	}
	body := rec.Body.String()
	err := cache.Put(context.Background(), saveAs, body+`<!-- Generated at `+time.Now().Format(time.DateTime)+` -->`, 0)
	if err != nil {
		log.Error(err)
	} else {
		err = cache.Put(context.Background(), saveAs+`.hash`, com.Md5(body), 0)
		if err != nil {
			log.Error(err)
		}
	}
	return err
}

func Remove(cacheKey string) error {
	cache.Delete(context.Background(), cacheKey+`.hash`)
	return cache.Delete(context.Background(), cacheKey)
}
