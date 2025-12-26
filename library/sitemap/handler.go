package sitemap

import (
	"path/filepath"
	"strings"

	"github.com/webx-top/echo"
)

func handleFile(c echo.Context, sitemapDir string, fileName string, getSubDirName func(echo.Context) string) error {
	root := sitemapDir + echo.FilePathSeparator
	if getSubDirName != nil {
		subDir := getSubDirName(c)
		if len(subDir) > 0 {
			root += subDir + echo.FilePathSeparator
		}
	}
	file := root + echo.FilePathSeparator + fileName
	return c.File(file)
}

func handleStatic(c echo.Context, sitemapDir string, getSubDirName func(echo.Context) string) error {
	root := sitemapDir + echo.FilePathSeparator
	if getSubDirName != nil {
		subDir := getSubDirName(c)
		if len(subDir) > 0 {
			root += subDir + echo.FilePathSeparator
		}
	}
	root += `sitemaps`
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
	return err
}
