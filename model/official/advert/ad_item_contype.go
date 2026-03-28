package advert

import (
	"context"
	"fmt"
	"strings"

	"github.com/webx-top/echo"
)

type (
	ContentRenderer func(Adverter) string
	Adverter        interface {
		GetWidth() uint
		GetHeight() uint
		GetURL() string
		GetContent() string
		GetContype() string
		GetTitle() string
		GetDescription() string
	}
)

var (
	Contype = echo.NewKVData()
	AdMode  = echo.NewKVData()
)

func Render(v Adverter) string {
	item := Contype.GetItem(v.GetContype())
	if item == nil {
		return ``
	}
	if fn, ok := item.Fn()(nil).(ContentRenderer); ok {
		return fn(v)
	}
	return ``
}

func GenStyle(p Adverter) string {
	if p == nil {
		return ``
	}
	styles := make([]string, 0, 3)
	if p.GetWidth() > 0 {
		styles = append(styles, fmt.Sprintf(`width:%dpx`, p.GetWidth()))
	} else {
		styles = append(styles, `width:100%`)
	}
	if p.GetHeight() > 0 {
		styles = append(styles, fmt.Sprintf(`height:%dpx`, p.GetHeight()))
	}
	styles = append(styles, `max-width:100%`)
	style := strings.Join(styles, `;`)
	if len(style) > 0 {
		style = ` style="` + style + `"`
	}
	return style
}

func GenLink(v Adverter, content string) string {
	url := v.GetURL()
	if len(url) == 0 {
		return content
	}
	return `<a href="` + url + `" target="_blank" title="` + v.GetTitle() + `">` + content + `</a>`
}

func init() {
	Contype.AddItem(echo.NewKV(`text`, echo.T(`文字广告`)).SetHKV(`description`, echo.T(`输入广告文字`)).SetFn(func(c context.Context) interface{} {
		return ContentRenderer(func(v Adverter) string {
			return GenLink(v, v.GetContent())
		})
	}))
	Contype.AddItem(echo.NewKV(`image`, echo.T(`图片广告`)).SetHKV(`description`, echo.T(`输入图片文件网址`)).SetFn(func(c context.Context) interface{} {
		return ContentRenderer(func(v Adverter) string {
			return GenLink(v, `<img rel="`+v.GetContent()+`" class="previewable" src="`+v.GetContent()+`"`+GenStyle(v)+` />`)
		})
	}))
	Contype.AddItem(echo.NewKV(`video`, echo.T(`视频广告`)).SetHKV(`description`, echo.T(`输入视频文件网址`)).SetFn(func(c context.Context) interface{} {
		return ContentRenderer(func(v Adverter) string {
			return GenLink(v, `<video src="`+v.GetContent()+`" controls="controls"`+GenStyle(v)+`></video>`)
		})
	}))
	Contype.AddItem(echo.NewKV(`audio`, echo.T(`音频广告`)).SetHKV(`description`, echo.T(`输入音频文件网址`)).SetFn(func(c context.Context) interface{} {
		return ContentRenderer(func(v Adverter) string {
			return GenLink(v, `<audio src="`+v.GetContent()+`" controls="controls"`+GenStyle(v)+`></audio>`)
		})
	}))

	AdMode.Add(`CPA`, `CPA`)
	AdMode.Add(`CPM`, `CPM`)
	AdMode.Add(`CPC`, `CPC`)
	AdMode.Add(`CPS`, `CPS`)
	AdMode.Add(`CPT`, `CPT`)
}
