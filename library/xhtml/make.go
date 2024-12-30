package xhtml

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/admpub/log"
	"github.com/webx-top/com"

	"github.com/coscms/webfront/library/cache"
	"github.com/coscms/webfront/library/xcommon"
	"github.com/coscms/webfront/registry/route"
	test "github.com/webx-top/echo/testing"
	"golang.org/x/sync/singleflight"
)

var ErrGenerateHTML = errors.New(`failed to generate html`)
var makeSinglefight = singleflight.Group{}

func Make(method string, path string, saveAs string, reqRewrite ...func(*http.Request)) error {
	_, err, _ := makeSinglefight.Do(method+`_`+path, func() (interface{}, error) {
		return nil, makeDo(method, path, saveAs, reqRewrite...)
	})
	return err
}

func makeDo(method string, path string, saveAs string, reqRewrite ...func(*http.Request)) error {
	siteURL := xcommon.SiteURL(nil)
	if len(siteURL) == 0 {
		return fmt.Errorf(`%w: frontend URL cannot be empty`, ErrGenerateHTML)
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
