//+build faker

package travisci

import (
	uiConfig "github.com/monitoror/monitoror/api/config"
	"github.com/monitoror/monitoror/api/config/versions"
	coreConfig "github.com/monitoror/monitoror/config"
	"github.com/monitoror/monitoror/internal/pkg/monitorable"
	coreModels "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorables/travisci/api"
	travisciDelivery "github.com/monitoror/monitoror/monitorables/travisci/api/delivery/http"
	travisciModels "github.com/monitoror/monitoror/monitorables/travisci/api/models"
	travisciUsecase "github.com/monitoror/monitoror/monitorables/travisci/api/usecase"
	"github.com/monitoror/monitoror/service/store"
)

type Monitorable struct {
	monitorable.DefaultMonitorableFaker

	store *store.Store

	// Config tile settings
	buildTileSetting uiConfig.TileEnabler
}

func NewMonitorable(store *store.Store) *Monitorable {
	m := &Monitorable{}
	m.store = store

	// Register Monitorable Tile in config manager
	m.buildTileSetting = store.TileSettingManager.Register(api.TravisCIBuildTileType, versions.MinimalVersion, m.GetVariants())

	return m
}

func (m *Monitorable) GetDisplayName() string { return "Travis CI (faker)" }

func (m *Monitorable) Enable(variant coreModels.VariantName) {
	usecase := travisciUsecase.NewTravisCIUsecase()
	delivery := travisciDelivery.NewTravisCIDelivery(usecase)

	// EnableTile route to echo
	routeGroup := m.store.MonitorableRouter.Group("/travisci", variant)
	route := routeGroup.GET("/build", delivery.GetBuild)

	// EnableTile data for config hydration
	m.buildTileSetting.Enable(variant, &travisciModels.BuildParams{}, route.Path, coreConfig.DefaultInitialMaxDelay)
}
