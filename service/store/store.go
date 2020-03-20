package store

import (
	uiConfig "github.com/monitoror/monitoror/api/config"
	coreConfig "github.com/monitoror/monitoror/config"
	"github.com/monitoror/monitoror/service/middlewares"
	"github.com/monitoror/monitoror/service/router"

	"github.com/jsdidierlaurent/echo-middleware/cache"
)

type (
	// Store is used to share Data in every monitorable
	Store struct {
		// Global CoreConfig
		CoreConfig *coreConfig.Config

		// CacheStore for every memory persistent data
		CacheStore cache.Store
		// MidCacheMiddlewaredleware using CacheStore to return cached data
		CacheMiddleware *middlewares.CacheMiddleware

		// MonitorableRouter helper wrapping echo Router monitorable
		MonitorableRouter router.MonitorableRouter

		// MonitorableConfigManager used to register Tile for verify / hydrate
		UIConfigManager uiConfig.Manager
	}
)
