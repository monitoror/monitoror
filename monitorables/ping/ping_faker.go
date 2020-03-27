//+build faker

package ping

import (
	uiConfig "github.com/monitoror/monitoror/api/config/usecase"
	coreConfig "github.com/monitoror/monitoror/config"
	"github.com/monitoror/monitoror/internal/pkg/monitorable"
	coreModels "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorables/ping/api"
	pingDelivery "github.com/monitoror/monitoror/monitorables/ping/api/delivery/http"
	pingModels "github.com/monitoror/monitoror/monitorables/ping/api/models"
	pingUsecase "github.com/monitoror/monitoror/monitorables/ping/api/usecase"
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
	store.UIConfigManager.RegisterTile(api.PingTileType, monitorable.GetVariants(), uiConfig.MinimalVersion)

	return monitorable
}

func (m *Monitorable) GetDisplayName() string { return "Ping (faker)" }

func (m *Monitorable) Enable(variant coreModels.VariantName) {

	usecase := pingUsecase.NewPingUsecase()
	delivery := pingDelivery.NewPingDelivery(usecase)

	// EnableTile route to echo
	routeGroup := m.store.MonitorableRouter.Group("/ping", variant)
	route := routeGroup.GET("/ping", delivery.GetPing)

	// EnableTile data for config hydration
	m.store.UIConfigManager.EnableTile(api.PingTileType, variant,
		&pingModels.PingParams{}, route.Path, coreConfig.DefaultInitialMaxDelay)
}
