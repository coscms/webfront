package mapping

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

type from struct {
	Id    int
	Title string
}

func (f from) GetField(key string) interface{} {
	switch key {
	case `Title`:
		return f.Title
	case `Id`:
		return f.Id
	default:
		return nil
	}
}

type to struct {
	Id    int
	Title string
	URL   string
}

func (f to) GetField(key string) interface{} {
	switch key {
	case `Title`:
		return f.Title
	case `URL`:
		return f.URL
	case `Id`:
		return f.Id
	default:
		return nil
	}
}

func (f *to) Set(key interface{}, value ...interface{}) {
	switch key.(string) {
	case `Title`:
		f.Title = value[0].(string)
	case `URL`:
		f.URL = value[0].(string)
	}
}

func TestSlice(t *testing.T) {
	s := []*from{
		{Id: 100, Title: `1`},
	}
	d := []*to{
		{Id: 100, Title: `0`},
		{Id: 200, Title: ``},
	}
	r := Slice(s, d, `Id`, `Id`, map[interface{}]string{
		`Title`:                       `Title`,
		Layout(`https://abs/{Title}`): `URL`,
	})
	assert.Equal(t, []*to{
		{Id: 100, Title: `1`, URL: `https://abs/1`},
		{Id: 200, Title: ``},
	}, r)
}

func TestSlice2(t *testing.T) {
	s := []*from{
		{Id: 100, Title: `1`},
	}
	d := []*to{
		{Id: 100, Title: `0`},
		{Id: 200, Title: ``},
	}
	r := Slice(s, d, `Id`, `Id`, map[interface{}]string{
		`Title`: `Title`,
		Layout(`https://abs/?id=` + url.QueryEscape(`{Title}`)): `URL`,
	})
	assert.Equal(t, []*to{
		{Id: 100, Title: `1`, URL: `https://abs/?id=1`},
		{Id: 200, Title: ``},
	}, r)
}
