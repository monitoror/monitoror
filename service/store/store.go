package store

import (
	uiConfig "github.com/monitoror/monitoror/api/config"
	"github.com/monitoror/monitoror/cli"
	coreConfig "github.com/monitoror/monitoror/config"
	"github.com/monitoror/monitoror/service/middlewares"
	"github.com/monitoror/monitoror/service/router"

	"github.com/jsdidierlaurent/echo-middleware/cache"
)

type (
	// Store is used to share Data in every monitorable
	Store struct {
		// CLI helper
		Cli cli.CLI

		// Global CoreConfig
		CoreConfig *coreConfig.Config

		// CacheStore for every memory persistent data
		CacheStore cache.Store
		// CacheMiddleware using CacheStore to return cached data
		CacheMiddleware *middlewares.CacheMiddleware

		// MonitorableRouter helper wrapping echo Router monitorable
		MonitorableRouter router.MonitorableRouter

		// MonitorableConfigManager used to register Tile for verify / hydrate
		TileSettingManager uiConfig.TileSettingManager
	}
)
