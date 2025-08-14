package sitemap

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/admpub/log"
	"github.com/admpub/sitemap-generator/smg"
	"github.com/webx-top/echo"
)

func GenerateIndex(rootURL string) error {
	now := time.Now().UTC()
	outputPath := filepath.Join(echo.Wd(), `public`)
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
		err = item.X.Do(sm.Add)
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

func GenerateSingle(rootURL string, f func(Adder) error) error {
	now := time.Now().UTC()
	outputPath := filepath.Join(echo.Wd(), `public`)
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

	err = f(sm.Add)
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
	for _, xmlFilePath := range []string{
		filepath.Join(echo.Wd(), `public`, `sitemap.xml`),
		filepath.Join(echo.Wd(), `public`, `sitemap_index.xml`),
		filepath.Join(echo.Wd(), `public`, `sitemaps`),
	} {
		os.RemoveAll(xmlFilePath)
	}
}
