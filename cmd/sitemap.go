package cmd

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/defaults"
	"github.com/webx-top/echo/param"

	"github.com/coscms/webcore/cmd"
	"github.com/coscms/webcore/library/config"
	"github.com/coscms/webcore/library/filecache"
	"github.com/coscms/webfront/initialize/frontend"
	"github.com/coscms/webfront/library/sitemap"
)

var sitemapCmd = &cobra.Command{
	Use:   "sitemap",
	Short: "sitemap generate",
	Long: `Usage ./webx sitemap <rootURL> [languageCode...]
删除所有: ./webx sitemap --mode=clear
删除指定域名或者指定语言: ./webx sitemap --mode=clear <rootURL> en,zh-CN
删除所有域名下指定语言: ./webx sitemap --mode=clear all en,zh-CN`,
	RunE: sitemapRunE,
}

var sitemapCfg = struct {
	Mode     string // full(全量生成) / incr(增量生成)
	Group    string
	AllChild bool // 是否同时生成所有子页面中的网址
}{
	Mode:     `incr`,
	AllChild: true,
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

func clearAllSubdirLanguages(langCodes []string) error {
	root := filepath.Join(echo.Wd(), `public`, `sitemap`)
	dirs, err := os.ReadDir(root)
	if err != nil {
		return err
	}
	var subDirs []string
	for _, info := range dirs {
		if info.IsDir() {
			subDirs = append(subDirs, info.Name())
		}
	}
	if len(subDirs) == 0 {
		return err
	}
	for _, _lang := range langCodes {
		sitemap.RemoveLanguage(_lang, subDirs...)
	}
	return err
}

func sitemapRunE(cmd *cobra.Command, args []string) error {
	var rootURL, langCode string
	com.SliceExtract(args, &rootURL, &langCode)
	if sitemapCfg.Mode == `clear` {
		switch rootURL {
		case ``, `all`:
			if len(langCode) == 0 {
				sitemap.RemoveAll()
				return nil
			}
			langCodes := param.Split(langCode, `,`).Filter(normalizeLangCode).String()
			err := validateLangCode(langCodes)
			if err != nil {
				return err
			}
			err = clearAllSubdirLanguages(langCodes)
			fmt.Println(`removing sitemap is done`)
			return err
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
	err = config.ParseConfig()
	if err != nil {
		return err
	}
	subDir := u.Hostname()
	if sitemapCfg.Mode == `clear` {
		if len(langCode) == 0 {
			sitemap.RemoveAll(subDir)
			fmt.Println(`removing sitemap is done`)
			return nil
		}
		langCodes := param.Split(langCode, `,`).Filter(normalizeLangCode).String()
		err = validateLangCode(langCodes)
		if err != nil {
			return err
		}
		for _, _lang := range langCodes {
			sitemap.RemoveLanguage(_lang, subDir)
		}
		fmt.Println(`removing sitemap is done`)
		return nil
	}

	var groupItems []*echo.KVx[sitemap.Sitemap, any]
	if len(sitemapCfg.Group) > 0 {
		groups := param.Split(sitemapCfg.Group, `,`).Filter().String()
		for _, group := range groups {
			item := sitemap.Registry.GetItem(group)
			if item != nil {
				groupItems = append(groupItems, item)
			}
		}
	} else {
		groupItems = sitemap.Registry.Slice()
	}

	if len(groupItems) == 0 {
		return errors.New(`No group found`)
	}

	eCtx := defaults.NewMockContext()

	frontend.TempInitRoute(u.Host)

	var prepare func(langCodes []string) error
	if sitemapCfg.Mode == `incr` {
		prepare = func(langCodes []string) error {
			for _, v := range groupItems {
				for _, _lang := range langCodes {
					b, err := filecache.ReadCache(`sitemap`, v.K+`_`+_lang+`_`+subDir+`.txt`)
					if err != nil && !os.IsNotExist(err) {
						return err
					}
					if len(b) > 0 {
						lastID := param.AsUint64(string(b))
						eCtx.Internal().Set(subDir+`.`+_lang+`.`+v.K+`LastID`, lastID)
					}
				}
			}
			return nil
		}
	}
	if prepare == nil {
		prepare = func(_ []string) error { return nil }
	}
	save := func(langCodes []string) {
		for _, v := range groupItems {
			for _, _lang := range langCodes {
				lastID := eCtx.Internal().Uint64(subDir + `.` + _lang + `.` + v.K + `LastID`)
				if lastID <= 0 {
					continue
				}
				err := filecache.WriteCache(`sitemap`, v.K+`_`+_lang+`_`+subDir+`.txt`, []byte(param.AsString(lastID)))
				if err != nil {
					fmt.Println(err.Error())
				}
			}
		}
	}
	var langCodes []string
	if len(langCode) == 0 {
		langCodes = make([]string, len(config.FromFile().Language.AllList))
		for index, lang := range config.FromFile().Language.AllList {
			langCodes[index] = echo.NewLangCode(lang).Normalize()
		}
	} else {
		langCodes = param.Split(langCode, `,`).Filter(normalizeLangCode).String()
		err = validateLangCode(langCodes)
		if err != nil {
			return err
		}
	}
	if err = prepare(langCodes); err != nil {
		return err
	}
	for _, _lang := range langCodes {
		err = sitemap.GenerateIndex(eCtx, rootURL, _lang, sitemapCfg.AllChild, subDir)
		if err != nil {
			fmt.Println(err.Error())
		}
	}
	fmt.Println(`sitemap generation is complete`)
	save(langCodes)
	return err
}

func init() {
	cmd.Add(sitemapCmd)
	sitemapCmd.Flags().StringVar(&sitemapCfg.Mode, `mode`, sitemapCfg.Mode, `模式。支持的值有: full(全量生成) / incr(增量生成) / clear(删除)`)

	slic := sitemap.Registry.Slice()
	groups := make([]string, len(slic))
	for i, v := range slic {
		groups[i] = v.K
	}
	sitemapCmd.Flags().StringVar(&sitemapCfg.Group, `group`, sitemapCfg.Group, `组。指定多个时用半角逗号“,”隔开, 支持的值有: `+strings.Join(groups, ` / `)+` 等。不指定时代表所有组`)
	sitemapCmd.Flags().BoolVar(&sitemapCfg.AllChild, `allChild`, sitemapCfg.AllChild, `是否同时生成所有子页面中的网址`)
}
