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

func GenerateIndex(ctx echo.Context, rootURL string, langCode string, generateChildPageItems bool, subDir ...string) error {
	if ctx == nil {
		ctx = defaults.NewMockContext()
	}
	now := time.Now().UTC()
	var outputPath string
	if len(subDir) > 0 {
		outputPath = filepath.Join(echo.Wd(), `public`, `sitemap`, subDir[0], langCode)
	} else {
		outputPath = filepath.Join(echo.Wd(), `public`, `sitemap`, langCode)
	}
	err := os.MkdirAll(outputPath, os.ModePerm)
	if err != nil {
		return err
	}
	smi := smg.NewSitemapIndex(true)
	smi.SetCompress(false)
	smi.SetSitemapIndexName(`sitemap_index`)
	smi.SetHostname(rootURL)
	smi.SetOutputPath(outputPath)
	serverURI := `/sitemaps/`
	if len(langCode) > 0 && langCode != config.FromFile().Language.Default {
		serverURI = `/` + langCode + serverURI
	}
	smi.SetServerURI(serverURI)

	var subDirName string
	if len(subDir) > 0 {
		subDirName = subDir[0]
	}
	for _, item := range Registry.Slice() {
		sm := smi.NewSitemap()
		sm.SetName(item.K)
		sm.SetLastMod(&now)
		sm.SetOutputPath(outputPath + echo.FilePathSeparator + `sitemaps`)
		if !generateChildPageItems {
			continue
		}
		err = item.X.Do(ctx, sm, langCode, subDirName)
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

func GenerateSingle(ctx echo.Context, rootURL string, langCode string, f func(echo.Context, *smg.Sitemap, string, string) error, subDir ...string) error {
	if ctx == nil {
		ctx = defaults.NewMockContext()
	}
	now := time.Now().UTC()
	var outputPath string
	if len(subDir) > 0 {
		outputPath = filepath.Join(echo.Wd(), `public`, `sitemap`, subDir[0], langCode)
	} else {
		outputPath = filepath.Join(echo.Wd(), `public`, `sitemap`, langCode)
	}
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

	var subDirName string
	if len(subDir) > 0 {
		subDirName = subDir[0]
	}
	err = f(ctx, sm, langCode, subDirName)
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

func RemoveAll(subDirs ...string) {
	if len(subDirs) == 0 {
		os.RemoveAll(filepath.Join(echo.Wd(), `public`, `sitemap`))
		return
	}
	for _, subDir := range subDirs {
		os.RemoveAll(filepath.Join(echo.Wd(), `public`, `sitemap`, subDir))
	}
}

func RemoveLanguage(langCode string, subDirs ...string) {
	if len(subDirs) == 0 {
		os.RemoveAll(filepath.Join(echo.Wd(), `public`, `sitemap`, langCode))
		return
	}
	for _, subDir := range subDirs {
		os.RemoveAll(filepath.Join(echo.Wd(), `public`, `sitemap`, subDir, langCode))
	}
}

func GenerateIndexAllLanguage(ctx echo.Context, rootURL string, generateChildPageItems bool, subDir ...string) (err error) {
	for _, lang := range config.FromFile().Language.AllList {
		lang = echo.NewLangCode(lang).Normalize()
		err = GenerateIndex(ctx, rootURL, lang, generateChildPageItems, subDir...)
		if err != nil {
			return
		}
	}
	return
}

func GenerateSingleAllLanguage(ctx echo.Context, rootURL string, f func(echo.Context, *smg.Sitemap, string, string) error, subDir ...string) (err error) {
	for _, lang := range config.FromFile().Language.AllList {
		lang = echo.NewLangCode(lang).Normalize()
		err = GenerateSingle(ctx, rootURL, lang, f, subDir...)
		if err != nil {
			return
		}
	}
	return
}
