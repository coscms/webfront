package rssgenerator

import (
	"mime"
	"path"
	"strings"
	"time"

	"github.com/admpub/feeds"
	"github.com/webx-top/com"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"

	"github.com/coscms/webcore/library/common"
	"github.com/coscms/webcore/library/httpserver"
	"github.com/coscms/webfront/middleware/sessdata"
	modelArticle "github.com/coscms/webfront/model/official/article"
)

var Registry = echo.NewKVxData[RSS, any]().
	Add(`article`, echo.T(`文章`), echo.KVxOptX[RSS, any](RSS{Do: articleRSS}))

func Register(k, v string, x RSS) {
	Registry.Add(k, v, echo.KVxOptX[RSS, any](x))
}

type RSS struct {
	Do func(ctx echo.Context, feed *feeds.Feed) error
}

func articleRSS(ctx echo.Context, feed *feeds.Feed) error {
	source := ctx.Form(`source`)
	articleM := modelArticle.NewArticle(ctx)
	list := []*modelArticle.ArticleWithOwner{}
	cond := db.NewCompounds()
	cond.AddKV(`display`, common.BoolY)
	if len(source) > 0 {
		cond.AddKV(`source_table`, source)
	}
	_, err := articleM.ListByOffset(&list, func(r db.Result) db.Result {
		return r.OrderBy(`-id`)
	}, 0, 20, cond.And())
	if err != nil {
		return err
	}
	for idx, row := range list {
		link := sessdata.URLByName(`article.detail`, row.Id)
		if strings.HasPrefix(link, `/`) {
			link = ctx.Site() + link[1:]
		}
		item := &feeds.Item{
			Id:          link,
			Title:       row.Title,
			Link:        &feeds.Link{Href: link},
			Description: row.Summary,
			Author:      nil,
			Created:     time.Unix(int64(row.Created), 0),
		}
		if row.Updated > 0 {
			item.Updated = time.Unix(int64(row.Updated), 0)
		}
		if idx == 0 {
			feed.Created = item.Created
			feed.Updated = item.Updated
		}
		if len(item.Description) == 0 {
			switch row.Contype {
			case common.ContentTypeHTML:
				item.Description = row.Content
			case common.ContentTypeMarkdown:
				item.Description = MarkdownToHTML(row.Content)
			case common.ContentTypeText:
				item.Description = row.Content
			}
		}
		if row.Customer != nil {
			item.Author = &feeds.Author{Name: row.Customer.Name, Email: ``}
		} else if row.User != nil {
			item.Author = &feeds.Author{Name: row.User.Username, Email: ``}
		}
		if len(row.Image) > 0 {
			mtype := mime.TypeByExtension(path.Ext(row.Image))
			row.Image = com.AbsURL(ctx.Site()+`rss`, row.Image)
			item.Enclosure = &feeds.Enclosure{
				Url:  row.Image,
				Type: mtype,
			}
		}
		feed.Add(item)
	}
	return err
}

func RegisterRoute(r echo.RouteRegister) {
	r.Get(`/rss`, Handle).SetMetaKV(httpserver.PermGuestKV())
	r.Get(`/rss/:group`, Handle).SetMetaKV(httpserver.PermGuestKV())
}
