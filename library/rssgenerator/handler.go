package rssgenerator

import (
	"github.com/coscms/webcore/library/config"
	"github.com/gorilla/feeds"
	"github.com/webx-top/echo"
)

func RSS(ctx echo.Context, f func(*feeds.Feed) error) error {
	cfg := config.Setting(`base`)
	contact := config.Setting(`contact`)
	feed := NewFeed(
		cfg.String(`siteName`),
		cfg.String(`siteSlogan`),
		ctx.Site(),
		cfg.String(`siteName`),
		contact.String(`customerServiceEmail`),
	)
	feed.Link = &feeds.Link{Href: ctx.FullRequestURI()}
	if err := f(feed); err != nil {
		return err
	}

	ctx.Response().Header().Set(echo.HeaderContentType, "application/rss+xml; charset=utf-8")

	return feed.WriteRss(ctx.Response())
}
