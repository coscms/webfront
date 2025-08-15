package sitemap

import (
	"path/filepath"
	"strings"

	"github.com/admpub/sitemap-generator/smg"
	"github.com/webx-top/echo"
)

var Registry = echo.NewKVxData[Sitemap, any]()

func Register(k, v string, x Sitemap) {
	Registry.Add(k, v, echo.KVxOptX[Sitemap, any](x))
}

type Sitemap struct {
	Do func(echo.Context, *smg.Sitemap) error
}

func RegisterRoute(r echo.RouteRegister) {
	sitemapDir := filepath.Join(echo.Wd(), `public`, `sitemap`)
	r.Get(`/sitemap.xml`, func(c echo.Context) error {
		lang := c.Lang().Normalize()
		return c.File(sitemapDir + echo.FilePathSeparator + lang + echo.FilePathSeparator + `sitemap.xml`)
	})
	r.Get(`/sitemap_index.xml`, func(c echo.Context) error {
		lang := c.Lang().Normalize()
		return c.File(sitemapDir + echo.FilePathSeparator + lang + echo.FilePathSeparator + `sitemap_index.xml`)
	})
	r.Get("/sitemaps/*", func(c echo.Context) error {
		return static(c, sitemapDir)
	})
}

func static(c echo.Context, sitemapDir string) error {
	lang := c.Lang().Normalize()
	root := sitemapDir + echo.FilePathSeparator + lang + echo.FilePathSeparator + `sitemaps`
	var err error
	root, err = filepath.Abs(root)
	if err != nil {
		return err
	}
	name := filepath.Join(root, c.Param("*"))
	if !strings.HasPrefix(name, root) {
		return echo.ErrNotFound
	}
	return c.File(name)
}
