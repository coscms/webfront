package top

import (
	"strings"

	"github.com/webx-top/com"
)

// ParseTags 将字符串解析为字符串切片
// 如果字符串以 `[` 开头，则认为是 JSON 字符串，否则认为是逗号分隔的字符串
// 例如 `["a","b","c"]` 或 `a,b,c`
func ParseTags(tags string) ([]string, error) {
	var (
		r   []string
		err error
	)

	if len(tags) > 0 {
		if tags[0] == '[' && tags[len(tags)-1] == ']' {
			b := com.Str2bytes(tags)
			err = com.JSONDecode(b, &r)
			if err != nil {
				return r, err
			}
		} else {
			r = strings.Split(tags, `,`)
		}
	}

	return r, err
}
