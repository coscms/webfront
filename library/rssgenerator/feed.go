package rssgenerator

import (
	"time"

	"github.com/admpub/feeds"
	"github.com/coscms/webcore/library/common"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/parser"
)

func NewFeed(title, slogan, host, authorName, authorEmail string) *feeds.Feed {
	now := time.Now()
	feed := &feeds.Feed{
		Title:       title,
		Description: slogan,
		Created:     now,
		Items:       []*feeds.Item{},
	}
	if len(authorName) > 0 || len(authorEmail) > 0 {
		feed.Author = &feeds.Author{Name: authorName, Email: authorEmail}
	}
	/*
		feed.Items = append(feed.Items, &feeds.Item{
			Id:          link,
			Title:       com.Substr(title,`...`,TitleMaxLength),
			Link:        &feeds.Link{Href: link},
			Description: MarkdownToHTML(``),
			Author:      &feeds.Author{Name: authorName, Email: authorEmail},
			Created:     *news.CreatedAt,
			Updated:     *news.UpdatedAt,
		})
	//*/
	return feed
}

func NewRssFeed(title, slogan, host, authorName string) *feeds.RssFeed {
	now := time.Now()
	feed := &feeds.RssFeed{
		Link:          host,
		Title:         title,
		Description:   slogan,
		LastBuildDate: now.Format(time.RFC1123Z),
		Items:         []*feeds.RssItem{},
	}
	if len(authorName) > 0 {
		feed.ManagingEditor = authorName
	}
	/*
		feed.Items = append(feed.Items, &feeds.RssItem{
			Guid:        &feeds.RssGuid{Id: link},
			Title:       title,
			Link:        link,
			Description: MarkdownToHTML(``),
			Author:      authorName,
			Category:    ``,
			PubDate:     *news.CreatedAt.Format(time.RFC1123Z),
		})
		//*/
	return feed
}

func MarkdownToHTML(md string) string {
	extensions := parser.NoIntraEmphasis | // 忽略单词内部的强调标记
		parser.Tables | // 解析表格语法
		parser.FencedCode | // 解析围栏代码块
		parser.Strikethrough | // 支持删除线语法
		parser.HardLineBreak | // 将换行符（\n）转换为 <br> 标签
		parser.Footnotes | // 支持脚注语法
		parser.MathJax | // 支持 MathJax 数学公式语法
		parser.SuperSubscript | // 支持上标和下标语法
		parser.EmptyLinesBreakList // 允许两个空行中断列表
	p := parser.NewWithExtensions(extensions)
	html := markdown.ToHTML([]byte(md), p, nil)
	cleanHTML := common.RemoveBytesXSS(html)
	return string(cleanHTML)
}
