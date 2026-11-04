//go:build faker

package uptimecom

import (
	"github.com/monitoror/monitoror/api/config/versions"
	"github.com/monitoror/monitoror/internal/pkg/monitorable"
	coreModels "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorables/uptimecom/api"
	uptimecomDelivery "github.com/monitoror/monitoror/monitorables/uptimecom/api/delivery/http"
	uptimecomModels "github.com/monitoror/monitoror/monitorables/uptimecom/api/models"
	uptimecomUsecase "github.com/monitoror/monitoror/monitorables/uptimecom/api/usecase"
	"github.com/monitoror/monitoror/registry"
	"github.com/monitoror/monitoror/store"
)

type Monitorable struct {
	monitorable.DefaultMonitorableFaker

	store *store.Store

	// Config tile settings
	checkTileEnabler            registry.TileEnabler
	transactionCheckTileEnabler registry.TileEnabler
}

func NewMonitorable(store *store.Store) *Monitorable {
	m := &Monitorable{}
	m.store = store

	// Register Monitorable Tile in config manager
	m.checkTileEnabler = store.Registry.RegisterTile(api.UptimecomCheckTileType, versions.Version2001, m.GetVariantsNames())

	return m
}

func (m *Monitorable) GetDisplayName() string { return "Uptimecom" }

func (m *Monitorable) Enable(variantName coreModels.VariantName) {
	usecase := uptimecomUsecase.NewUptimecomUsecase()
	delivery := uptimecomDelivery.NewUptimecomDelivery(usecase)

	// EnableTile route to echo
	routeGroup := m.store.MonitorableRouter.Group("/uptimecom", variantName)
	checkRoute := routeGroup.GET("/check", delivery.GetCheck)

	// EnableTile data for config hydration
	m.checkTileEnabler.Enable(variantName, &uptimecomModels.CheckParams{}, checkRoute.Path)
}
