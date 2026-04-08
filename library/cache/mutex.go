package cache

import (
	"context"
	"errors"
	"sync/atomic"
	"time"
)

var ErrFailedToAcquireLock = errors.New("failed to acquire lock")

var (
	defaultLockType int32 = int32(LockTypeMemory)
	tryLockers            = map[LockType]TryLocker{
		LockTypeMemory: &mutexMemory{},
		LockTypeRedis:  NewMutexRedis(2 * time.Minute),
	}
)

type LockType int32

const (
	LockTypeMemory LockType = iota
	LockTypeRedis
)

func DefaultLockType() LockType {
	v := atomic.LoadInt32(&defaultLockType)
	return LockType(v)
}

func SetDefaultLockType(lockType LockType) {
	atomic.AddInt32(&defaultLockType, int32(lockType))
}

func RegisterTryLocker(t LockType, fn TryLocker) {
	tryLockers[t] = fn
}

type UnlockFunc func(context.Context) error
type TryLocker interface {
	Lock(ctx context.Context, key string) (unlock UnlockFunc, err error)
	TryLock(ctx context.Context, key string) (unlock UnlockFunc, err error)
	TryLockWithTimeout(ctx context.Context, key string, timeout time.Duration) (unlock UnlockFunc, err error)
	Forget(ctx context.Context, key string)
}

func GetLocker(types ...LockType) TryLocker {
	var t LockType
	if len(types) > 0 {
		t = types[0]
	} else {
		t = DefaultLockType()
	}
	if tryLocker, ok := tryLockers[t]; ok {
		return tryLocker
	}
	return tryLockers[LockTypeMemory]
}

func Lock(ctx context.Context, key string, types ...LockType) (unlock UnlockFunc, err error) {
	return GetLocker(types...).Lock(ctx, key)
}

func TryLock(ctx context.Context, key string, types ...LockType) (unlock UnlockFunc, err error) {
	return GetLocker(types...).TryLock(ctx, key)
}

func TryLockWithTimeout(ctx context.Context, key string, timeout time.Duration, types ...LockType) (unlock UnlockFunc, err error) {
	return GetLocker(types...).TryLockWithTimeout(ctx, key, timeout)
}

func ForgetLock(ctx context.Context, key string, types ...LockType) {
	GetLocker(types...).Forget(ctx, key)
}
