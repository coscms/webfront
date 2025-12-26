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

func GenerateIndex(ctx echo.Context, rootURL string, langCodes []string, generateChildPageItems bool, subDir ...string) error {
	if ctx == nil {
		ctx = defaults.NewMockContext()
	}
	now := time.Now().UTC()
	var outputPath string
	if len(subDir) > 0 {
		outputPath = filepath.Join(echo.Wd(), `public`, `sitemap`, subDir[0])
	} else {
		outputPath = filepath.Join(echo.Wd(), `public`, `sitemap`)
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
		err = item.X.Run(ctx, sm, langCodes, item.K, subDirName)
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

func GenerateSingle(ctx echo.Context, rootURL string, langCodes []string, f *echo.KVx[Sitemap, any], subDir ...string) error {
	if ctx == nil {
		ctx = defaults.NewMockContext()
	}
	now := time.Now().UTC()
	var outputPath string
	if len(subDir) > 0 {
		outputPath = filepath.Join(echo.Wd(), `public`, `sitemap`, subDir[0])
	} else {
		outputPath = filepath.Join(echo.Wd(), `public`, `sitemap`)
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
	err = f.X.Run(ctx, sm, langCodes, f.K, subDirName)
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

func GenerateIndexAllLanguage(ctx echo.Context, rootURL string, generateChildPageItems bool, subDir ...string) (err error) {
	err = GenerateIndex(ctx, rootURL, config.FromFile().Language.AllList, generateChildPageItems, subDir...)
	return
}

func GenerateSingleAllLanguage(ctx echo.Context, rootURL string, f *echo.KVx[Sitemap, any], subDir ...string) (err error) {
	err = GenerateSingle(ctx, rootURL, config.FromFile().Language.AllList, f, subDir...)
	return
}
