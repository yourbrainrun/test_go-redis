package drivers

import (
	"github.com/go-redis/redis/v8"
	"sync"
	"test_go-redis/configs"
)

var redisInstance *redis.Client
var redisOnce sync.Once

func GetRedisClient() *redis.Client {
	redisOnce.Do(func() {
		config := configs.GetRedisConfig()
		redisInstance = redis.NewClient(&redis.Options{
			Addr:               config.Host + ":" + config.Port,
			Password:           config.Password,
			DB:                 config.DB,
			PoolSize:           config.PoolSize,
			MinIdleConns:       0,
			MaxConnAge:         0,
			PoolTimeout:        0,
			IdleTimeout:        0,
			IdleCheckFrequency: 0,
			TLSConfig:          nil,
			Limiter:            nil,
		})
	})
	return redisInstance
}
