package cache

import (
	"context"

	"github.com/admpub/cache"
	"github.com/admpub/color"
	"github.com/admpub/log"
	"github.com/coscms/webcore/cmd/bootconfig"
	"github.com/coscms/webcore/library/captcha/driver/captcha_go"
	"github.com/coscms/webcore/library/common"
	"github.com/coscms/webcore/library/config"
	"github.com/coscms/webcore/library/config/extend"
	"github.com/webx-top/com"
	"github.com/webx-top/db/lib/factory"
)

var cacheRootContext context.Context
var cacheRootContextCancel context.CancelFunc

func CloseCache() {
	if cacheRootContextCancel != nil {
		cacheRootContextCancel()
	}
}

func connection(diffs config.Diffs) error {
	cfg := common.Setting(`cache`)
	adapters := make([]string, 0, 2)
	defaultCfg, ok := cfg.Get(`default`).(*cache.Options)
	if ok {
		err := CacheNew(cacheRootContext, *defaultCfg, `default`)
		if err != nil {
			logPrefix := color.GreenString(`[cache]`) + `[default][` + defaultCfg.Adapter + `]`
			log.Error(logPrefix, err)
		} else {
			adapters = append(adapters, defaultCfg.Adapter)
			log.Okayf(`cache system uses default storage engine: %s`, defaultCfg.Adapter)
		}
	}
	fallbackCfg, ok := cfg.Get(`fallback`).(*cache.Options)
	if ok {
		err := CacheNew(cacheRootContext, *fallbackCfg, `fallback`)
		if err != nil {
			logPrefix := color.GreenString(`[cache]`) + `[fallback][` + fallbackCfg.Adapter + `]`
			log.Error(logPrefix, err)
		} else {
			adapters = append(adapters, fallbackCfg.Adapter)
			log.Okayf(`cache system uses fallback storage engine: %s`, fallbackCfg.Adapter)
		}
	}
	lockerCfg, ok := config.FromFile().Extend.Get(`locker`).(*ReloadableOptions)
	if !ok || !lockerCfg.IsValid() {
		if com.InSlice(`redis`, adapters) {
			resetRedsync()
			SetDefaultLockType(LockTypeRedis)
			log.Okay(`locker system uses storage engine: redis`)
		} else {
			SetDefaultLockType(LockTypeMemory)
			log.Okay(`locker system uses storage engine: memory`)
		}
	}
	return nil
}

func init() {
	cacheRootContext, cacheRootContextCancel = context.WithCancel(context.Background())
	bootconfig.OnStart(-1, func() {
		config.OnGroupSetSettings(`cache`, connection)
		factory.SetCacher(DBCacher)
		extend.Register(`locker`, func() interface{} { return NewReloadableOptions() })
		captcha_go.DefaultStorer = captcha_go.NewStoreCacheWrapper(DBCacher)
	})
}
