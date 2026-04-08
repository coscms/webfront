package cache

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/admpub/once"
	"github.com/admpub/redsync/v4"
	goredis "github.com/admpub/redsync/v4/redis/goredis/v9"
	goredislib "github.com/redis/go-redis/v9"
)

var (
	redsyncClient *redsync.Redsync
	redsyncOnce   once.Once
)

func resetRedsync() {
	redsyncOnce.Reset()
}

func initRedsync() {
	client, _ := Cache(cacheRootContext, `locker`).Client().(*goredislib.Client)
	if client == nil {
		client = RedisClient()
	}
	pool := goredis.NewPool(client)
	redsyncClient = redsync.New(pool)
}

func onceInitRedsync() {
	initRedsync()
}

func RedsyncClient() *redsync.Redsync {
	redsyncOnce.Do(onceInitRedsync)

	return redsyncClient
}

// RedisMutex 分布式锁
// example:
// mutex := RedisMutex(`goods_1`)
// err = mutex.Lock(ctx)
//
//	if err != nil {
//		panic(err)
//	}
//
// mutex.Unlock(ctx)
func RedisMutex(key string, options ...redsync.Option) *redsync.Mutex {
	return RedsyncClient().NewMutex(key, options...)
}

func NewMutexRedis(maxLockDuration time.Duration) TryLocker {
	if maxLockDuration <= 0 {
		maxLockDuration = time.Minute
	}
	return &mutexRedis{
		maxLockDuration: maxLockDuration,
	}
}

type mutexRedis struct {
	maxLockDuration time.Duration
}

func (r *mutexRedis) Lock(ctx context.Context, key string) (unlock UnlockFunc, err error) {
	delay := 100 * time.Millisecond
	m := RedisMutex(key,
		redsync.WithExpiry(r.maxLockDuration),
		redsync.WithTries(1000),
		redsync.WithRetryDelayFunc(func(tries int) time.Duration {
			return delay * time.Duration(tries)
		}),
		redsync.WithRetryDelay(delay),
	)

	err = m.LockContext(ctx)
	if err != nil {
		if err == redsync.ErrFailed {
			err = ErrFailedToAcquireLock
		}
		return
	}
	unlock = func(ctx context.Context) error {
		ok, err := m.UnlockContext(ctx)
		if !ok || err != nil {
			return fmt.Errorf("unlock unsuccessful: %w", err)
		}
		return nil
	}
	return
}

func isRedsyncErrTaken(err error) bool {
	_, ok := errors.AsType[*redsync.ErrTaken](err)
	if ok {
		return ok
	}
	_, ok = errors.AsType[*redsync.ErrNodeTaken](err)
	return ok
}

func (r *mutexRedis) TryLock(ctx context.Context, key string) (unlock UnlockFunc, err error) {
	m := RedisMutex(key,
		redsync.WithExpiry(r.maxLockDuration),
		redsync.WithTries(1),
		redsync.WithRetryDelay(50*time.Millisecond),
	)
	err = m.LockContext(ctx)
	if err != nil {
		if isRedsyncErrTaken(err) {
			err = ErrFailedToAcquireLock
		}
		return
	}
	done := make(chan struct{})
	go func() {
		ticker := time.NewTicker(r.maxLockDuration / 3)
		defer ticker.Stop()
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				m.ExtendContext(ctx)
			}
		}
	}()

	unlock = func(ctx context.Context) error {
		close(done)
		ok, err := m.UnlockContext(ctx)
		if !ok || err != nil {
			return fmt.Errorf("unlock unsuccessful: %w", err)
		}
		return nil
	}
	return
}

func (r *mutexRedis) TryLockWithTimeout(ctx context.Context, key string, maxLockDuration time.Duration) (unlock UnlockFunc, err error) {
	return r.tryLockWithContext(ctx, key, maxLockDuration)
}

func (r *mutexRedis) tryLockWithContext(ctx context.Context, key string, maxLockDuration time.Duration) (unlock UnlockFunc, err error) {
	m := RedisMutex(key,
		redsync.WithExpiry(maxLockDuration),
		redsync.WithTries(1),
		redsync.WithRetryDelay(50*time.Millisecond),
	)
	err = m.LockContext(ctx)
	if err != nil {
		if isRedsyncErrTaken(err) {
			err = ErrFailedToAcquireLock
		}
		return
	}
	done := make(chan struct{})
	go func() {
		ticker := time.NewTicker(maxLockDuration / 3)
		defer ticker.Stop()
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				m.ExtendContext(ctx)
			}
		}
	}()

	unlock = func(ctx context.Context) error {
		close(done)
		ok, err := m.UnlockContext(ctx)
		if !ok || err != nil {
			return fmt.Errorf("unlock unsuccessful: %w", err)
		}
		return nil
	}
	return
}

func (*mutexRedis) Forget(ctx context.Context, key string) {
	RedisMutex(key).UnlockContext(ctx)
}
