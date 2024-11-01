package top

import (
	"testing"
	"time"

	"github.com/webx-top/echo/testing/test"
)

func TestCN2Number(t *testing.T) {
	r, e := CN2Number(`一千三百二十八`)
	if e != nil {
		panic(e)
	}
	test.Eq(t, uint64(1328), r)

	r, e = CN2Number(`一千三百`)
	if e != nil {
		panic(e)
	}
	test.Eq(t, uint64(1300), r)

	r, e = CN2Number(`一`)
	if e != nil {
		panic(e)
	}
	test.Eq(t, uint64(1), r)

	r, e = CN2Number(`十`)
	if e != nil {
		panic(e)
	}
	test.Eq(t, uint64(10), r)

	r, e = CN2Number(`十五`)
	if e != nil {
		panic(e)
	}
	test.Eq(t, uint64(15), r)

	r, e = CN2Number(`十万零三`)
	if e != nil {
		panic(e)
	}
	test.Eq(t, uint64(100003), r)
	ti := time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC)
	t.Logf(`%s: %d`, ti.String(), ti.Unix())
	test.True(t, ti.IsZero())
}
