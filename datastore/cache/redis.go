package cache

import (
	"context"
	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/extra/redisotel/v8"
	"github.com/go-redis/redis/v8"
	"github.com/marcoEgger/genki/config"
	"github.com/marcoEgger/genki/logger"
	"github.com/spf13/pflag"
	"time"
)

const (
	CacheUrl      = "cache-url"
	CacheDatabase = "cache-database"
)

//goland:noinspection GoUnusedExportedFunction
func NewRedisCache() *cache.Cache {
	redisClient := redis.NewClient(&redis.Options{
		Addr: config.GetString(CacheUrl),
		DB:   config.GetInt(CacheDatabase),
		OnConnect: func(ctx context.Context, cn *redis.Conn) error {
			logger.Debugf("connected to redis")
			return nil
		},
	})
	redisClient.AddHook(redisotel.NewTracingHook())
	return cache.New(&cache.Options{
		Redis:      redisClient,
		LocalCache: cache.NewTinyLFU(5000, time.Minute),
	})
}

//goland:noinspection GoUnusedExportedFunction
func Flags() *pflag.FlagSet {
	fs := pflag.NewFlagSet("redis-cache", pflag.ContinueOnError)
	fs.String(CacheUrl, "cache-1-redis-master.cache.svc.cluster.local:6379", "the cache URL including port")
	fs.Int(CacheDatabase, 0, "the redis cache database number including port")
	return fs
}
