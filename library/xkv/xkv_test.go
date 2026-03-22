package xkv

import (
	"context"
	"sync/atomic"
	"testing"

	"github.com/coscms/webfront/library/cache"
	"github.com/stretchr/testify/assert"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/defaults"
)

func TestCacheString(t *testing.T) {
	testFn := func() (echo.KVList, error) {
		var value echo.KVList
		err := cache.XFunc(context.Background(), `test`, &value, func() (err error) {
			value = echo.KVList{
				echo.NewKV(`k`, `v`),
			}
			println(`~~~~~~~~~~~~~~~~~~~~~~~~~~~~~query`)
			return
		}, cache.TTL(DefaultTTL))
		return value, err
	}
	expected := echo.KVList{
		echo.NewKV(`k`, `v`),
	}
	value, err := testFn()
	assert.NoError(t, err)
	assert.Equal(t, expected, value)

	value, err = testFn()
	assert.NoError(t, err)
	assert.Equal(t, expected, value)

	ctx := defaults.NewMockContext()
	var i int32
	f := func() (int, error) {
		return int(atomic.AddInt32(&i, 1)), nil
	}
	resultN, err := GetOnce(ctx, ``, f)
	assert.NoError(t, err)
	assert.Equal(t, 1, resultN)

	resultN, err = GetOnce(ctx, ``, f)
	assert.NoError(t, err)
	assert.Equal(t, 1, resultN)
}
