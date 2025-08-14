package sitemap

import (
	"path/filepath"

	"github.com/admpub/sitemap-generator/smg"
	"github.com/webx-top/echo"
)

var Registry = echo.NewKVxData[Sitemap, any]()

type Adder func(*smg.SitemapLoc) error

type Sitemap struct {
	Do func(add Adder) error
}

func RegisterRoute(r echo.RouteRegister) {
	r.File(`/sitemap.xml`, filepath.Join(echo.Wd(), `public`, `sitemap.xml`))
	r.File(`/sitemap_index.xml`, filepath.Join(echo.Wd(), `public`, `sitemap_index.xml`))
	r.Static(`/sitemaps`, filepath.Join(echo.Wd(), `public`, `sitemaps`))
}
