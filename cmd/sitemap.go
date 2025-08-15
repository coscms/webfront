package cmd

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/defaults"
	"github.com/webx-top/echo/param"
	"github.com/webx-top/echo/subdomains"

	"github.com/coscms/webcore/cmd"
	"github.com/coscms/webcore/library/config"
	"github.com/coscms/webcore/library/filecache"
	"github.com/coscms/webcore/library/httpserver"
	"github.com/coscms/webfront/initialize/frontend"
	"github.com/coscms/webfront/library/sitemap"
	"github.com/coscms/webfront/registry/route"
)

var sitemapCmd = &cobra.Command{
	Use:   "sitemap",
	Short: "sitemap generate",
	Long: `Usage ./webx sitemap <rootURL> [languageCode]
删除所有: ./webx sitemap --mode=clear
删除指定域名或者指定语言: ./webx sitemap --mode=clear <rootURL> [languageCode]
删除所有域名下指定语言: ./webx sitemap --mode=clear all [languageCode]`,
	RunE: sitemapRunE,
}

var sitemapCfg = struct {
	Mode     string // full(全量生成) / incr(增量生成)
	AllChild bool   // 是否同时生成所有子页面中的网址
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

func validateLangCode(langs []string) error {
	validLangs := make([]string, len(config.FromFile().Language.AllList))
	for idx, lng := range config.FromFile().Language.AllList {
		validLangs[idx] = echo.NewLangCode(lng).Normalize()
	}
	var incorrectLang []string
	for _, _lang := range langs {
		if !com.InSlice(_lang, validLangs) {
			incorrectLang = append(incorrectLang, _lang)
		}
	}
	if len(incorrectLang) > 0 {
		return fmt.Errorf(`unsupported language: %+v`, incorrectLang)
	}
	return nil
}

func clearAllSubdirLanguages(langs []string) error {
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
	for _, _lang := range langs {
		sitemap.RemoveLanguage(_lang, subDirs...)
	}
	return err
}

func sitemapRunE(cmd *cobra.Command, args []string) error {
	var rootURL, lang string
	com.SliceExtract(args, &rootURL, &lang)
	if sitemapCfg.Mode == `clear` {
		switch rootURL {
		case ``, `all`:
			if len(lang) == 0 {
				sitemap.RemoveAll()
				return nil
			}
			langs := param.Split(lang, `,`).Filter(normalizeLangCode).String()
			err := validateLangCode(langs)
			if err != nil {
				return err
			}
			err = clearAllSubdirLanguages(langs)
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
	eCtx := defaults.NewMockContext()

	e := route.IRegister().Echo()
	subdomains.Default.Default = httpserver.KindFrontend
	subdomains.Default.Add(httpserver.KindFrontend+`@`+u.Host, e)
	route.Apply()
	frontend.SetRewriter(e)
	e.Commit()

	switch sitemapCfg.Mode {
	case `full`:

	case `incr`:
		for _, v := range sitemap.Registry.Slice() {
			b, err := filecache.ReadCache(`sitemap`, v.K)
			if err != nil && !os.IsNotExist(err) {
				return err
			}
			if len(b) > 0 {
				lastID := param.AsUint64(string(b))
				eCtx.Internal().Set(v.K+`LastID`, lastID)
			}
		}
	case `clear`:
		if len(lang) == 0 {
			sitemap.RemoveAll(subDir)
			fmt.Println(`removing sitemap is done`)
			return nil
		}
		langs := param.Split(lang, `,`).Filter(normalizeLangCode).String()
		err = validateLangCode(langs)
		if err != nil {
			return err
		}
		for _, _lang := range langs {
			sitemap.RemoveLanguage(_lang, subDir)
		}
		fmt.Println(`removing sitemap is done`)
		return nil
	}

	defer func() {
		for _, v := range sitemap.Registry.Slice() {
			lastID := eCtx.Internal().Uint64(v.K + `LastID`)
			if lastID > 0 {
				err := filecache.WriteCache(`sitemap`, v.K, []byte(param.AsString(lastID)))
				if err != nil {
					fmt.Println(err.Error())
				}
			}
		}
	}()
	if len(lang) == 0 {
		err = sitemap.GenerateIndexAllLanguage(eCtx, rootURL, sitemapCfg.AllChild, subDir)
		fmt.Println(`sitemap generation is complete`)
		return err
	}
	langs := param.Split(lang, `,`).Filter(normalizeLangCode).String()
	err = validateLangCode(langs)
	if err != nil {
		return err
	}
	for _, _lang := range langs {
		err = sitemap.GenerateIndex(eCtx, rootURL, _lang, sitemapCfg.AllChild, subDir)
		if err != nil {
			fmt.Println(err.Error())
		}
	}
	fmt.Println(`sitemap generation is complete`)
	return err
}

func init() {
	cmd.Add(sitemapCmd)
	sitemapCmd.Flags().StringVar(&sitemapCfg.Mode, `mode`, sitemapCfg.Mode, `模式。支持的值有: full(全量生成) / incr(增量生成) / clear(删除)`)
	sitemapCmd.Flags().BoolVar(&sitemapCfg.AllChild, `allChild`, sitemapCfg.AllChild, `是否同时生成所有子页面中的网址`)
}
