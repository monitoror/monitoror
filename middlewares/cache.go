package middlewares

import (
	"time"

	"github.com/labstack/echo/v4"

	"github.com/jsdidierlaurent/echo-middleware/cache"
	"github.com/jsdidierlaurent/monitowall/config"
)

type (
	CacheMiddleware struct {
		ForeverCacheConfig cache.CacheMiddlewareConfig
		DefaultCacheConfig cache.CacheMiddlewareConfig
	}
)

//NewCacheMiddleware create common store for 2 types of cache middleware. One without expiration and one with default expiration (set in config)
func NewCacheMiddleware(config *config.Config) *CacheMiddleware {
	store := cache.NewGoCacheStore(
		time.Second*time.Duration(config.Cache.Duration),
		time.Second*time.Duration(config.Cache.CleanupInterval),
	)

	return &CacheMiddleware{
		ForeverCacheConfig: cache.CacheMiddlewareConfig{
			Store:  store,
			Expire: cache.FOREVER,
		},
		DefaultCacheConfig: cache.CacheMiddlewareConfig{
			Store: store,
		},
	}
}

//ForeverCache Cache middleware without expiration
func (cm *CacheMiddleware) ForeverCache(handle echo.HandlerFunc) echo.HandlerFunc {
	return cache.CacheHandlerWithConfig(cm.ForeverCacheConfig, handle)
}

//DefaultCache Cache middleware with default expiration (set in config)
func (cm *CacheMiddleware) DefaultCache(handle echo.HandlerFunc) echo.HandlerFunc {
	return cache.CacheHandlerWithConfig(cm.DefaultCacheConfig, handle)
}
