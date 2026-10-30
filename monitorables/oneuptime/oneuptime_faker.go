//go:build faker

package oneuptime

import (
	"github.com/monitoror/monitoror/api/config/versions"
	"github.com/monitoror/monitoror/internal/pkg/monitorable"
	coreModels "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorables/oneuptime/api"
	oneuptimeDelivery "github.com/monitoror/monitoror/monitorables/oneuptime/api/delivery/http"
	oneuptimeModels "github.com/monitoror/monitoror/monitorables/oneuptime/api/models"
	oneuptimeUsecase "github.com/monitoror/monitoror/monitorables/oneuptime/api/usecase"
	"github.com/monitoror/monitoror/registry"
	"github.com/monitoror/monitoror/store"
)

type Monitorable struct {
	monitorable.DefaultMonitorableFaker

	store *store.Store

	// Config tile settings
	monitorTileEnabler registry.TileEnabler
}

func NewMonitorable(store *store.Store) *Monitorable {
	m := &Monitorable{}
	m.store = store

	// Register Monitorable Tile in config manager
	m.monitorTileEnabler = store.Registry.RegisterTile(api.OneuptimeMonitorTileType, versions.Version2003, m.GetVariantsNames())

	return m
}

func (m *Monitorable) GetDisplayName() string {
	return "OneUptime"
}

func (m *Monitorable) Enable(variantName coreModels.VariantName) {
	usecase := oneuptimeUsecase.NewOneuptimeUsecase()
	delivery := oneuptimeDelivery.NewOneuptimeDelivery(usecase)

	// EnableTile route to echo
	routeGroup := m.store.MonitorableRouter.Group("/oneuptime", variantName)
	monitorRoute := routeGroup.GET("/monitor", delivery.GetMonitor)

	// EnableTile data for config hydration
	m.monitorTileEnabler.Enable(variantName, &oneuptimeModels.MonitorParams{}, monitorRoute.Path)
}
