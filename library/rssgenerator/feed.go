package rssgenerator

import (
	"time"

	"github.com/coscms/webcore/library/common"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/parser"
	"github.com/gorilla/feeds"
)

const TitleMaxLength = 20

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
			Author:      &feeds.Author{Name: authorName, Email: authorName},
			Created:     *news.CreatedAt,
			Updated:     *news.UpdatedAt,
		})
	*/
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
