package official

import "github.com/webx-top/echo"

var (
	// TagGroups 标签组
	TagGroups = echo.NewKVData()
)

func AddTagGroup(group string, title string) {
	TagGroups.Add(group, title)
}
