package top

import (
	"strings"

	"github.com/webx-top/com"
)

// ParseTags 将字符串解析为多个标签
// tags: string类型,可传入null、空字符串、以逗号分隔的多个标签字符串、或者以JSON格式的字符串
// 返回多个标签字符串的切片,或者错误
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
			if tags == `null` {
				return r, nil
			}
			r = strings.Split(tags, `,`)
		}
	}

	return r, err
}
