package sitemap

import (
	"path/filepath"

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
	r.File(`/sitemap.xml`, filepath.Join(echo.Wd(), `public`, `sitemap.xml`))
	r.File(`/sitemap_index.xml`, filepath.Join(echo.Wd(), `public`, `sitemap_index.xml`))
	r.Static(`/sitemaps`, filepath.Join(echo.Wd(), `public`, `sitemaps`))
}
