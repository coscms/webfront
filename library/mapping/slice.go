package mapping

import (
	"regexp"

	"github.com/webx-top/com"
)

type MappingFrom interface {
	GetField(string) interface{}
}

type MappingTo interface {
	MappingFrom
	Set(key interface{}, value ...interface{})
}

type MappingCfgGetter interface {
	GetFrom() interface{}
	GetTo() string
}

func From(from interface{}, to ...string) M {
	if len(to) > 0 {
		return M{KFrom: from, KTo: to[0]}
	}
	return M{KFrom: from}
}

type M struct {
	KFrom interface{}
	KTo   string
}

func (a M) To(to string) M {
	a.KTo = to
	return a
}

func (a M) GetFrom() interface{} {
	return a.KFrom
}

func (a M) GetTo() string {
	return a.KTo
}

type Layout string

var placeholder = regexp.MustCompile(`(?:%7B|\{)([^}%]+)(?:%7D|\})`)

func Slice[V MappingFrom, T MappingTo](queried []V, rows []T, linkKeyFrom string, linkKeyTo string, mapping ...MappingCfgGetter) []T {
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
		for _, mp := range mapping {
			switch v := mp.GetFrom().(type) {
			case string:
				newVal := srcRow.GetField(v)
				if newVal == nil {
					continue
				}
				rows[index].Set(mp.GetTo(), newVal)
			case Layout: // https://aaa/{Id} or https://aaa/?id=%7BId%7D
				newVal := placeholder.ReplaceAllStringFunc(string(v), func(s string) string {
					var k string
					if s[0] == '{' {
						k = s[1 : len(s)-1]
					} else {
						k = s[3 : len(s)-3]
					}
					v := srcRow.GetField(k)
					if v == nil {
						return ``
					}
					return com.String(v)
				})
				rows[index].Set(mp.GetTo(), newVal)
			case func(V) interface{}:
				newVal := v(srcRow)
				if newVal != nil {
					rows[index].Set(mp.GetTo(), newVal)
				}
			}
		}
	}
	return rows
}
