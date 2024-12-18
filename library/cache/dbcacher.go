package cache

import (
	"context"

	"github.com/admpub/cache/x"
	"github.com/coscms/captcha"
	"github.com/webx-top/db/lib/factory"
)

var DBCacher = NewDBCacher()
var _ factory.Cacher = (*dbCacher)(nil)
var _ captcha.Storer = (*dbCacher)(nil)

func NewDBCacher() *dbCacher {
	return &dbCacher{}
}

type dbCacher struct{}

func (d *dbCacher) Put(ctx context.Context, key string, value interface{}, ttlSeconds int64) error {
	return Cache(cacheRootContext).Put(ctx, key, value, ttlSeconds)
}

func (d *dbCacher) Del(ctx context.Context, key string) error {
	return Cache(cacheRootContext).Delete(ctx, key)
}

func (d *dbCacher) Delete(ctx context.Context, key string) error {
	return d.Del(ctx, key)
}

func (d *dbCacher) Get(ctx context.Context, key string, value interface{}) error {
	return Cache(cacheRootContext).Get(ctx, key, value)
}

func (d *dbCacher) Do(ctx context.Context, key string, recv interface{}, fn func() error, ttlSeconds int64) error {
	return XQuery(ctx, key, recv, x.QueryFunc(fn), TTL(ttlSeconds))
}
