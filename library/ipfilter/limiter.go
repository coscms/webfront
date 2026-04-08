package ipfilter

import (
	"context"
	"crypto/tls"
	"time"

	"github.com/admpub/color"
	"github.com/admpub/log"
	goredislib "github.com/redis/go-redis/v9"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/defaults"
	"github.com/webx-top/echo/middleware/ratelimiter"
	"github.com/webx-top/echo/param"

	dbschemaDBMgr "github.com/nging-plugins/dbmanager/application/dbschema"
)

func NewRateLimiterConfig() *RateLimiterConfig {
	return &RateLimiterConfig{}
}

type RateLimiterConfig struct {
	On bool
	// The max count in duration for no policy, default is 100.
	Max int

	// Count duration for no policy, default is 1 Minute (60s).
	Duration int64 //seconds
	//key prefix, default is "LIMIT:".
	Prefix string

	//If request gets a  internal limiter error, just skip the limiter and let it go to next middleware
	SkipInternalError bool

	RedisAddr     string
	RedisPassword string
	RedisDB       int
	DBAccountID   uint
}

func (o *RateLimiterConfig) FromStore(r echo.H) *RateLimiterConfig {
	o.On = r.Bool(`On`)
	o.Max = r.Int(`Max`)
	o.Duration = r.Int64(`Duration`)
	o.Prefix = r.String(`Prefix`)
	o.SkipInternalError = r.Bool(`SkipInternalError`)
	o.RedisAddr = r.String(`RedisAddr`)
	o.DBAccountID = r.Uint(`DBAccountID`)
	return o
}

func (o *RateLimiterConfig) Apply(opts *ratelimiter.RateLimiterConfig) *RateLimiterConfig {
	opts.Max = o.Max
	opts.Duration = time.Duration(o.Duration) * time.Second
	opts.Prefix = o.Prefix
	opts.SkipRateLimiterInternalError = o.SkipInternalError
	if o.DBAccountID > 0 {
		m := dbschemaDBMgr.NewNgingDbAccount(defaults.NewMockContext())
		err := m.Get(nil, `id`, o.DBAccountID)
		if err == nil {
			if len(m.Name) == 0 {
				m.Name = `0`
			}
			o.RedisAddr = m.Host
			o.RedisPassword = m.Password
			o.RedisDB = param.AsInt(m.Name)
		}
	}
	if len(o.RedisAddr) > 0 {
		client := NewRedisClient(
			RedisAddr(o.RedisAddr),
			RedisPassword(o.RedisPassword),
			RedisDB(o.RedisDB),
		)
		if err := client.Ping(context.Background()).Err(); err != nil {
			log.Error(color.RedString(`[rateLimiter]`), ` `, err.Error())
		} else {
			opts.Client = client
		}
	}
	return o
}

func RedisAddr(addr string) func(*goredislib.Options) {
	return func(opts *goredislib.Options) {
		opts.Addr = addr
	}
}

func RedisMaxRetries(maxRetries int) func(*goredislib.Options) {
	return func(opts *goredislib.Options) {
		opts.MaxRetries = maxRetries
	}
}

func RedisNetwork(network string) func(*goredislib.Options) {
	return func(opts *goredislib.Options) {
		opts.Network = network
	}
}

func RedisPassword(password string) func(*goredislib.Options) {
	return func(opts *goredislib.Options) {
		opts.Password = password
	}
}

func RedisDB(db int) func(*goredislib.Options) {
	return func(opts *goredislib.Options) {
		opts.DB = db
	}
}

func RedisPoolSize(poolSize int) func(*goredislib.Options) {
	return func(opts *goredislib.Options) {
		opts.PoolSize = poolSize
	}
}

func RedisDialTimeout(timeout time.Duration) func(*goredislib.Options) {
	return func(opts *goredislib.Options) {
		opts.DialTimeout = timeout
	}
}

func RedisReadTimeout(timeout time.Duration) func(*goredislib.Options) {
	return func(opts *goredislib.Options) {
		opts.ReadTimeout = timeout
	}
}

func RedisWriteTimeout(timeout time.Duration) func(*goredislib.Options) {
	return func(opts *goredislib.Options) {
		opts.WriteTimeout = timeout
	}
}

func RedisPoolTimeout(timeout time.Duration) func(*goredislib.Options) {
	return func(opts *goredislib.Options) {
		opts.PoolTimeout = timeout
	}
}

func RedisIdleTimeout(timeout time.Duration) func(*goredislib.Options) {
	return func(opts *goredislib.Options) {
		opts.ConnMaxIdleTime = timeout
	}
}

func RedisTLSConfig(config *tls.Config) func(*goredislib.Options) {
	return func(opts *goredislib.Options) {
		opts.TLSConfig = config
	}
}

func NewRedisClient(settings ...func(*goredislib.Options)) *RedisClient {
	options := goredislib.Options{
		Network: "tcp",
		Addr:    "127.0.0.1:6379",
	}
	for _, option := range settings {
		option(&options)
	}
	c := &RedisClient{
		Client: goredislib.NewClient(&options),
	}
	return c
}

// RedisClient Implements RedisClient for goredislib.Client
type RedisClient struct {
	*goredislib.Client
}

func (c *RedisClient) DeleteKey(ctx context.Context, key string) error {
	return c.Del(ctx, key).Err()
}

func (c *RedisClient) EvalulateSha(ctx context.Context, sha1 string, keys []string, args ...interface{}) (interface{}, error) {
	return c.EvalSha(ctx, sha1, keys, args...).Result()
}

func (c *RedisClient) LuaScriptLoad(ctx context.Context, script string) (string, error) {
	return c.ScriptLoad(ctx, script).Result()
}
