package redis

import (
	"net/url"
	"strconv"
	"strings"

	"github.com/RichardKnop/machinery/v2/config"
	"github.com/coscms/webfront/library/queue/machinery"

	redisBackend "github.com/RichardKnop/machinery/v2/backends/redis"
	redisBroker "github.com/RichardKnop/machinery/v2/brokers/redis"
	redisLock "github.com/RichardKnop/machinery/v2/locks/redis"

	backendsiface "github.com/RichardKnop/machinery/v2/backends/iface"
	brokersiface "github.com/RichardKnop/machinery/v2/brokers/iface"
	locksiface "github.com/RichardKnop/machinery/v2/locks/iface"
)

func init() {
	machinery.RegisterBackendEngine(`redis`, initBackend)
	machinery.RegisterBrokerEngine(`redis`, initBroker)
	machinery.RegisterLockEngine(`redis`, initLock)
}

func initBackend(u *url.URL, cfg *config.Config) (backendsiface.Backend, error) {
	var db int
	var err error
	sdb := u.Query().Get(`db`)
	if len(sdb) > 0 {
		db, err = strconv.Atoi(sdb)
		if err != nil {
			return nil, err
		}
	} else if after, found := strings.CutPrefix(u.Path, `/`); found && len(after) > 0 {
		db, err = strconv.Atoi(after)
		if err != nil {
			return nil, err
		}
	}
	if u.User != nil && len(u.User.Username()) > 0 {
		password, _ := u.User.Password()
		return redisBackend.New(cfg, u.Host, u.User.Username(), password, u.Query().Get(`socketPath`), db), err
	}
	return redisBackend.NewGR(cfg, strings.Split(u.Host, `,`), db), err
}
func initBroker(u *url.URL, cfg *config.Config) (brokersiface.Broker, error) {
	var db int
	var err error
	sdb := u.Query().Get(`db`)
	if len(sdb) > 0 {
		db, err = strconv.Atoi(sdb)
		if err != nil {
			return nil, err
		}
	} else if after, found := strings.CutPrefix(u.Path, `/`); found && len(after) > 0 {
		db, err = strconv.Atoi(after)
		if err != nil {
			return nil, err
		}
	}
	if u.User != nil && len(u.User.Username()) > 0 {
		password, _ := u.User.Password()
		return redisBroker.New(cfg, u.Host, u.User.Username(), password, u.Query().Get(`socketPath`), db), err
	}
	return redisBroker.NewGR(cfg, strings.Split(u.Host, `,`), db), err
}
func initLock(u *url.URL, cfg *config.Config) (locksiface.Lock, error) {
	var db int
	var err error
	sdb := u.Query().Get(`db`)
	if len(sdb) > 0 {
		db, err = strconv.Atoi(sdb)
		if err != nil {
			return nil, err
		}
	} else if after, found := strings.CutPrefix(u.Path, `/`); found && len(after) > 0 {
		db, err = strconv.Atoi(after)
		if err != nil {
			return nil, err
		}
	}
	var retries int
	sRetries := u.Query().Get(`retries`)
	if len(sRetries) > 0 {
		retries, err = strconv.Atoi(sRetries)
		if err != nil {
			return nil, err
		}
	}
	return redisLock.New(cfg, strings.Split(u.Host, `,`), db, retries), err
}
