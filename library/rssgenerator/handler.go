package rssgenerator

import (
	"github.com/admpub/feeds"
	"github.com/webx-top/echo"

	"github.com/coscms/webcore/library/config"
)

// Handle RSS handler
//
// This function is used to generate RSS feeds. You can use it in your route by
// calling it with the context and the group name of the RSS feed.
//
// For example:
//
// echo.Get("/rss/:group", rssgenerator.Handle)
//
// This will generate RSS feeds for the given group name.
//
// If no group name is provided, it will use the default group name "article".
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

// HandleWith is a helper function that takes a context and a function to generate an RSS feed.
// It creates a new RSS feed with the site name, slogan, and link, and then calls the given function with the context and the feed.
// If the given function returns an error, it will be returned.
// Otherwise, it will write the RSS feed to the response and set the content type to "application/rss+xml; charset=utf-8".
// It is useful for generating RSS feeds for different groups of items.
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
