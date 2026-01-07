// Package rssgenerator contains functions for generating RSS feeds.
package rssgenerator

import (
	"time"

	"github.com/admpub/feeds"
	"github.com/coscms/webcore/library/common"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/parser"
	"github.com/webx-top/com"
)

// NewFeed creates a new feeds.Feed instance.
//
// title is the title of the rss feed.
//
// slogan is the slogan of the rss feed.
//
// host is the hostname of the rss feed.
//
// authorName and authorEmail are the name and email of the author of the rss feed.
//
// The function returns a new feeds.Feed instance with the Title, Description, Created, Items, and Author fields set.
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

// NewRssFeed creates a new feeds.RssFeed instance.
//
// title is the title of the rss feed.
//
// slogan is the slogan of the rss feed.
//
// host is the hostname of the rss feed.
//
// authorName is the name of the author, it is optional.
//
// The function returns a pointer to the new feeds.RssFeed instance.
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

// CDATA wraps a string with CDATA tags, which are used to escape special characters in XML.
// See https://www.w3.org/TR/REC-xml/#sec-cdata for more information.
func CDATA(s string) string {
	return `<![CDATA[` + s + `]]>`
}

// MarkdownToHTML 将 Markdown 字符串转换为 HTML 字符串。
// 该函数将 Markdown 字符串解析成 HTML 字符串，并将结果返回。
// 该函数使用了以下 Markdown 语法扩展项：
//   - NoIntraEmphasis：忽略单词内部的强调标记
//   - Tables：解析表格语法
//   - FencedCode：解析围栏代码块
//   - Strikethrough：支持删除线语法
//   - HardLineBreak：将换行符（\n）转换为 <br> 标签
//   - Footnotes：支持脚注语法
//   - MathJax：支持 MathJax 数学公式语法
//   - SuperSubscript：支持上标和下标语法
//   - EmptyLinesBreakList：允许两个空行中断列表
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
	b := com.Str2bytes(md)
	b = markdown.ToHTML(b, p, nil)
	b = common.RemoveBytesXSS(b)
	return com.Bytes2str(b)
}
