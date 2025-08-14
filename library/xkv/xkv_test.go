package xkv

import (
	"context"
	"testing"

	"github.com/coscms/webfront/library/cache"
	"github.com/stretchr/testify/assert"
)

func TestCacheString(t *testing.T) {
	testFn := func() (string, error) {
		var value string
		err := cache.XFunc(context.Background(), `test`, &value, func() (err error) {
			value = `pong`
			println(`~~~~~~~~~~~~~~~~~~~~~~~~~~~~~query`)
			return
		}, cache.TTL(DefaultTTL))
		return value, err
	}
	value, err := testFn()
	assert.NoError(t, err)
	assert.Equal(t, `pong`, value)

	value, err = testFn()
	assert.NoError(t, err)
	assert.Equal(t, `pong`, value)
}
