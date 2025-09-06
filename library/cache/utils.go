package cache

import (
	"github.com/admpub/cache"
	"github.com/admpub/cache/x"
	"github.com/admpub/color"
	"github.com/admpub/log"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
	"github.com/webx-top/pagination"

	"github.com/coscms/webcore/library/config/extend"
)

func IsDbAccount(v string) bool {
	return com.StrIsNumeric(v)
}

// ===============================
// reloader
// ===============================

var _ extend.Reloader = (*ReloadableOptions)(nil)

func NewReloadableOptions() *ReloadableOptions {
	return &ReloadableOptions{
		Options: &cache.Options{},
	}
}

type ReloadableOptions struct {
	*cache.Options
}

func (o *ReloadableOptions) Reload() error {
	err := CacheNew(cacheRootContext, *o.Options, `locker`)
	if err != nil {
		logPrefix := color.GreenString(`[cache]`) + `[locker][` + o.Adapter + `]`
		log.Error(logPrefix, err)
	} else {
		if o.Adapter == `redis` {
			resetRedsync()
			SetDefaultLockType(LockTypeRedis)
		} else {
			SetDefaultLockType(LockTypeMemory)
		}
	}
	return err
}

func (o *ReloadableOptions) IsValid() bool {
	return o != nil && o.Options != nil && len(o.Options.Adapter) > 0
}

// ===============================
// list
// ===============================

type List[T any] struct {
	List          []T
	PagingOptions echo.H
}

func (c *List[T]) SetList(list []T) *List[T] {
	c.List = list
	return c
}

func (c *List[T]) SetOptions(options echo.H) *List[T] {
	c.PagingOptions = options
	return c
}

func (c *List[T]) Do(ctx echo.Context, cacheKey string, fn func() ([]T, error), opts ...x.GetOption) error {
	//Delete(ctx, cacheKey)
	err := XFunc(ctx, cacheKey, &c, func() error {
		var err error
		c.List, err = fn()
		if err != nil {
			return err
		}
		c.PagingOptions = ctx.Get(`pagination`).(*pagination.Pagination).Options()
		return err
	}, opts...)
	return err
}

func (c *List[T]) GetPaginator(ctx echo.Context) *pagination.Pagination {
	return pagination.New(ctx).SetOptions(c.PagingOptions)
}

func (c *List[T]) GetList() []T {
	return c.List
}

func (c *List[T]) GetOptions() echo.H {
	return c.PagingOptions
}
