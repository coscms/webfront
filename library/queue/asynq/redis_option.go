package asynq

import (
	"errors"

	"github.com/coscms/webfront/library/cache"
	"github.com/hibiken/asynq"
)

func ParseRedisURI(connURI string) (asynq.RedisConnOpt, error) {
	return asynq.ParseRedisURI(connURI)
}

func RedisOptFromCache() (asynq.RedisConnOpt, error) {
	cfg := cache.RedisOptions()
	if cfg == nil {
		return nil, errors.New(`supported`)
	}
	return asynq.RedisClientOpt{
		Network: cfg.Network,
		Addr:    cfg.Addr,
		//Username:     cfg.Username,
		Password:     cfg.Password,
		DB:           cfg.DB,
		DialTimeout:  cfg.DialTimeout,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		PoolSize:     cfg.PoolSize,
		TLSConfig:    cfg.TLSConfig,
	}, nil
}

type RedisClientOpt = asynq.RedisClientOpt
type RedisFailoverClientOpt = asynq.RedisFailoverClientOpt
type RedisClusterClientOpt = asynq.RedisClusterClientOpt
