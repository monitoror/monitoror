package middlewares

import (
	"time"

	"github.com/jsdidierlaurent/monitoror/config"

	"github.com/jsdidierlaurent/echo-middleware/cache"
	"github.com/labstack/echo/v4"
)

/**
* Cache Middleware for monitoror
*
* We need two types of Cache.
* - UpstreamCache : serves as a circuit breaker to answer before executing the request. (By default, short TTL)
* - DownstreamCache : serves as backup to answer the old result in case of service timeout. (By default, long TTL)
*
* UpstreamCache must be implemented on some routes only (and with variable expiration).
* He is implemented as a decorator on the handler of each route
*
* DownstreamCache should be used instead of a timeout response.
* So we look at the cache in the global error handler (see handlers/error.go)
*
* To fill both store at the same time, I implemented a store wrapper that performs every actions on both store
 */

const (
	DownstreamStoreContextKey = "jsdidierlaurent.monitoror.downstreamStore"
	CachePrefix               = "jsdidierlaurent.monitoror.cache"

	DownstreamCacheHeader = "Timeout-Recover"
)

type (
	CacheMiddleware struct {
		store responsesStore
	}

	//responsesStore implement cache.Store to provide it to CacheMiddleware
	responsesStore struct {
		UpstreamStore   cache.Store
		DownstreamStore cache.Store
	}
)

//NewCacheMiddleware used config to instantiate CacheMiddleware
func NewCacheMiddleware(config *config.Config) *CacheMiddleware {
	store := responsesStore{
		UpstreamStore: cache.NewGoCacheStore(
			time.Second*time.Duration(config.UpstreamCache.Expire),
			time.Second*time.Duration(config.UpstreamCache.CleanupInterval),
		),
		DownstreamStore: cache.NewGoCacheStore(
			time.Second*time.Duration(config.DownstreamCache.Expire),
			time.Second*time.Duration(config.DownstreamCache.CleanupInterval),
		),
	}

	return &CacheMiddleware{store: store}
}

//==============================================================================
// UPSTREAM MIDDLEWARE
//==============================================================================

// UpstreamCache return the cached response if he finds it in the store. (Decorator Handlers)
func (cm *CacheMiddleware) UpstreamCacheHandler(handle echo.HandlerFunc) echo.HandlerFunc {
	return cache.CacheHandlerWithConfig(cache.CacheMiddlewareConfig{
		Store:     &cm.store,
		KeyPrefix: CachePrefix,
	}, handle)
}

//UpstreamCacheWithExpiration return the cached response if he finds it in the store. (Decorator Handlers)
func (cm *CacheMiddleware) UpstreamCacheHandlerWithExpiration(expire time.Duration, handle echo.HandlerFunc) echo.HandlerFunc {
	return cache.CacheHandlerWithConfig(cache.CacheMiddlewareConfig{
		Store:     &cm.store,
		KeyPrefix: CachePrefix,
		Expire:    expire,
	}, handle)
}

//==============================================================================
// DOWNSTREAM MIDDLEWARE
//==============================================================================

// DownstreamStoreMiddleware Provide Downstream Store to all route. Used when route return timeout error
func (cm *CacheMiddleware) DownstreamStoreMiddleware() echo.MiddlewareFunc {
	config := cache.StoreMiddlewareConfig{
		Store:      cm.store.DownstreamStore,
		ContextKey: DownstreamStoreContextKey,
	}
	return cache.StoreMiddlewareWithConfig(config)
}

//==============================================================================
// ResponsesStore methods (implementation of cache.Store)
//==============================================================================
func (c *responsesStore) Get(key string, value interface{}) error {
	return c.UpstreamStore.Get(key, value)
}

func (c *responsesStore) Set(key string, val interface{}, expires time.Duration) (err error) {
	err = c.UpstreamStore.Set(key, val, expires)
	_ = c.DownstreamStore.Set(key, val, cache.DEFAULT)
	return
}

func (c *responsesStore) Add(key string, value interface{}, expires time.Duration) error {
	panic("unimplemented")
}

func (c *responsesStore) Replace(key string, value interface{}, expires time.Duration) error {
	panic("unimplemented")
}

func (c *responsesStore) Delete(key string) error {
	panic("unimplemented")
}

func (c *responsesStore) Increment(key string, n uint64) (uint64, error) {
	panic("unimplemented")
}

func (c *responsesStore) Decrement(key string, n uint64) (uint64, error) {
	panic("unimplemented")
}

func (c *responsesStore) Flush() error {
	panic("unimplemented")

}
