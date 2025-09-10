package sitemap

import (
	"path/filepath"
	"strings"

	"github.com/coscms/webcore/library/config"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
)

func handleFile(c echo.Context, sitemapDir string, fileName string, getSubDirName func(echo.Context) string) error {
	lang := c.Lang().Normalize()
	root := sitemapDir + echo.FilePathSeparator
	if getSubDirName != nil {
		subDir := getSubDirName(c)
		if len(subDir) > 0 {
			root += subDir + echo.FilePathSeparator
		}
	}
	file := root + lang + echo.FilePathSeparator + fileName
	err := c.File(file)
	if err == nil || err != echo.ErrNotFound {
		return err
	}
	if lang != config.FromFile().Language.Default {
		lang = config.FromFile().Language.Default
		file = root + lang + echo.FilePathSeparator + fileName
		err = c.File(file)
	}
	return err
}

func handleStatic(c echo.Context, sitemapDir string, getSubDirName func(echo.Context) string) error {
	lang := c.Lang().Normalize()
	root := sitemapDir + echo.FilePathSeparator
	if getSubDirName != nil {
		subDir := getSubDirName(c)
		if len(subDir) > 0 {
			root += subDir + echo.FilePathSeparator
		}
	}
	root += lang + echo.FilePathSeparator + `sitemaps`
	var err error
	root, err = filepath.Abs(root)
	if err != nil {
		return err
	}
	reqFile := c.Param("*")
	reqFile = echo.CleanPath(reqFile)
	name := filepath.Join(root, reqFile)
	if !strings.HasPrefix(name, root) {
		return echo.ErrNotFound
	}
	err = c.File(name)
	if err == nil || err != echo.ErrNotFound || com.IsDir(root) {
		return err
	}
	langDefault := config.FromFile().Language.Default
	if lang != langDefault {
		root = strings.TrimSuffix(root, lang+echo.FilePathSeparator+`sitemaps`) + langDefault + echo.FilePathSeparator + `sitemaps`
		name = filepath.Join(root, reqFile)
		err = c.File(name)
	}
	return err
}
