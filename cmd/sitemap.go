package cmd

import (
	"fmt"
	"net/url"
	"os"

	"github.com/spf13/cobra"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/defaults"
	"github.com/webx-top/echo/param"

	"github.com/coscms/webcore/cmd"
	"github.com/coscms/webcore/library/config"
	"github.com/coscms/webcore/library/filecache"
	"github.com/coscms/webfront/library/sitemap"
)

var sitemapCmd = &cobra.Command{
	Use:   "sitemap",
	Short: "sitemap generate",
	Long:  `Usage ./webx sitemap <rootURL> [languageCode]`,
	RunE:  sitemapRunE,
}

var sitemapCfg = struct {
	Mode     string // full(全量生成) / incr(增量生成)
	AllChild bool   // 是否同时生成所有子页面中的网址
}{
	Mode:     `incr`,
	AllChild: true,
}

func sitemapRunE(cmd *cobra.Command, args []string) error {
	var rootURL, lang string
	com.SliceExtract(args, &rootURL, &lang)
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
	eCtx := defaults.NewMockContext()
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
		err = sitemap.GenerateIndexAllLanguage(eCtx, rootURL, sitemapCfg.AllChild)
		return err
	}
	langs := param.Split(lang, `,`).Filter(func(s *string) bool {
		if s == nil {
			return false
		}
		if len(*s) == 0 {
			return false
		}
		*s = echo.NewLangCode(*s).Normalize()
		return true
	}).String()
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
	for _, _lang := range langs {
		err = sitemap.GenerateIndex(eCtx, rootURL, _lang, sitemapCfg.AllChild)
		if err != nil {
			fmt.Println(err.Error())
		}
	}
	return err
}

func init() {
	cmd.Add(sitemapCmd)
	minifyCmd.Flags().StringVar(&sitemapCfg.Mode, `incr`, sitemapCfg.Mode, `模式。支持的值有: full(全量生成) / incr(增量生成)`)
	minifyCmd.Flags().BoolVar(&sitemapCfg.AllChild, `allChild`, sitemapCfg.AllChild, `是否同时生成所有子页面中的网址`)
}
