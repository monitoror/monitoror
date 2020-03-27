//+build faker

package http

import (
	uiConfig "github.com/monitoror/monitoror/api/config/usecase"
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
}

func NewMonitorable(store *store.Store) *Monitorable {
	monitorable := &Monitorable{}
	monitorable.store = store

	// Register Monitorable Tile in config manager
	store.UIConfigManager.RegisterTile(api.HTTPStatusTileType, monitorable.GetVariants(), uiConfig.MinimalVersion)
	store.UIConfigManager.RegisterTile(api.HTTPRawTileType, monitorable.GetVariants(), uiConfig.MinimalVersion)
	store.UIConfigManager.RegisterTile(api.HTTPFormattedTileType, monitorable.GetVariants(), uiConfig.MinimalVersion)

	return monitorable
}

func (m *Monitorable) GetDisplayName() string { return "HTTP (faker)" }

func (m *Monitorable) Enable(variant coreModels.VariantName) {
	usecase := httpUsecase.NewHTTPUsecase()
	delivery := httpDelivery.NewHTTPDelivery(usecase)

	// EnableTile route to echo
	routeGroup := m.store.MonitorableRouter.Group("/http", variant)
	routeStatus := routeGroup.GET("/status", delivery.GetHTTPStatus)
	routeRaw := routeGroup.GET("/raw", delivery.GetHTTPRaw)
	routeJSON := routeGroup.GET("/formatted", delivery.GetHTTPFormatted)

	// EnableTile data for config hydration
	m.store.UIConfigManager.EnableTile(api.HTTPStatusTileType, variant,
		&httpModels.HTTPStatusParams{}, routeStatus.Path, coreConfig.DefaultInitialMaxDelay)
	m.store.UIConfigManager.EnableTile(api.HTTPRawTileType, variant,
		&httpModels.HTTPRawParams{}, routeRaw.Path, coreConfig.DefaultInitialMaxDelay)
	m.store.UIConfigManager.EnableTile(api.HTTPFormattedTileType, variant,
		&httpModels.HTTPFormattedParams{}, routeJSON.Path, coreConfig.DefaultInitialMaxDelay)
}
