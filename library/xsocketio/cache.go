package xsocketio

import (
	"sync"

	"github.com/admpub/log"
	"github.com/coscms/webfront/library/cache"
	socketio "github.com/googollee/go-socket.io"
	esi "github.com/webx-top/echo-socket.io"
)

var (
	instances = map[string]*esi.Wrapper{}
	mu        = &sync.RWMutex{}
)

func SocketIO(namespace string, cfg *Config) *esi.Wrapper {
	mu.RLock()
	v, y := instances[namespace]
	mu.RUnlock()
	if y {
		return v
	}

	v = socketIOWrapper(namespace)
	if cfg.EnableRedis {
		redisCfg := cache.RedisOptions()
		if redisCfg != nil {
			adpCfg := &socketio.RedisAdapterOptions{
				Addr:     redisCfg.Addr,
				Prefix:   `SOCKETIO`,
				Network:  redisCfg.Network,
				Password: redisCfg.Password,
				DB:       redisCfg.DB,
			}
			if cfg.RedisDB > 0 {
				adpCfg.DB = cfg.RedisDB
			}
			v.Server.Adapter(adpCfg)
			log.Okayf(`socket.io enable redis adapter`)
		}
	}
	v.Serve()
	mu.Lock()
	instances[namespace] = v
	mu.Unlock()
	return v
}

func Close(namespace string) bool {
	mu.RLock()
	v, y := instances[namespace]
	mu.RUnlock()
	if y {
		v.Close()
		mu.Lock()
		delete(instances, namespace)
		mu.Unlock()
	}

	return y
}

func CloseAll() {
	mu.Lock()
	for namespace, instance := range instances {
		instance.Close()
		delete(instances, namespace)
	}
	mu.Unlock()
}
