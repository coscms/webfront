package sitemap

import (
	"errors"
	"fmt"
	"net/url"
	"os"

	"github.com/admpub/log"
	"github.com/coscms/webcore/cmd/bootconfig"
	"github.com/coscms/webcore/library/config"
	"github.com/coscms/webcore/library/filecache"
	"github.com/coscms/webfront/initialize/frontend"
	"github.com/coscms/webfront/registry/route"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/defaults"
	"github.com/webx-top/echo/param"
)

func NewConfig() Config {
	return Config{
		Mode:     `incr`,
		AllChild: true,
	}
}

type Config struct {
	Mode     string // full(全量生成) / incr(增量生成) / clear(删除)
	Group    string
	AllChild bool // 是否同时生成所有子页面中的网址
}

func normalizeLangCode(s *string) bool {
	if s == nil {
		return false
	}
	if len(*s) == 0 {
		return false
	}
	*s = echo.NewLangCode(*s).Normalize()
	return true
}

func validateLangCode(langCodes []string) error {
	validLangs := make([]string, len(config.FromFile().Language.AllList))
	for idx, lng := range config.FromFile().Language.AllList {
		validLangs[idx] = echo.NewLangCode(lng).Normalize()
	}
	var incorrectLang []string
	for _, _lang := range langCodes {
		if !com.InSlice(_lang, validLangs) {
			incorrectLang = append(incorrectLang, _lang)
		}
	}
	if len(incorrectLang) > 0 {
		return fmt.Errorf(`unsupported language: %+v`, incorrectLang)
	}
	return nil
}

func CmdGenerate(rootURL, langCode string, sitemapCfg Config) error {
	if sitemapCfg.Mode == `clear` {
		switch rootURL {
		case ``, `all`:
			RemoveAll()
			log.Info(`removing sitemap is done`)
			return nil
		}
	}
	if len(rootURL) == 0 {
		return fmt.Errorf(`please specify the website root URL: %s sitemap <rootURL>`, os.Args[0])
	}
	u, err := url.Parse(rootURL)
	if err != nil {
		return err
	}
	if len(u.Host) == 0 || len(u.Scheme) == 0 {
		return fmt.Errorf(`invalid root URL: %s`, rootURL)
	}
	if !bootconfig.IsWeb() {
		err = config.ParseConfig()
		if err != nil {
			return err
		}
	}
	subDir := u.Hostname()
	if sitemapCfg.Mode == `clear` {
		RemoveAll(subDir)
		log.Info(`removing sitemap is done`)
		return nil
	}

	var groupItems []*echo.KVx[Sitemap, any]
	if len(sitemapCfg.Group) > 0 {
		groups := param.Split(sitemapCfg.Group, `,`).Filter().String()
		for _, group := range groups {
			item := Registry.GetItem(group)
			if item != nil {
				groupItems = append(groupItems, item)
			}
		}
	} else {
		groupItems = Registry.Slice()
	}

	if len(groupItems) == 0 {
		return errors.New(`No group found`)
	}

	eCtx := defaults.NewMockContext(route.IRegister().Echo())
	eCtx.Request().SetHost(u.Host)
	eCtx.Request().Header().Set(echo.HeaderXForwardedProto, u.Scheme)

	if !bootconfig.IsWeb() {
		frontend.TempInitRoute(u.Host)
	}

	var prepare func() error
	if sitemapCfg.Mode == `incr` {
		prepare = func() error {
			for _, v := range groupItems {
				b, err := filecache.ReadCache(`sitemap`, v.K+`_`+subDir+`.txt`)
				if err != nil && !os.IsNotExist(err) {
					return err
				}
				if len(b) > 0 {
					lastID := param.AsUint64(string(b))
					eCtx.Internal().Set(`sitemapGen.`+subDir+`.`+v.K+`LastID`, lastID) // sitemapGen.<Hostname>.articleLastID
				}
			}
			return nil
		}
	}
	if prepare == nil {
		prepare = func() error { return nil }
	}
	save := func() {
		var err error
		for _, v := range groupItems {
			lastID := eCtx.Internal().Uint64(`sitemapGen.` + subDir + `.` + v.K + `LastID`) // sitemapGen.<Hostname>.articleLastID
			if lastID <= 0 {
				continue
			}
			err = filecache.WriteCache(`sitemap`, v.K+`_`+subDir+`.txt`, []byte(param.AsString(lastID)))
			if err != nil {
				log.Error(err.Error())
			}
		}
	}
	cfg := config.FromFile()
	var langCodes []string
	if len(langCode) == 0 {
		langCodes = make([]string, len(cfg.Language.AllList))
		for index, lang := range cfg.Language.AllList {
			langCodes[index] = echo.NewLangCode(lang).Normalize()
		}
	} else {
		langCodes = param.Split(langCode, `,`).Filter(normalizeLangCode).String()
		err = validateLangCode(langCodes)
		if err != nil {
			return err
		}
	}
	if err = prepare(); err != nil {
		return err
	}
	err = GenerateIndex(eCtx, rootURL, langCodes, sitemapCfg.AllChild, subDir)
	log.Info(`sitemap generation is complete`)
	save()
	return err
}
