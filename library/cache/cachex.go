package cache

import (
	"context"

	"github.com/admpub/cache/x"
	"github.com/coscms/webfront/dbschema"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/defaults"
)

type QueryFunc = x.QueryFunc

const (
	CacheDisabled int64 = -1 // 禁用缓存
	CacheFresh    int64 = -2 // 强制缓存新数据
)

var (
	TTL               = x.TTL
	Query             = x.Query
	DisableCacheUsage = x.DisableCacheUsage
	UseFreshData      = x.UseFreshData
	Disabled          = DisableCacheUsage(true) // 禁用缓存
	Fresh             = UseFreshData(true)      // 强制缓存新数据
	Noop              = func(o *x.Options) {}
)

func AdminRefreshable(ctx echo.Context, customer *dbschema.OfficialCustomer, ttl x.GetOption) x.GetOption {
	if !defaults.IsMockContext(ctx) {
		if customer == nil {
			return ttl
		}
		if customer.Uid <= 0 {
			return ttl
		}
	}

	nocache := ctx.Formx(`nocache`).Bool()
	return TTLIf(nocache, Fresh, ttl)
}

func GetTTLByNumber(ttl int64, b x.GetOption) x.GetOption {
	switch ttl {
	case CacheDisabled: // 禁用缓存
		return func(o *x.Options) {
			if b != nil {
				b(o)
			}
			Disabled(o)
		}
	case CacheFresh: // 强制缓存新数据
		return func(o *x.Options) {
			if b != nil {
				b(o)
			}
			Fresh(o)
		}
	default:
		if b != nil {
			return b
		}
		if ttl <= 0 {
			return Noop
		}
		return TTL(ttl)
	}
}

func TTLIf(condition bool, a x.GetOption, b x.GetOption) x.GetOption {
	if condition {
		return func(o *x.Options) {
			b(o)
			a(o)
		}
	}
	return b
}

func TTLIfCallback(condition func() bool, a x.GetOption, b x.GetOption) x.GetOption {
	if condition() {
		return func(o *x.Options) {
			b(o)
			a(o)
		}
	}
	return b
}

func AddTTL(a x.GetOption, ttl int64) x.GetOption {
	return func(o *x.Options) {
		a(o)
		o.AddTTL(ttl)
	}
}

func GenOptions(ctx echo.Context, cacheSeconds int64) []x.GetOption {
	nocache := ctx.Formx(`nocache`).Int()
	opts := []x.GetOption{TTL(cacheSeconds)}
	switch nocache {
	case 1:
		opts = append(opts, DisableCacheUsage(true)) // 禁用缓存
	case 2:
		opts = append(opts, UseFreshData(true)) // 强制缓存新数据
	case 3:
		opts = append(opts, UseFreshData(true))
	}
	return opts
}

func xNewFromPool(query x.Querier, ttlSeconds int64, args ...string) *x.Cachex {
	c := Cache(cacheRootContext, args...)
	return x.NewFromPool(c, query, ttlSeconds)
}

// XQuery 获取缓存，如果不存在则执行函数获取数据并缓存【自动避免缓存穿透】
func XQuery(ctx context.Context, key string, recv interface{}, query x.Querier, options ...x.GetOption) error {
	a := xNewFromPool(query, 0)
	defer a.Release()
	return a.Get(ctx, key, recv, options...)
}

// XFunc 获取缓存，如果不存在则执行函数获取数据并缓存【自动避免缓存穿透】
func XFunc(ctx context.Context, key string, recv interface{}, fn func() error, options ...x.GetOption) error {
	a := xNewFromPool(QueryFunc(fn), 0)
	defer a.Release()
	return a.Get(ctx, key, recv, options...)
}

// Delete 删除缓存
func Delete(ctx context.Context, key string, args ...string) error {
	return Cache(cacheRootContext, args...).Delete(ctx, key)
}

// Put 设置缓存
func Put(ctx context.Context, key string, val interface{}, timeout int64, args ...string) error {
	return Cache(cacheRootContext, args...).Put(ctx, key, val, timeout)
}

// Get 获取缓存
func Get(ctx context.Context, key string, recv interface{}, args ...string) error {
	return Cache(cacheRootContext, args...).Get(ctx, key, recv)
}

// Incr 自增
func Incr(ctx context.Context, key string, args ...string) error {
	return Cache(cacheRootContext, args...).Incr(ctx, key)
}

// Decr 自减
func Decr(ctx context.Context, key string, args ...string) error {
	return Cache(cacheRootContext, args...).Decr(ctx, key)
}

// IsExist 是否存在缓存
func IsExist(ctx context.Context, key string, args ...string) (bool, error) {
	return Cache(cacheRootContext, args...).IsExist(ctx, key)
}

// Flush 删除所有缓存数据
func Flush(ctx context.Context, args ...string) error {
	return Cache(cacheRootContext, args...).Flush(ctx)
}
