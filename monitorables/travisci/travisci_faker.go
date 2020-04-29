//+build faker

package travisci

import (
	"github.com/monitoror/monitoror/api/config/versions"
	"github.com/monitoror/monitoror/internal/pkg/monitorable"
	coreModels "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorables/travisci/api"
	travisciDelivery "github.com/monitoror/monitoror/monitorables/travisci/api/delivery/http"
	travisciModels "github.com/monitoror/monitoror/monitorables/travisci/api/models"
	travisciUsecase "github.com/monitoror/monitoror/monitorables/travisci/api/usecase"
	"github.com/monitoror/monitoror/registry"
	"github.com/monitoror/monitoror/store"
)

type Monitorable struct {
	monitorable.DefaultMonitorableFaker

	store *store.Store

	// Config tile settings
	buildTileEnabler registry.TileEnabler
}

func NewMonitorable(store *store.Store) *Monitorable {
	m := &Monitorable{}
	m.store = store

	// Register Monitorable Tile in config manager
	m.buildTileEnabler = store.Registry.RegisterTile(api.TravisCIBuildTileType, versions.MinimalVersion, m.GetVariantsNames())

	return m
}

func (m *Monitorable) GetDisplayName() string { return "Travis CI (faker)" }

func (m *Monitorable) Enable(variantName coreModels.VariantName) {
	usecase := travisciUsecase.NewTravisCIUsecase()
	delivery := travisciDelivery.NewTravisCIDelivery(usecase)

	// EnableTile route to echo
	routeGroup := m.store.MonitorableRouter.Group("/travisci", variantName)
	route := routeGroup.GET("/build", delivery.GetBuild)

	// EnableTile data for config hydration
	m.buildTileEnabler.Enable(variantName, &travisciModels.BuildParams{}, route.Path)
}
