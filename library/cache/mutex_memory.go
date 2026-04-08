package cache

import (
	"context"
	"time"

	lock "github.com/admpub/go-lock"
)

var (
	mutexGroup = lock.NewGroup(nil)
)

type mutexMemory struct{}

func (*mutexMemory) Lock(ctx context.Context, key string) (unlock UnlockFunc, err error) {
	mutexGroup.Lock(key)
	unlock = func(ctx context.Context) error {
		mutexGroup.UnlockAndFree(key)
		return nil
	}
	return
}

func (*mutexMemory) TryLock(ctx context.Context, key string) (unlock UnlockFunc, err error) {
	if !mutexGroup.TryLock(key) {
		err = ErrFailedToAcquireLock
		return
	}
	unlock = func(ctx context.Context) error {
		mutexGroup.UnlockAndFree(key)
		return nil
	}
	return
}

func (*mutexMemory) TryLockWithTimeout(ctx context.Context, key string, timeout time.Duration) (unlock UnlockFunc, err error) {
	if !mutexGroup.TryLockWithTimeout(key, timeout) {
		err = ErrFailedToAcquireLock
		return
	}
	unlock = func(ctx context.Context) error {
		mutexGroup.UnlockAndFree(key)
		return nil
	}
	return
}

func (*mutexMemory) Forget(ctx context.Context, key string) {
	mutexGroup.UnlockAndFree(key)
}
