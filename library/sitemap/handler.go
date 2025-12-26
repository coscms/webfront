package sitemap

import (
	"path/filepath"
	"regexp"
	"strings"

	"github.com/webx-top/echo"
)

var hostVerifierR = regexp.MustCompile(`^[a-zA-Z0-9-]+(\.[a-zA-Z0-9-]+)*$`)

func VerifyHost(host string) bool {
	return hostVerifierR.MatchString(host)
}

func handleFile(c echo.Context, sitemapDir string, fileName string, getSubDirName func(echo.Context) string) error {
	root := sitemapDir + echo.FilePathSeparator
	if getSubDirName != nil {
		subDir := getSubDirName(c)
		if len(subDir) > 0 {
			if !VerifyHost(subDir) {
				return echo.ErrBadRequest
			}
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
			if !VerifyHost(subDir) {
				return echo.ErrBadRequest
			}
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
