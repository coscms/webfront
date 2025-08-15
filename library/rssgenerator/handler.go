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

func HandleWith(ctx echo.Context, f func(echo.Context, *feeds.Feed) error) error {
	cfg := config.Setting(`base`)
	contact := config.Setting(`contact`)
	feed := NewFeed(
		cfg.String(`siteName`),
		cfg.String(`siteSlogan`),
		ctx.Site(),
		cfg.String(`siteName`),
		contact.String(`customerServiceEmail`),
	)
	// 常见rel值

	// ‌alternate‌：表示文档存在不同格式的版本（如PDF、Word文档）
	// ‌enclosure‌：表示文档包含附件（如ZIP文件）
	// ‌stylesheet‌：指定样式表文件
	// ‌help‌：提供帮助文档
	// ‌search‌：提供搜索功能
	// self：当前文档网址

	feed.Links = []*feeds.Link{
		&feeds.Link{Rel: `self`, Href: ctx.FullRequestURI()},
	}
	feed.Copyright = `Copyright © ` + ctx.Domain()
	if err := f(ctx, feed); err != nil {
		return err
	}

	ctx.Response().Header().Set(echo.HeaderContentType, "application/rss+xml; charset=utf-8")

	return feed.WriteRss(ctx.Response())
}
