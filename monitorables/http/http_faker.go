//+build faker

package http

import (
	uiConfig "github.com/monitoror/monitoror/api/config"
	"github.com/monitoror/monitoror/api/config/versions"
	coreConfig "github.com/monitoror/monitoror/config"
	"github.com/monitoror/monitoror/internal/pkg/monitorable"
	coreModels "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorables/http/api"
	httpDelivery "github.com/monitoror/monitoror/monitorables/http/api/delivery/http"
	httpModels "github.com/monitoror/monitoror/monitorables/http/api/models"
	httpUsecase "github.com/monitoror/monitoror/monitorables/http/api/usecase"
	"github.com/monitoror/monitoror/service/store"
)

type Monitorable struct {
	monitorable.DefaultMonitorableFaker

	store *store.Store

	// Config tile settings
	statusTileSetting    uiConfig.TileEnabler
	rawTileSetting       uiConfig.TileEnabler
	formattedTileSetting uiConfig.TileEnabler
}

func NewMonitorable(store *store.Store) *Monitorable {
	m := &Monitorable{}
	m.store = store

	// Register Monitorable Tile in config manager
	m.statusTileSetting = store.TileSettingManager.Register(api.HTTPStatusTileType, versions.MinimalVersion, m.GetVariantNames())
	m.rawTileSetting = store.TileSettingManager.Register(api.HTTPRawTileType, versions.MinimalVersion, m.GetVariantNames())
	m.formattedTileSetting = store.TileSettingManager.Register(api.HTTPFormattedTileType, versions.MinimalVersion, m.GetVariantNames())

	return m
}

func (m *Monitorable) GetDisplayName() string { return "HTTP (faker)" }

func (m *Monitorable) Enable(variantName coreModels.VariantName) {
	usecase := httpUsecase.NewHTTPUsecase()
	delivery := httpDelivery.NewHTTPDelivery(usecase)

	// EnableTile route to echo
	routeGroup := m.store.MonitorableRouter.Group("/http", variantName)
	routeStatus := routeGroup.GET("/status", delivery.GetHTTPStatus)
	routeRaw := routeGroup.GET("/raw", delivery.GetHTTPRaw)
	routeJSON := routeGroup.GET("/formatted", delivery.GetHTTPFormatted)

	// EnableTile data for config hydration
	m.statusTileSetting.Enable(variant, &httpModels.HTTPStatusParams{}, routeStatus.Path, coreConfig.DefaultInitialMaxDelay)
	m.rawTileSetting.Enable(variant, &httpModels.HTTPRawParams{}, routeRaw.Path, coreConfig.DefaultInitialMaxDelay)
	m.formattedTileSetting.Enable(variant, &httpModels.HTTPFormattedParams{}, routeJSON.Path, coreConfig.DefaultInitialMaxDelay)
}
