package mapping

import (
	"regexp"
)

type MappingFrom interface {
	GetField(string) interface{}
}

type MappingTo interface {
	MappingFrom
	Set(key interface{}, value ...interface{})
}

type Layout string

var placeholder = regexp.MustCompile(`(?:%7B|\{)(.+)(?:%7D|\})`)

func Slice[V MappingFrom, T MappingTo](queried []V, rows []T, linkKeyFrom string, linkKeyTo string, mapping map[interface{}]string) []T {
	kk := map[interface{}]int{}
	for index, row := range rows {
		k := row.GetField(linkKeyTo)
		if k == nil {
			continue
		}
		kk[k] = index
	}
	for _, srcRow := range queried {
		k := srcRow.GetField(linkKeyFrom)
		if k == nil {
			continue
		}
		index, ok := kk[k]
		if !ok {
			continue
		}
		for from, to := range mapping {
			switch v := from.(type) {
			case string:
				newVal := srcRow.GetField(v)
				if newVal == nil {
					continue
				}
				rows[index].Set(to, newVal)
			case Layout:
				newVal := placeholder.ReplaceAllStringFunc(string(v), func(s string) string {
					var k string
					if s[0] == '{' {
						k = s[1 : len(s)-1]
					} else {
						k = s[3 : len(s)-3]
					}
					v, _ := srcRow.GetField(k).(string)
					return v
				})
				rows[index].Set(to, newVal)
			}
		}
	}
	return rows
}
