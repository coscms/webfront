package xsocketio

import (
	"github.com/coscms/webcore/library/common"
	"github.com/coscms/webcore/library/config"
	"github.com/webx-top/echo"
)

func GetConfig() *Config {
	cfg, ok := common.Setting(`socketio`).Get(`config`).(*Config)
	if ok {
		return cfg
	}
	cfg = NewConfig()
	cfg.FromStore(config.FromFile().Extend.GetStore(`socketio`))
	return cfg
}

func NewConfig() *Config {
	return &Config{}
}

type Config struct {
	EnableRedis   bool   `json:"enableRedis"`
	RedisDB       int    `json:"redisDB"`
	RedisAddr     string `json:"redisAddr"`
	RedisPrefix   string `json:"redisPrefix"`
	RedisNetwork  string `json:"redisNetwork"`
	RedisPassword string `json:"redisPassword"`
}

func (c *Config) FromStore(v echo.H) *Config {
	c.EnableRedis = v.Bool(`enableRedis`)
	c.RedisDB = v.Int(`redisDB`)
	c.RedisAddr = v.String(`redisAddr`)
	c.RedisPrefix = v.String(`redisPrefix`)
	c.RedisNetwork = v.String(`redisNetwork`)
	c.RedisPassword = v.String(`redisPassword`)
	return c
}
