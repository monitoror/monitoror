//+build !faker

package travisci

import (
	"github.com/monitoror/monitoror/api/config/versions"
	pkgMonitorable "github.com/monitoror/monitoror/internal/pkg/monitorable"

	coreModels "github.com/monitoror/monitoror/models"

	"github.com/monitoror/monitoror/monitorables/travisci/api"
	travisciDelivery "github.com/monitoror/monitoror/monitorables/travisci/api/delivery/http"
	travisciModels "github.com/monitoror/monitoror/monitorables/travisci/api/models"
	travisciRepository "github.com/monitoror/monitoror/monitorables/travisci/api/repository"
	travisciUsecase "github.com/monitoror/monitoror/monitorables/travisci/api/usecase"
	travisciConfig "github.com/monitoror/monitoror/monitorables/travisci/config"
	"github.com/monitoror/monitoror/service/registry"
	"github.com/monitoror/monitoror/service/store"
)

type Monitorable struct {
	store *store.Store

	config map[coreModels.VariantName]*travisciConfig.TravisCI

	// Config tile settings
	buildTileEnabler registry.TileEnabler
}

func NewMonitorable(store *store.Store) *Monitorable {
	m := &Monitorable{}
	m.store = store
	m.config = make(map[coreModels.VariantName]*travisciConfig.TravisCI)

	// Load core config from env
	pkgMonitorable.LoadConfig(&m.config, travisciConfig.Default)

	// Register Monitorable Tile in config manager
	m.buildTileEnabler = store.Registry.RegisterTile(api.TravisCIBuildTileType, versions.MinimalVersion, m.GetVariantsNames())

	return m
}

func (m *Monitorable) GetDisplayName() string {
	return "Travis CI"
}

func (m *Monitorable) GetVariantsNames() []coreModels.VariantName {
	return pkgMonitorable.GetVariantsNames(m.config)
}

func (m *Monitorable) Validate(variantName coreModels.VariantName) (bool, []error) {
	conf := m.config[variantName]

	// Validate Config
	if errors := pkgMonitorable.ValidateConfig(conf, variantName); errors != nil {
		return false, errors
	}

	return true, nil
}

func (m *Monitorable) Enable(variantName coreModels.VariantName) {
	conf := m.config[variantName]

	repository := travisciRepository.NewTravisCIRepository(conf)
	usecase := travisciUsecase.NewTravisCIUsecase(repository)
	delivery := travisciDelivery.NewTravisCIDelivery(usecase)

	// EnableTile route to echo
	routeGroup := m.store.MonitorableRouter.Group("/travisci", variantName)
	route := routeGroup.GET("/build", delivery.GetBuild)

	// EnableTile data for config hydration
	m.buildTileEnabler.Enable(variantName, &travisciModels.BuildParams{}, route.Path)
}
