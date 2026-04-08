package cache

import (
	"context"

	"github.com/admpub/once"
)

var defaultSG, cancelDefaultSG = once.OnceValue(func() Singleflighter {
	return Singleflight()
})

func SingleflightDo(ctx context.Context, key string, fn func() (interface{}, error)) (v interface{}, err error, shared bool) {
	return defaultSG().Do(ctx, key, fn)
}

func SingleflightDoChan(ctx context.Context, key string, fn func() (interface{}, error)) <-chan SinglefightResult {
	return defaultSG().DoChan(ctx, key, fn)
}

func SingleflightForget(ctx context.Context, key string) {
	defaultSG().Forget(ctx, key)
}

func ResetSingleflight() {
	cancelDefaultSG()
}
