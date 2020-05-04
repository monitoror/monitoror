package store

import (
	"github.com/monitoror/monitoror/cli"
	coreConfig "github.com/monitoror/monitoror/config"
	"github.com/monitoror/monitoror/service/middlewares"
	"github.com/monitoror/monitoror/service/registry"
	"github.com/monitoror/monitoror/service/router"

	"github.com/jsdidierlaurent/echo-middleware/cache"
)

type (
	// Store is used to share Data in every monitorable
	Store struct {
		// CLI helper
		Cli cli.CLI

		// Global CoreConfig
		CoreConfig *coreConfig.CoreConfig

		// CacheStore for every memory persistent data
		CacheStore cache.Store
		// CacheMiddleware using CacheStore to return cached data
		CacheMiddleware *middlewares.CacheMiddleware

		// MonitorableRouter helper wrapping echo Router monitorable
		MonitorableRouter router.MonitorableRouter

		// Registry used to register Tile for verify / hydrate
		Registry registry.Registry
	}
)
