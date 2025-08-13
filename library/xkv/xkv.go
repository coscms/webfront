package xkv

import (
	"context"
	"errors"

	"github.com/admpub/events"
	"github.com/coscms/webcore/dbschema"
	"github.com/coscms/webcore/model"
	"github.com/coscms/webfront/library/cache"
	"github.com/coscms/webfront/middleware/sessdata"
	"github.com/webx-top/echo"
)

func init() {
	echo.OnCallback(`nging.kv.delete`, func(e events.Event) error {
		kv := e.Context.Get(`kv`).(*dbschema.NgingKv)
		return RemoveCache(kv)
	})
	echo.OnCallback(`nging.kv.edit`, func(e events.Event) error {
		kv := e.Context.Get(`kv`).(*dbschema.NgingKv)
		return RemoveCache(kv)
	})
	echo.OnCallback(`nging.kv.add`, func(e events.Event) error {
		kv := e.Context.Get(`kv`).(*dbschema.NgingKv)
		return RemoveCache(kv)
	})
}

func RemoveCache(kv *dbschema.NgingKv) error {
	err1 := cache.Delete(context.Background(), `nging.kv.key.`+kv.Key)
	err2 := cache.Delete(context.Background(), `nging.kv.type.`+kv.Type)
	return errors.Join(err1, err2)
}

var DefaultTTL int64 = 86400 * 7

// GetValue 获取 key 的值
// defaultValue: 0. 默认值; 1. 说明; 2. 帮助说明 (1 和 2 仅在自动创建时有用)
func GetValue(ctx echo.Context, key string, defaultValue ...string) (string, error) {
	var value string
	err := cache.XFunc(ctx, `nging.kv.key.`+key, &value, func() (err error) {
		kvM := model.NewKv(ctx)
		value, err = kvM.GetValue(key, defaultValue...)
		return
	}, cache.AdminRefreshable(ctx, sessdata.Customer(ctx), cache.TTL(DefaultTTL)))
	return value, err
}

func GetTypeValues(ctx echo.Context, typ string, defaultValue ...string) ([]string, error) {
	var value []string
	err := cache.XFunc(ctx, `nging.kv.type.`+typ, &value, func() (err error) {
		kvM := model.NewKv(ctx)
		value, err = kvM.GetTypeValues(typ, defaultValue...)
		return
	}, cache.AdminRefreshable(ctx, sessdata.Customer(ctx), cache.TTL(DefaultTTL)))
	return value, err
}
