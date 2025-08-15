package sitemap

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/admpub/log"
	"github.com/admpub/sitemap-generator/smg"
	"github.com/coscms/webcore/library/config"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/defaults"
)

func GenerateIndex(ctx echo.Context, rootURL string, lang string, generateChildPageItems bool) error {
	if ctx == nil {
		ctx = defaults.NewMockContext()
	}
	now := time.Now().UTC()
	outputPath := filepath.Join(echo.Wd(), `public`, `sitemap`, lang)
	err := os.MkdirAll(outputPath, os.ModePerm)
	if err != nil {
		return err
	}
	smi := smg.NewSitemapIndex(true)
	smi.SetCompress(false)
	smi.SetSitemapIndexName(`sitemap_index`)
	smi.SetHostname(rootURL)
	smi.SetOutputPath(outputPath)
	smi.SetServerURI(`/sitemaps/`)

	for _, item := range Registry.Slice() {
		sm := smi.NewSitemap()
		sm.SetName(item.K)
		sm.SetLastMod(&now)
		sm.SetOutputPath(outputPath + echo.FilePathSeparator + `sitemaps`)
		if !generateChildPageItems {
			continue
		}
		err = item.X.Do(ctx, sm)
		if err != nil {
			return fmt.Errorf("unable to add sitemapLoc#%v: %v", item.K, err)
		}
	}

	var filename string
	filename, err = smi.Save()
	if err != nil {
		return fmt.Errorf("unable to save sitemap: %v", err)
	}
	log.Okayf(`sitemap generated successfully: %s`, filename)
	return err
}

func GenerateSingle(ctx echo.Context, rootURL string, lang string, f func(echo.Context, *smg.Sitemap) error) error {
	if ctx == nil {
		ctx = defaults.NewMockContext()
	}
	now := time.Now().UTC()
	outputPath := filepath.Join(echo.Wd(), `public`, `sitemap`, lang)
	err := os.MkdirAll(outputPath, os.ModePerm)
	if err != nil {
		return err
	}
	sm := smg.NewSitemap(true)
	sm.SetCompress(false)
	sm.SetName(`sitemap`)
	sm.SetHostname(rootURL)
	sm.SetOutputPath(outputPath)
	sm.SetLastMod(&now)

	err = f(ctx, sm)
	if err != nil {
		return fmt.Errorf("unable to add sitemapLoc: %v", err)
	}

	var filename []string
	filename, err = sm.Save()
	if err != nil {
		return fmt.Errorf("unable to save sitemap: %v", err)
	}
	log.Okayf(`sitemap generated successfully: %+v`, filename)
	return err
}

func RemoveAll() {
	os.RemoveAll(filepath.Join(echo.Wd(), `public`, `sitemap`))
}

func RemoveLanguage(lang string) {
	os.RemoveAll(filepath.Join(echo.Wd(), `public`, `sitemap`, lang))
}

func GenerateIndexAllLanguage(ctx echo.Context, rootURL string, generateChildPageItems bool) (err error) {
	for _, lang := range config.FromFile().Language.AllList {
		lang = echo.NewLangCode(lang).Normalize()
		err = GenerateIndex(ctx, rootURL, lang, generateChildPageItems)
		if err != nil {
			return
		}
	}
	return
}

func GenerateSingleAllLanguage(ctx echo.Context, rootURL string, f func(echo.Context, *smg.Sitemap) error) (err error) {
	for _, lang := range config.FromFile().Language.AllList {
		lang = echo.NewLangCode(lang).Normalize()
		err = GenerateSingle(ctx, rootURL, lang, f)
		if err != nil {
			return
		}
	}
	return
}
