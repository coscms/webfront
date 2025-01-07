package official

import "github.com/webx-top/echo"

var (
	Contype = echo.NewKVData()
)

func init() {
	Contype.Add(`text`, echo.T(`纯文本`))
	Contype.Add(`html`, echo.T(`富文本`))
	Contype.Add(`markdown`, `Markdown`)
}
