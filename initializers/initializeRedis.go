package initializers

import (
	"github.com/go-redis/cache/v9"
	"github.com/redis/go-redis/v9"
	"time"
)

var rdb *redis.Client
var Cache *cache.Cache

// Инициализируем кэширование
func InitializeRedis(config *Config) {
	rdb = redis.NewClient(&redis.Options{
		Addr: "localhost:5432",
	})
	Cache = cache.New(&cache.Options{
		Redis: rdb,
		// Cache 10k keys for 1 minute.
		LocalCache: cache.NewTinyLFU(10000, time.Minute),
	})

}
