//+build !faker

package http

import (
	uiConfig "github.com/monitoror/monitoror/api/config/usecase"
	coreConfig "github.com/monitoror/monitoror/config"
	"github.com/monitoror/monitoror/monitorables/http/api"
	httpDelivery "github.com/monitoror/monitoror/monitorables/http/api/delivery/http"
	httpModels "github.com/monitoror/monitoror/monitorables/http/api/models"
	httpRepository "github.com/monitoror/monitoror/monitorables/http/api/repository"
	httpUsecase "github.com/monitoror/monitoror/monitorables/http/api/usecase"
	httpCoreConfig "github.com/monitoror/monitoror/monitorables/http/config"
	"github.com/monitoror/monitoror/service/store"
)

type Monitorable struct {
	store *store.Store

	config map[string]*httpCoreConfig.HTTP
}

func NewMonitorable(store *store.Store) *Monitorable {
	monitorable := &Monitorable{}
	monitorable.store = store
	monitorable.config = make(map[string]*httpCoreConfig.HTTP)

	// Load core config from env
	coreConfig.LoadMonitorableConfig(&monitorable.config, httpCoreConfig.Default)

	// Register Monitorable Tile in config manager
	store.UIConfigManager.RegisterTile(api.HTTPStatusTileType, monitorable.GetVariants(), uiConfig.MinimalVersion)
	store.UIConfigManager.RegisterTile(api.HTTPRawTileType, monitorable.GetVariants(), uiConfig.MinimalVersion)
	store.UIConfigManager.RegisterTile(api.HTTPFormattedTileType, monitorable.GetVariants(), uiConfig.MinimalVersion)

	return monitorable
}

func (m *Monitorable) GetVariants() []string { return coreConfig.GetVariantsFromConfig(m.config) }
func (m *Monitorable) IsValid(_ string) bool { return true }

func (m *Monitorable) Enable(variant string) {
	conf := m.config[variant]

	repository := httpRepository.NewHTTPRepository(conf)
	usecase := httpUsecase.NewHTTPUsecase(repository, m.store.CacheStore, m.store.CoreConfig.UpstreamCacheExpiration)
	delivery := httpDelivery.NewHTTPDelivery(usecase)

	// EnableTile route to echo
	httpGroup := m.store.MonitorableRouter.Group("/http", variant)
	routeStatus := httpGroup.GET("/status", delivery.GetHTTPStatus)
	routeRaw := httpGroup.GET("/raw", delivery.GetHTTPRaw)
	routeJSON := httpGroup.GET("/formatted", delivery.GetHTTPFormatted)

	// EnableTile data for config hydration
	m.store.UIConfigManager.EnableTile(api.HTTPStatusTileType, variant, &httpModels.HTTPStatusParams{}, routeStatus.Path, conf.InitialMaxDelay)
	m.store.UIConfigManager.EnableTile(api.HTTPRawTileType, variant, &httpModels.HTTPRawParams{}, routeRaw.Path, conf.InitialMaxDelay)
	m.store.UIConfigManager.EnableTile(api.HTTPFormattedTileType, variant, &httpModels.HTTPFormattedParams{}, routeJSON.Path, conf.InitialMaxDelay)
}
