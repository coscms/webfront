package mapping

import (
	"fmt"
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
	case `Id`:
		f.Id = value[0].(int)
	}
}

func TestSlice(t *testing.T) {
	s := []*from{
		{Id: 100, Title: `title`},
	}
	d := []*to{
		{Id: 100, Title: `0`},
		{Id: 200, Title: ``},
	}
	r := Slice(s, d, `Id`, `Id`, M{`Title`, `Title`}, M{Layout(`https://abs/{Id}/{Title}`), `URL`})
	assert.Equal(t, []*to{
		{Id: 100, Title: `title`, URL: `https://abs/100/title`},
		{Id: 200, Title: ``},
	}, r)
}

func TestSlice2(t *testing.T) {
	s := []*from{
		{Id: 100, Title: `title`},
	}
	d := []*to{
		{Id: 100, Title: `0`},
		{Id: 200, Title: ``},
	}
	c := func(v *from) interface{} {
		return fmt.Sprintf(`https://abs/?id=%d&title=%s`, v.Id, v.Title)
	}
	r := Slice(s, d, `Id`, `Id`, From(`Title`).To(`Title`),
		M{c, `URL`},
	)
	assert.Equal(t, []*to{
		{Id: 100, Title: `title`, URL: `https://abs/?id=100&title=title`},
		{Id: 200, Title: ``},
	}, r)
}

func TestSlice3(t *testing.T) {
	s := []*from{
		{Id: 100, Title: `title`},
	}
	d := []*to{
		{Id: 100, Title: `0`},
		{Id: 200, Title: ``},
	}
	r := Slice(s, d, `Id`, `Id`, M{`Title`, `Title`},
		M{Layout(`https://abs/?id=` + url.QueryEscape(`{Id}`) + `&title=` + url.QueryEscape(`{Title}`)), `URL`},
	)
	assert.Equal(t, []*to{
		{Id: 100, Title: `title`, URL: `https://abs/?id=100&title=title`},
		{Id: 200, Title: ``},
	}, r)
}
