package cache

import (
	"github.com/admpub/once"
	goredislib "github.com/redis/go-redis/v9"
	"github.com/webx-top/echo"
)

var (
	redisClient *goredislib.Client
	redisOnce   once.Once
)

func init() {
	echo.OnCallback(`webx.cache.connected.redis.after`, func(data echo.Event) error {
		resetRedisClient()
		return nil
	})
}

func resetRedisClient() {
	redisOnce.Reset()
}

func initRedisClient() {
	defer resetRedsync()
	rc, ok := Cache(cacheRootContext, `default`).Client().(*goredislib.Client)
	if ok {
		redisClient = rc
		return
	}
	rc, _ = Cache(cacheRootContext, `fallback`).Client().(*goredislib.Client)
	redisClient = rc
}

func onceInitRedisClient() {
	initRedisClient()
}

func RedisClient() *goredislib.Client {
	redisOnce.Do(onceInitRedisClient)
	return redisClient
}

func RedisOptions() *goredislib.Options {
	opt, ok := Cache(cacheRootContext, `default`).(redisOptions)
	if ok {
		return opt.Options()
	}
	opt, ok = Cache(cacheRootContext, `fallback`).(redisOptions)
	if ok {
		return opt.Options()
	}
	return nil
}

type redisOptions interface {
	Options() *goredislib.Options
}
