package middlewares

import (
	"strings"
	"time"

	"github.com/monitoror/monitoror/models"

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
	DownstreamStoreContextKey = "monitoror.downstreamStore"
	DownstreamCacheHeader     = "Timeout-Recover"

	TempKeyPrefix = "T"
)

type (
	CacheMiddleware struct {
		store                       cache.Store
		downstreamDefaultExpiration time.Duration
		upstreamDefaultExpiration   time.Duration
	}

	// Wrapper for setting value in store with 2 keys for timeout
	upstreamStore struct {
		store                       cache.Store
		downstreamDefaultExpiration time.Duration
	}
)

// NewCacheMiddleware used config to instantiate CacheMiddleware
func NewCacheMiddleware(store cache.Store, downstreamDefaultExpiration, upstreamDefaultExpiration time.Duration) *CacheMiddleware {
	return &CacheMiddleware{store, downstreamDefaultExpiration, upstreamDefaultExpiration}
}

//==============================================================================
// UPSTREAM MIDDLEWARE
//==============================================================================

// UpstreamCache return the cached response if he finds it in the store. (Decorator Handlers)
func (cm *CacheMiddleware) UpstreamCacheHandler(handle echo.HandlerFunc) echo.HandlerFunc {
	return cache.CacheHandlerWithConfig(cache.CacheMiddlewareConfig{
		Store:     &upstreamStore{cm.store, cm.downstreamDefaultExpiration},
		KeyPrefix: TempKeyPrefix, // hack for use Sprintf inside set methode
		Expire:    cm.upstreamDefaultExpiration,
	}, handle)
}

//UpstreamCacheWithExpiration return the cached response if he finds it in the store. (Decorator Handlers)
func (cm *CacheMiddleware) UpstreamCacheHandlerWithExpiration(expire time.Duration, handle echo.HandlerFunc) echo.HandlerFunc {
	return cache.CacheHandlerWithConfig(cache.CacheMiddlewareConfig{
		Store:     &upstreamStore{cm.store, cm.downstreamDefaultExpiration},
		KeyPrefix: TempKeyPrefix, // hack for use Sprintf inside set methode
		Expire:    expire,
	}, handle)
}

//==============================================================================
// DOWNSTREAM MIDDLEWARE
//==============================================================================

// DownstreamStoreMiddleware Provide Downstream Store to all route. Used when route return timeout error
func (cm *CacheMiddleware) DownstreamStoreMiddleware() echo.MiddlewareFunc {
	config := cache.StoreMiddlewareConfig{
		Store:      cm.store,
		ContextKey: DownstreamStoreContextKey,
	}
	return cache.StoreMiddlewareWithConfig(config)
}

//==============================================================================
// ResponsesStore methods (implementation of cache.Store)
//==============================================================================
func (c *upstreamStore) Get(key string, value interface{}) error {
	return c.store.Get(getCustomKey(key, models.UpstreamStoreKeyPrefix), value)
}

func (c *upstreamStore) Set(key string, val interface{}, expires time.Duration) (err error) {
	err = c.store.Set(getCustomKey(key, models.UpstreamStoreKeyPrefix), val, expires)
	_ = c.store.Set(getCustomKey(key, models.DownstreamStoreKeyPrefix), val, c.downstreamDefaultExpiration)
	return
}

func (c *upstreamStore) Add(key string, value interface{}, expires time.Duration) error {
	panic("unimplemented")
}

func (c *upstreamStore) Replace(key string, value interface{}, expires time.Duration) error {
	panic("unimplemented")
}

func (c *upstreamStore) Delete(key string) error {
	panic("unimplemented")
}

func (c *upstreamStore) Increment(key string, n uint64) (uint64, error) {
	panic("unimplemented")
}

func (c *upstreamStore) Decrement(key string, n uint64) (uint64, error) {
	panic("unimplemented")
}

func (c *upstreamStore) Flush() error {
	panic("unimplemented")
}

func getCustomKey(key, keyPrefix string) string {
	return strings.Replace(key, TempKeyPrefix, keyPrefix, 1)
}
