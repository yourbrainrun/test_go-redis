package helpers

import (
	"context"
	"test_go-redis/drivers"
	"time"
)

type RedisLock struct {
	ReleaseDuration time.Duration
}

func (lock *RedisLock) Lock(key string) bool {
	redisDb := drivers.GetRedisClient()
	ctx := context.Background()
	res, err := redisDb.SetNX(ctx, key, 1, lock.ReleaseDuration).Result()
	if err != nil || res == false {
		return false
	}
	return true
}

func (lock *RedisLock) Unlock(name string) bool {
	ctx := context.Background()
	redisDb := drivers.GetRedisClient()

	_, err := redisDb.Del(ctx, name).Result()
	if err != nil {
		return false
	}
	return true
}
