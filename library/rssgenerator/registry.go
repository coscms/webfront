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
	Do func(ctx echo.Context, feed *feeds.RssFeed) error
}

func T2s(t time.Time) string {
	return t.Format(time.RFC1123Z)
}

func articleRSS(ctx echo.Context, feed *feeds.RssFeed) error {
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
		item := &feeds.RssItem{
			Guid:        &feeds.RssGuid{Id: link},
			Title:       row.Title,
			Link:        link,
			Description: row.Summary,
		}
		if row.Updated > 0 {
			item.PubDate = T2s(time.Unix(int64(row.Updated), 0))
		} else {
			item.PubDate = T2s(time.Unix(int64(row.Created), 0))
		}
		if idx == 0 {
			feed.PubDate = item.PubDate
			feed.LastBuildDate = item.PubDate
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
		if len(row.Image) > 0 {
			mtype := mime.TypeByExtension(path.Ext(row.Image))
			row.Image = com.AbsURL(ctx.Site()+`rss`, row.Image)
			item.Enclosure = &feeds.RssEnclosure{
				Url:  row.Image,
				Type: mtype,
			}
		}
		if row.Customer != nil {
			item.Author = row.Customer.Name
		} else if row.User != nil {
			item.Author = row.User.Username
		}
		if row.Category != nil {
			item.Category = row.Category.Name
		}
		feed.Items = append(feed.Items, item)
	}
	return err
}

func RegisterRoute(r echo.RouteRegister) {
	r.Get(`/rss`, Handle).SetMetaKV(httpserver.PermGuestKV())
	r.Get(`/rss/:group`, Handle).SetMetaKV(httpserver.PermGuestKV())
}
