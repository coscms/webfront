package rssgenerator

import (
	"github.com/admpub/feeds"
	"github.com/webx-top/echo"

	"github.com/coscms/webcore/library/config"
)

func Handle(ctx echo.Context) error {
	group := ctx.Param(`group`)
	if len(group) == 0 {
		group = ctx.Form(`group`)
	}
	if len(group) == 0 {
		group = `article`
	}
	item := Registry.GetItem(group)
	if item == nil {
		return echo.ErrNotFound
	}

	return HandleWith(ctx, item.X.Do)
}

func HandleWith(ctx echo.Context, f func(echo.Context, *feeds.RssFeed) error) error {
	cfg := config.Setting(`base`)
	feed := NewRssFeed(
		cfg.String(`siteName`),
		cfg.String(`siteSlogan`),
		ctx.Site(),
		cfg.String(`siteName`),
	)
	feed.Link = ctx.FullRequestURI()
	feed.Copyright = ctx.Domain()
	if err := f(ctx, feed); err != nil {
		return err
	}

	ctx.Response().Header().Set(echo.HeaderContentType, "application/rss+xml; charset=utf-8")

	return feeds.WriteXML(feed, ctx.Response())
}
