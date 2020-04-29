//+build faker

package port

import (
	"github.com/monitoror/monitoror/api/config/versions"
	"github.com/monitoror/monitoror/internal/pkg/monitorable"
	coreModels "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorables/port/api"
	portDelivery "github.com/monitoror/monitoror/monitorables/port/api/delivery/http"
	portModels "github.com/monitoror/monitoror/monitorables/port/api/models"
	portUsecase "github.com/monitoror/monitoror/monitorables/port/api/usecase"
	"github.com/monitoror/monitoror/registry"
	"github.com/monitoror/monitoror/store"
)

type Monitorable struct {
	monitorable.DefaultMonitorableFaker

	store *store.Store

	// Config tile settings
	portTileEnabler registry.TileEnabler
}

func NewMonitorable(store *store.Store) *Monitorable {
	m := &Monitorable{}
	m.store = store

	// Register Monitorable Tile in config manager
	m.portTileEnabler = store.Registry.RegisterTile(api.PortTileType, versions.MinimalVersion, m.GetVariantsNames())

	return m
}

func (m *Monitorable) GetDisplayName() string { return "Port (faker)" }

func (m *Monitorable) Enable(variantName coreModels.VariantName) {
	usecase := portUsecase.NewPortUsecase()
	delivery := portDelivery.NewPortDelivery(usecase)

	// EnableTile route to echo
	routeGroup := m.store.MonitorableRouter.Group("/port", variantName)
	route := routeGroup.GET("/port", delivery.GetPort)

	// EnableTile data for config hydration
	m.portTileEnabler.Enable(variantName, &portModels.PortParams{}, route.Path)
}
