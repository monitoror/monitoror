//+build !faker

package http

import (
	uiConfig "github.com/monitoror/monitoror/api/config"
	"github.com/monitoror/monitoror/api/config/versions"
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

	// Config tile settings
	statusTileSetting    uiConfig.TileEnabler
	rawTileSetting       uiConfig.TileEnabler
	formattedTileSetting uiConfig.TileEnabler
}

func NewMonitorable(store *store.Store) *Monitorable {
	m := &Monitorable{}
	m.store = store
	m.config = make(map[coreModels.VariantName]*httpConfig.HTTP)

	// Load core config from env
	pkgMonitorable.LoadConfig(&m.config, httpConfig.Default)

	// Register Monitorable Tile in config manager
	m.statusTileSetting = store.TileSettingManager.Register(api.HTTPStatusTileType, versions.MinimalVersion, m.GetVariantNames())
	m.rawTileSetting = store.TileSettingManager.Register(api.HTTPRawTileType, versions.MinimalVersion, m.GetVariantNames())
	m.formattedTileSetting = store.TileSettingManager.Register(api.HTTPFormattedTileType, versions.MinimalVersion, m.GetVariantNames())

	return m
}

func (m *Monitorable) GetDisplayName() string {
	return "HTTP"
}

func (m *Monitorable) GetVariantNames() []coreModels.VariantName {
	return pkgMonitorable.GetVariants(m.config)
}

func (m *Monitorable) Validate(_ coreModels.VariantName) (bool, error) {
	return true, nil
}

func (m *Monitorable) Enable(variantName coreModels.VariantName) {
	conf := m.config[variantName]

	repository := httpRepository.NewHTTPRepository(conf)
	usecase := httpUsecase.NewHTTPUsecase(repository, m.store.CacheStore, m.store.CoreConfig.UpstreamCacheExpiration)
	delivery := httpDelivery.NewHTTPDelivery(usecase)

	// EnableTile route to echo
	routeGroup := m.store.MonitorableRouter.Group("/http", variantName)
	routeStatus := routeGroup.GET("/status", delivery.GetHTTPStatus)
	routeRaw := routeGroup.GET("/raw", delivery.GetHTTPRaw)
	routeJSON := routeGroup.GET("/formatted", delivery.GetHTTPFormatted)

	// EnableTile data for config hydration
	m.statusTileSetting.Enable(variant, &httpModels.HTTPStatusParams{}, routeStatus.Path, conf.InitialMaxDelay)
	m.rawTileSetting.Enable(variant, &httpModels.HTTPRawParams{}, routeRaw.Path, conf.InitialMaxDelay)
	m.formattedTileSetting.Enable(variant, &httpModels.HTTPFormattedParams{}, routeJSON.Path, conf.InitialMaxDelay)
}
