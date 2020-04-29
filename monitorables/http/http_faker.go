//+build faker

package http

import (
	"github.com/monitoror/monitoror/api/config/versions"
	"github.com/monitoror/monitoror/internal/pkg/monitorable"
	coreModels "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorables/http/api"
	httpDelivery "github.com/monitoror/monitoror/monitorables/http/api/delivery/http"
	httpModels "github.com/monitoror/monitoror/monitorables/http/api/models"
	httpUsecase "github.com/monitoror/monitoror/monitorables/http/api/usecase"
	"github.com/monitoror/monitoror/registry"
	"github.com/monitoror/monitoror/store"
)

type Monitorable struct {
	monitorable.DefaultMonitorableFaker

	store *store.Store

	// Config tile settings
	statusTileEnabler    registry.TileEnabler
	rawTileEnabler       registry.TileEnabler
	formattedTileEnabler registry.TileEnabler
}

func NewMonitorable(store *store.Store) *Monitorable {
	m := &Monitorable{}
	m.store = store

	// Register Monitorable Tile in config manager
	m.statusTileEnabler = store.Registry.RegisterTile(api.HTTPStatusTileType, versions.MinimalVersion, m.GetVariantsNames())
	m.rawTileEnabler = store.Registry.RegisterTile(api.HTTPRawTileType, versions.MinimalVersion, m.GetVariantsNames())
	m.formattedTileEnabler = store.Registry.RegisterTile(api.HTTPFormattedTileType, versions.MinimalVersion, m.GetVariantsNames())

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
	m.statusTileEnabler.Enable(variantName, &httpModels.HTTPStatusParams{}, routeStatus.Path)
	m.rawTileEnabler.Enable(variantName, &httpModels.HTTPRawParams{}, routeRaw.Path)
	m.formattedTileEnabler.Enable(variantName, &httpModels.HTTPFormattedParams{}, routeJSON.Path)
}
