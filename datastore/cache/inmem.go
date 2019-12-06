package cache

import (
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/spf13/pflag"
)

const InMemoryExpirationTimeout = 5 * time.Second
const InMemoryPurgeTimeout = 5 * time.Second

func NewInMemoryCache() *cache.Cache {
	c := cache.New(InMemoryExpirationTimeout, InMemoryPurgeTimeout)
	return c
}

func Flags() *pflag.FlagSet {
	fs := pflag.NewFlagSet("inmem-cache", pflag.ContinueOnError)
	fs.Duration("inmem-cache-expiration", InMemoryExpirationTimeout, "default expiration timeout for cache entries")
	fs.Duration("inmem-cache-purge", InMemoryExpirationTimeout, "default purge timeout for cache entries")
	return fs
}

