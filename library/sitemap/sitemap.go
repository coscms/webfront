// Package sitemap contains functions for generating sitemap files.
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

// GenerateIndex 生成sitemap索引文件
//
//   - ctx: Context对象
//
//   - rootURL: 网站根URL
//
//   - langCodes: 语言代码,多个语言代码用逗号分隔
//
//   - generateChildPageItems: 是否生成所有子页面中的网址
//
//   - subDir: 生成sitemap文件的子目录,可选
//
//     生成的sitemap文件将被保存在public/sitemap/目录下
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

// GenerateSingle 生成sitemap单个文件
//
//   - ctx: Context对象
//
//   - rootURL: 网站根URL
//
//   - langCodes: 语言代码,多个语言代码用逗号分隔
//
//   - f: 需要生成sitemap的函数
//
//   - subDir: 生成sitemap文件的子目录,可选
//
//     生成的sitemap文件将被保存在public/sitemap/目录下
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

// RemoveAll removes all sitemap files from the public/sitemap/ directory.
//
// If subDirs is empty, it will remove the entire public/sitemap/ directory.
//
// Otherwise, it will remove the specified subdirectories from the public/sitemap/ directory.
//
// For example, if subDirs is ["zh-CN", "en-US"], it will remove the public/sitemap/zh-CN/ and public/sitemap/en-US/ directories.
//
// Note that this function will not return an error even if the directory does not exist.
func RemoveAll(subDirs ...string) {
	if len(subDirs) == 0 {
		os.RemoveAll(filepath.Join(echo.Wd(), `public`, `sitemap`))
		return
	}
	for _, subDir := range subDirs {
		os.RemoveAll(filepath.Join(echo.Wd(), `public`, `sitemap`, subDir))
	}
}

// GenerateIndexAllLanguage generates the sitemap index file for all languages.
//
// It is a shorthand for calling GenerateIndex with config.FromFile().Language.AllList.
//
// See GenerateIndex for more information.
func GenerateIndexAllLanguage(ctx echo.Context, rootURL string, generateChildPageItems bool, subDir ...string) (err error) {
	err = GenerateIndex(ctx, rootURL, config.FromFile().Language.AllList, generateChildPageItems, subDir...)
	return
}

// GenerateSingleAllLanguage generates the sitemap index file for a single sitemap generator.
//
// It is a shorthand for calling GenerateSingle with config.FromFile().Language.AllList.
//
// See GenerateSingle for more information.
func GenerateSingleAllLanguage(ctx echo.Context, rootURL string, f *echo.KVx[Sitemap, any], subDir ...string) (err error) {
	err = GenerateSingle(ctx, rootURL, config.FromFile().Language.AllList, f, subDir...)
	return
}
