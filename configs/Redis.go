package configs

import (
	"os"
	"strconv"
	"sync"
)

type Redis struct {
	Host     string
	Password string
	Port     string
	DB       int
	PoolSize int
}

var RedisConfigInstance *Redis
var RedisConfigOnce sync.Once

func GetRedisConfig() *Redis {
	RedisConfigOnce.Do(func() {
		db, err := strconv.ParseInt(os.Getenv("REDIS_DB"), 10, 8)
		if err != nil {
			panic("redis db is invalid")
		}
		poolSize, err := strconv.ParseInt(os.Getenv("REDIS_POOL_SIZE"), 10, 8)
		if err != nil {
			poolSize = 0
		}
		RedisConfigInstance = &Redis{
			Host:     os.Getenv("REDIS_HOST"),
			Password: os.Getenv("REDIS_PASSWORD"),
			Port:     os.Getenv("REDIS_PORT"),
			DB:       int(db),
			PoolSize: int(poolSize),
		}
	})
	return RedisConfigInstance
}
