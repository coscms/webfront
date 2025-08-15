package sitemap

import (
	"path/filepath"
	"time"

	"github.com/admpub/sitemap-generator/smg"
	"github.com/webx-top/db"
	"github.com/webx-top/db/lib/factory/pagination"
	"github.com/webx-top/echo"

	"github.com/coscms/webcore/library/common"
	"github.com/coscms/webcore/library/httpserver"
	"github.com/coscms/webfront/middleware/sessdata"
	modelArticle "github.com/coscms/webfront/model/official/article"
)

var Registry = echo.NewKVxData[Sitemap, any]().
	Add(`article`, echo.T(`文章`), echo.KVxOptX[Sitemap, any](Sitemap{Do: articleSitemap}))

func Register(k, v string, x Sitemap) {
	Registry.Add(k, v, echo.KVxOptX[Sitemap, any](x))
}

type Sitemap struct {
	Do func(echo.Context, *smg.Sitemap) error
}

func articleSitemap(ctx echo.Context, sm *smg.Sitemap) error {
	source := ctx.Form(`source`)
	articleM := modelArticle.NewArticle(ctx)
	cond := db.NewCompounds()
	cond.AddKV(`display`, common.BoolY)
	if len(source) > 0 {
		cond.AddKV(`source_table`, source)
	}
	mw := func(r db.Result) db.Result {
		return r.Select(`id`, `created`, `updated`, `image`).OrderBy(`id`)
	}
	articleLastID := ctx.Internal().Uint64(`articleLastID`)
	if articleLastID > 0 {
		cond.AddKV(`id`, db.Gt(articleLastID))
	}
	ls := pagination.NewOffsetLister(articleM, nil, mw, cond.And())
	err := ls.ChunkList(func() (err error) {
		list := articleM.Objects()
		for _, row := range list {
			link := sessdata.URLByName(`article.detail`, row.Id)
			var lastMod time.Time
			if row.Updated > 0 {
				lastMod = time.Unix(int64(row.Updated), 0)
			} else {
				lastMod = time.Unix(int64(row.Created), 0)
			}
			item := &smg.SitemapLoc{
				Loc:        link,
				LastMod:    &lastMod,
				ChangeFreq: smg.Weekly,
				Priority:   PriorityDetail,
			}
			if len(row.Image) > 0 {
				item.Images = append(item.Images, &smg.SitemapImage{
					ImageLoc: row.Image,
				})
			}
			err = sm.Add(item)
			if err != nil {
				return
			}
			articleLastID = row.Id
		}
		return
	}, 100, 0)

	ctx.Internal().Set(`articleLastID`, articleLastID)
	return err
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
