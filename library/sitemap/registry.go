package sitemap

import (
	"path/filepath"

	"github.com/admpub/sitemap-generator/smg"
	"github.com/coscms/webcore/library/httpserver"
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
		return handleFile(c, sitemapDir, `sitemap.xml`)
	}).SetMetaKV(httpserver.PermGuestKV())
	r.Get(`/sitemap_index.xml`, func(c echo.Context) error {
		return handleFile(c, sitemapDir, `sitemap_index.xml`)
	}).SetMetaKV(httpserver.PermGuestKV())
	r.Get("/sitemaps/*", func(c echo.Context) error {
		return handleStatic(c, sitemapDir)
	}).SetMetaKV(httpserver.PermGuestKV())
}
