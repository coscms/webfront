package sitemap

import (
	"path/filepath"
	"time"

	"github.com/admpub/sitemap-generator/smg"
	"github.com/webx-top/db"
	"github.com/webx-top/db/lib/factory/pagination"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/param"

	"github.com/coscms/webcore/library/common"
	"github.com/coscms/webcore/library/httpserver"
	modelArticle "github.com/coscms/webfront/model/official/article"
)

var Registry = echo.NewKVxData[Sitemap, any]().
	Add(`article`, echo.T(`文章`), echo.KVxOptX[Sitemap, any](Sitemap{Do: articleSitemap}))

func Register(k, v string, x Sitemap) {
	Registry.Add(k, v, echo.KVxOptX[Sitemap, any](x))
}

type LocGenerator func(ctx echo.Context, sm *smg.Sitemap, langCodes []string, lastID string) (newLastID string, err error)

type Sitemap struct {
	Do LocGenerator
}

func (a Sitemap) Run(ctx echo.Context, sm *smg.Sitemap, langCodes []string, name string, subDirName string) error {
	inkey := `sitemapGen.` + subDirName + `.` + name + `LastID`
	lastID := ctx.Internal().String(inkey)
	var err error
	lastID, err = a.Do(ctx, sm, langCodes, lastID)
	if len(lastID) > 0 {
		ctx.Internal().Set(inkey, lastID)
	}
	return err
}

func articleSitemap(ctx echo.Context, sm *smg.Sitemap, langCodes []string, lastID string) (string, error) {
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
	articleLastID := param.AsUint64(lastID)
	if articleLastID > 0 {
		cond.AddKV(`id`, db.Gt(articleLastID))
	}
	ls := pagination.NewOffsetLister(articleM, nil, mw, cond.And())
	err := ls.ChunkList(func() (err error) {
		list := articleM.Objects()
		for _, row := range list {
			link := ctx.URLByName(`article.detail`, row.Id)
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
			if len(langCodes) > 1 {
				item.Alternate = make([]*smg.SitemapAlternateLoc, 0, len(langCodes)+1)
				item.Alternate = append(item.Alternate, &smg.SitemapAlternateLoc{
					Hreflang: `x-default`,
					Href:     link,
					Rel:      `alternate`,
				})
				relativeLink := ctx.RelativeURLByName(`article.detail`, row.Id)
				for _, langCode := range langCodes {
					item.Alternate = append(item.Alternate, &smg.SitemapAlternateLoc{
						Hreflang: langCode,
						Href:     ctx.SiteRoot() + `/` + langCode + relativeLink,
						Rel:      `alternate`,
					})
				}
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

	lastID = param.AsString(articleLastID)
	return lastID, err
}

func RegisterRoute(r echo.RouteRegister, getSubDirName func(echo.Context) string) {
	sitemapDir := filepath.Join(echo.Wd(), `public`, `sitemap`)
	r.Get(`/sitemap.xml`, func(c echo.Context) error {
		return handleFile(c, sitemapDir, `sitemap.xml`, getSubDirName)
	}).SetMetaKV(httpserver.PermGuestKV())
	r.Get(`/sitemap_index.xml`, func(c echo.Context) error {
		return handleFile(c, sitemapDir, `sitemap_index.xml`, getSubDirName)
	}).SetMetaKV(httpserver.PermGuestKV())
	r.Get("/sitemaps/*", func(c echo.Context) error {
		return handleStatic(c, sitemapDir, getSubDirName)
	}).SetMetaKV(httpserver.PermGuestKV())
}
