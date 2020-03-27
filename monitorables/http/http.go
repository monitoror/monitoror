//+build !faker

package http

import (
	uiConfig "github.com/monitoror/monitoror/api/config/usecase"
	pkgMonitorable "github.com/monitoror/monitoror/internal/pkg/monitorable"
	coreModels "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorables/http/api"
	httpDelivery "github.com/monitoror/monitoror/monitorables/http/api/delivery/http"
	httpModels "github.com/monitoror/monitoror/monitorables/http/api/models"
	httpRepository "github.com/monitoror/monitoror/monitorables/http/api/repository"
	httpUsecase "github.com/monitoror/monitoror/monitorables/http/api/usecase"
	httpConfig "github.com/monitoror/monitoror/monitorables/http/config"
	"github.com/monitoror/monitoror/service/store"
)

type Monitorable struct {
	store *store.Store

	config map[coreModels.VariantName]*httpConfig.HTTP
}

func NewMonitorable(store *store.Store) *Monitorable {
	monitorable := &Monitorable{}
	monitorable.store = store
	monitorable.config = make(map[coreModels.VariantName]*httpConfig.HTTP)

	// Load core config from env
	pkgMonitorable.LoadConfig(&monitorable.config, httpConfig.Default)

	// Register Monitorable Tile in config manager
	store.UIConfigManager.RegisterTile(api.HTTPStatusTileType, monitorable.GetVariants(), uiConfig.MinimalVersion)
	store.UIConfigManager.RegisterTile(api.HTTPRawTileType, monitorable.GetVariants(), uiConfig.MinimalVersion)
	store.UIConfigManager.RegisterTile(api.HTTPFormattedTileType, monitorable.GetVariants(), uiConfig.MinimalVersion)

	return monitorable
}

func (m *Monitorable) GetDisplayName() string {
	return "HTTP"
}

func (m *Monitorable) GetVariants() []coreModels.VariantName {
	return pkgMonitorable.GetVariants(m.config)
}

func (m *Monitorable) Validate(_ coreModels.VariantName) (bool, error) {
	return true, nil
}

func (m *Monitorable) Enable(variant coreModels.VariantName) {
	conf := m.config[variant]

	repository := httpRepository.NewHTTPRepository(conf)
	usecase := httpUsecase.NewHTTPUsecase(repository, m.store.CacheStore, m.store.CoreConfig.UpstreamCacheExpiration)
	delivery := httpDelivery.NewHTTPDelivery(usecase)

	// EnableTile route to echo
	routeGroup := m.store.MonitorableRouter.Group("/http", variant)
	routeStatus := routeGroup.GET("/status", delivery.GetHTTPStatus)
	routeRaw := routeGroup.GET("/raw", delivery.GetHTTPRaw)
	routeJSON := routeGroup.GET("/formatted", delivery.GetHTTPFormatted)

	// EnableTile data for config hydration
	m.store.UIConfigManager.EnableTile(api.HTTPStatusTileType, variant,
		&httpModels.HTTPStatusParams{}, routeStatus.Path, conf.InitialMaxDelay)
	m.store.UIConfigManager.EnableTile(api.HTTPRawTileType, variant,
		&httpModels.HTTPRawParams{}, routeRaw.Path, conf.InitialMaxDelay)
	m.store.UIConfigManager.EnableTile(api.HTTPFormattedTileType, variant,
		&httpModels.HTTPFormattedParams{}, routeJSON.Path, conf.InitialMaxDelay)
}
