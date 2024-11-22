package cache

import (
	"context"
	"path/filepath"
	"strings"
	"sync/atomic"

	"github.com/admpub/cache"
	"github.com/admpub/color"
	"github.com/admpub/log"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
	"golang.org/x/sync/singleflight"

	_ "github.com/admpub/cache/redis5" // redis
	_ "github.com/admpub/cache/sqlite" // sqlite

	"github.com/coscms/webcore/library/common"

	dbschemaDBMgr "github.com/nging-plugins/dbmanager/application/dbschema"
)

var (
	defaultCacheOptions = &cache.Options{
		Adapter:  `memory`,
		Interval: 300,
	}
	sg                     singleflight.Group
	defaultCacheInstance   atomic.Value // cache.Cache
	cacheConfigParsers     = map[string]func(cache.Options) (cache.Options, error){}
	instanceCachePrefix    = `Cache:`
	defaultConnectionName  = `default`
	fallbackConnectionName = `fallback`
)

func AddOptionParser(adapter string, parser func(cache.Options) (cache.Options, error)) {
	cacheConfigParsers[adapter] = parser
}

func Cache(ctx context.Context, args ...string) cache.Cache {
	if len(args) > 0 {
		defaultConnectionName = args[0]
		if len(args) > 1 {
			fallbackConnectionName = args[2]
		}
	}
	defaultKey := instanceCachePrefix + defaultConnectionName
	c, ok := echo.Get(defaultKey).(cache.Cache)
	if ok {
		return c
	}
	val, err, _ := sg.Do(`getCache`, func() (interface{}, error) {
		logPrefix := color.GreenString(`[cache]`)
		log.Warn(logPrefix, `[`+defaultConnectionName+`] 未找到已连接的实例`)
		var c cache.Cache
		var ok bool
		if defaultConnectionName != fallbackConnectionName {
			fallbackKey := instanceCachePrefix + fallbackConnectionName
			c, ok = echo.Get(fallbackKey).(cache.Cache)
			if ok {
				log.Warn(logPrefix, `[`+c.Name()+`] 使用备用实例`)
				echo.Set(defaultKey, c)
				return c, nil
			}
			log.Warn(logPrefix, `[`+fallbackConnectionName+`] 未找到已连接的实例`)
		}
		if c == nil {
			log.Warn(logPrefix, `[`+defaultCacheOptions.Adapter+`] 使用默认实例`)
			c, ok = defaultCacheInstance.Load().(cache.Cache)
			if !ok {
				if ctx == nil {
					ctx = context.Background()
				}
				var err error
				c, err = cache.Cacher(ctx, *defaultCacheOptions)
				if err != nil {
					log.Errorf(logPrefix, `[`+defaultCacheOptions.Adapter+`] 使用默认实例错误: %v`, err)
				} else {
					defaultCacheInstance.Store(c)
					echo.Set(defaultKey, c)
				}
			}
		}
		return c, nil
	})
	if err != nil {
		log.Error(err)
		return c
	}
	c = val.(cache.Cache)
	return c
}

func CacheNew(ctx context.Context, opts cache.Options, keys ...string) error {
	var connectionName string
	if len(keys) > 0 {
		connectionName = keys[0]
	}
	if len(connectionName) == 0 {
		connectionName = opts.Adapter
	}
	key := instanceCachePrefix + connectionName
	logPrefix := color.GreenString(`[cache]`) + `[` + connectionName + `][` + opts.Adapter + `]`
	if len(opts.AdapterConfig) == 0 {
		log.Okay(logPrefix, `缺少必要的配置参数，跳过`)
		return nil
	}
	c, ok := echo.Get(key).(cache.Cache)
	if ok {
		log.Okay(logPrefix, `断开连接`)
		echo.Delete(key)
		c.Close()
	}
	echo.FireByNameWithMap(`webx.cache.connected.`+opts.Adapter+`.before`, echo.H{`cache`: nil, `options`: opts})
	log.Info(logPrefix, `开始连接`)
	switch opts.Adapter {
	case `file`:
		if !filepath.IsAbs(opts.AdapterConfig) {
			opts.AdapterConfig = filepath.Join(echo.Wd(), opts.AdapterConfig)
		}
	case `sqlite`:
		if strings.HasSuffix(opts.AdapterConfig, echo.FilePathSeparator) {
			opts.AdapterConfig += `cache.sqlite`
		}
		if !filepath.IsAbs(opts.AdapterConfig) {
			opts.AdapterConfig = filepath.Join(echo.Wd(), opts.AdapterConfig)
		}
		if com.IsDir(opts.AdapterConfig) {
			opts.AdapterConfig = filepath.Join(opts.AdapterConfig, `cache.sqlite`)
		}
	case `redis`:
		if IsDbAccount(opts.AdapterConfig) {
			m := dbschemaDBMgr.NewNgingDbAccount(common.NewMockContext())
			err := m.Get(nil, `id`, opts.AdapterConfig)
			if err == nil {
				if len(m.Name) == 0 {
					m.Name = `0`
				}
				opts.AdapterConfig = `network=tcp,addr=` + m.Host + `,password=` + m.Password + `,db=` + m.Name + `,pool_size=100,idle_timeout=180,hset_name=Cache,prefix=cache:`
			}
		}
	}
	if cfgParser, ok := cacheConfigParsers[opts.Adapter]; ok {
		var err error
		opts, err = cfgParser(opts)
		if err != nil {
			log.Error(logPrefix, color.RedString(`配置解析失败:`+err.Error()))
			return err
		}
	}
	if ctx == nil {
		ctx = context.Background()
	}
	val, err := cache.Cacher(ctx, opts)
	if err != nil {
		log.Error(logPrefix, color.RedString(`连接失败:`+err.Error()))
		return err
	}
	echo.Set(key, val)
	return echo.FireByNameWithMap(`webx.cache.connected.`+opts.Adapter+`.after`, echo.H{`cache`: val, `options`: opts})
}
