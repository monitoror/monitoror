//+build faker

package travisci

import (
	uiConfig "github.com/monitoror/monitoror/api/config/usecase"
	coreConfig "github.com/monitoror/monitoror/config"
	"github.com/monitoror/monitoror/monitorables/travisci/api"
	travisciDelivery "github.com/monitoror/monitoror/monitorables/travisci/api/delivery/http"
	travisciModels "github.com/monitoror/monitoror/monitorables/travisci/api/models"
	travisciUsecase "github.com/monitoror/monitoror/monitorables/travisci/api/usecase"
	"github.com/monitoror/monitoror/service/store"
)

type Monitorable struct {
	store *store.Store
}

func NewMonitorable(store *store.Store) *Monitorable {
	monitorable := &Monitorable{}
	monitorable.store = store

	// Register Monitorable Tile in config manager
	store.UIConfigManager.RegisterTile(api.TravisCIBuildTileType, monitorable.GetVariants(), uiConfig.MinimalVersion)

	return monitorable
}

func (m *Monitorable) GetVariants() []string       { return []string{coreConfig.DefaultVariant} }
func (m *Monitorable) IsValid(variant string) bool { return true }

func (m *Monitorable) Enable(variant string) {
	usecase := travisciUsecase.NewTravisCIUsecase()
	delivery := travisciDelivery.NewTravisCIDelivery(usecase)

	// EnableTile route to echo
	route := m.store.MonitorableRouter.Group("/travisci", variant).GET("/build", delivery.GetBuild)

	// EnableTile data for config hydration
	m.store.UIConfigManager.EnableTile(api.TravisCIBuildTileType, variant, &travisciModels.BuildParams{}, route.Path, coreConfig.DefaultInitialMaxDelay)
}
