package cache

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/admpub/cache"
	"github.com/stretchr/testify/assert"
)

func init() {
	CacheNew(context.Background(), cache.Options{
		Adapter:       `redis`,
		AdapterConfig: `network=tcp,addr=127.0.0.1:6379,password=,db=0,pool_size=100,idle_timeout=180,hset_name=Cache,prefix=cache:`,
	}, `default`)
}

func TestMutex(t *testing.T) {
	unlock, err := TryLock(context.Background(), `test`)
	assert.NoError(t, err)
	err = unlock(context.Background())
	assert.NoError(t, err)
}

func TestMutex2(t *testing.T) {
	unlock, err := TryLock(context.Background(), `test`)
	assert.NoError(t, err)
	_, err2 := TryLock(context.Background(), `test`)
	assert.Equal(t, ErrFailedToAcquireLock, err2)
	_, err3 := TryLock(context.Background(), `test`)
	assert.Equal(t, ErrFailedToAcquireLock, err3)
	err = unlock(context.Background())
	assert.NoError(t, err)

	TestMutex(t)
}

func TestMutexRedis(t *testing.T) {
	unlock, err := TryLock(context.Background(), `test`, LockTypeRedis)
	if err == nil {
		err = unlock(context.Background())
		assert.NoError(t, err)
	}
}

func TestMutexRedisLock(t *testing.T) {
	var n int
	var wg sync.WaitGroup
	j := 10
	for i := 0; i < j; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			unlock, err := Lock(context.Background(), `test`, LockTypeRedis)
			assert.NoError(t, err)
			time.Sleep(1 * time.Second)
			fmt.Printf("TestMutexRedisLock ================>%d\n", i)
			n++
			err = unlock(context.Background())
			assert.NoError(t, err)
		}(i)
	}
	wg.Wait()
	assert.Equal(t, j, n)
}

func TestMutexRedis2(t *testing.T) {
	unlock, err := TryLock(context.Background(), `test`, LockTypeRedis)
	assert.NoError(t, err)
	_, err2 := TryLock(context.Background(), `test`, LockTypeRedis)
	assert.Equal(t, ErrFailedToAcquireLock, err2)
	_, err3 := TryLock(context.Background(), `test`, LockTypeRedis)
	assert.Equal(t, ErrFailedToAcquireLock, err3)
	err = unlock(context.Background())
	assert.NoError(t, err)

	TestMutexRedis(t)
}

func BenchmarkMutex(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			unlock, err := TryLock(context.Background(), `test`)
			if err != nil {
				if err == ErrFailedToAcquireLock {
					continue
				}
				panic(err)
			}
			err = unlock(context.Background())
			if err != nil {
				panic(err)
			}
		}
	})
}

func BenchmarkMutexRedis(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			unlock, err := TryLock(context.Background(), `test`, LockTypeRedis)
			if err != nil {
				if err == ErrFailedToAcquireLock {
					continue
				}
				panic(err)
			}
			err = unlock(context.Background())
			if err != nil {
				panic(err)
			}
		}
	})
}
