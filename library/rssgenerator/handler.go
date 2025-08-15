package rssgenerator

import (
	"github.com/admpub/feeds"
	"github.com/webx-top/echo"

	"github.com/coscms/webcore/library/config"
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
	if err := f(feed); err != nil {
		return err
	}

	ctx.Response().Header().Set(echo.HeaderContentType, "application/rss+xml; charset=utf-8")

	return feed.WriteRss(ctx.Response())
}
