//+build faker

package port

import (
	uiConfig "github.com/monitoror/monitoror/api/config"
	"github.com/monitoror/monitoror/api/config/versions"
	coreConfig "github.com/monitoror/monitoror/config"
	"github.com/monitoror/monitoror/internal/pkg/monitorable"
	coreModels "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorables/port/api"
	portDelivery "github.com/monitoror/monitoror/monitorables/port/api/delivery/http"
	portModels "github.com/monitoror/monitoror/monitorables/port/api/models"
	portUsecase "github.com/monitoror/monitoror/monitorables/port/api/usecase"
	"github.com/monitoror/monitoror/service/store"
)

type Monitorable struct {
	monitorable.DefaultMonitorableFaker

	store *store.Store

	// Config tile settings
	portTileSetting uiConfig.TileEnabler
}

func NewMonitorable(store *store.Store) *Monitorable {
	m := &Monitorable{}
	m.store = store

	// Register Monitorable Tile in config manager
	m.portTileSetting = store.TileSettingManager.Register(api.PortTileType, versions.MinimalVersion, m.GetVariants())

	return m
}

func (m *Monitorable) GetDisplayName() string { return "Port (faker)" }

func (m *Monitorable) Enable(variant coreModels.VariantName) {
	usecase := portUsecase.NewPortUsecase()
	delivery := portDelivery.NewPortDelivery(usecase)

	// EnableTile route to echo
	routeGroup := m.store.MonitorableRouter.Group("/port", variant)
	route := routeGroup.GET("/port", delivery.GetPort)

	// EnableTile data for config hydration
	m.portTileSetting.Enable(variant, &portModels.PortParams{}, route.Path, coreConfig.DefaultInitialMaxDelay)
}
