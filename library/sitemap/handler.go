package sitemap

import (
	"path/filepath"
	"regexp"
	"strings"

	"github.com/webx-top/echo"
)

var hostVerifierR = regexp.MustCompile(`^[a-zA-Z0-9-]+(\.[a-zA-Z0-9-]+)*$`)

// VerifyHost checks if the given host is valid.
// A valid host is one that only contains characters of [a-zA-Z0-9-] and
// at least one period (.). The host must also end with a period.
func VerifyHost(host string) bool {
	return hostVerifierR.MatchString(host)
}

// handleFile writes the file at the given path to the HTTP response.
// The sitemapDir is the root directory of the sitemap files.
// The fileName is the name of the file to write.
// The getSubDirName is a function that returns the subdirectory name of the file.
// If getSubDirName is not nil, the returned subdirectory name will be appended to the root directory.
// The subdirectory name must be a valid hostname.
// If the subdirectory name is invalid, echo.ErrBadRequest will be returned.
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

// handleStatic writes the static file at the given path to the HTTP response.
// The sitemapDir is the root directory of the sitemap files.
// The getSubDirName is a function that returns the subdirectory name of the file.
// If getSubDirName is not nil, the returned subdirectory name will be appended to the root directory.
// The subdirectory name must be a valid hostname.
// If the subdirectory name is invalid, echo.ErrBadRequest will be returned.
// If the file is not found, echo.ErrNotFound will be returned.
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
