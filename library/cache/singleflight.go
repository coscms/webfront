package cache

import (
	"context"

	"golang.org/x/sync/singleflight"
)

type SinglefightResult = singleflight.Result

type Singleflighter interface {
	Do(ctx context.Context, key string, fn func() (interface{}, error)) (v interface{}, err error, shared bool)
	DoChan(ctx context.Context, key string, fn func() (interface{}, error)) <-chan SinglefightResult
	Forget(ctx context.Context, key string)
}

type singleflightLock struct {
	mutex TryLocker
}

func (s *singleflightLock) Do(ctx context.Context, key string, fn func() (interface{}, error)) (v interface{}, err error, shared bool) {
	var unlock UnlockFunc
	unlock, err = s.mutex.TryLock(ctx, key)
	if err != nil {
		if err == ErrFailedToAcquireLock {
			err = nil
			shared = true
		}
		return
	}
	defer unlock(ctx)
	v, err = fn()
	return
}

func (s *singleflightLock) DoChan(ctx context.Context, key string, fn func() (interface{}, error)) <-chan SinglefightResult {
	ch := make(chan SinglefightResult, 1)
	var unlock UnlockFunc
	unlock, err := s.mutex.TryLock(ctx, key)
	if err != nil {
		r := SinglefightResult{}
		if err == ErrFailedToAcquireLock {
			r.Shared = true
		} else {
			r.Err = err
		}
		ch <- r
		return ch
	}
	go func() {
		defer unlock(ctx)
		r := SinglefightResult{}
		r.Val, r.Err = fn()
		ch <- r
	}()
	return ch
}

func (s *singleflightLock) Forget(ctx context.Context, key string) {
	s.mutex.Forget(ctx, key)
}

func NewSingleflight(mu TryLocker) Singleflighter {
	return &singleflightLock{mutex: mu}
}

func Singleflight(types ...LockType) Singleflighter {
	return NewSingleflight(GetLocker(types...))
}
