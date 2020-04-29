//+build !faker

package jenkins

import (
	"github.com/monitoror/monitoror/api/config/versions"
	pkgMonitorable "github.com/monitoror/monitoror/internal/pkg/monitorable"
	coreModels "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorables/jenkins/api"
	jenkinsDelivery "github.com/monitoror/monitoror/monitorables/jenkins/api/delivery/http"
	jenkinsModels "github.com/monitoror/monitoror/monitorables/jenkins/api/models"
	jenkinsRepository "github.com/monitoror/monitoror/monitorables/jenkins/api/repository"
	jenkinsUsecase "github.com/monitoror/monitoror/monitorables/jenkins/api/usecase"
	jenkinsConfig "github.com/monitoror/monitoror/monitorables/jenkins/config"
	"github.com/monitoror/monitoror/registry"
	"github.com/monitoror/monitoror/store"
)

type Monitorable struct {
	store *store.Store

	config map[coreModels.VariantName]*jenkinsConfig.Jenkins

	// Config tile settings
	buildTileEnabler      registry.TileEnabler
	buildGeneratorEnabler registry.GeneratorEnabler
}

func NewMonitorable(store *store.Store) *Monitorable {
	m := &Monitorable{}
	m.store = store
	m.config = make(map[coreModels.VariantName]*jenkinsConfig.Jenkins)

	// Load core config from env
	pkgMonitorable.LoadConfig(&m.config, jenkinsConfig.Default)

	// Register Monitorable Tile in config manager
	m.buildTileEnabler = store.Registry.RegisterTile(api.JenkinsBuildTileType, versions.MinimalVersion, m.GetVariantsNames())
	m.buildGeneratorEnabler = store.Registry.RegisterGenerator(api.JenkinsBuildTileType, versions.MinimalVersion, m.GetVariantsNames())

	return m
}

func (m *Monitorable) GetDisplayName() string {
	return "Jenkins"
}

func (m *Monitorable) GetVariantsNames() []coreModels.VariantName {
	return pkgMonitorable.GetVariantsNames(m.config)
}

func (m *Monitorable) Validate(variantName coreModels.VariantName) (bool, []error) {
	conf := m.config[variantName]

	// No configuration set
	if conf.URL == "" {
		return false, nil
	}

	// Validate Config
	if err := pkgMonitorable.ValidateConfig(conf, variantName); err != nil {
		return false, err
	}

	return true, nil
}

func (m *Monitorable) Enable(variantName coreModels.VariantName) {
	conf := m.config[variantName]

	repository := jenkinsRepository.NewJenkinsRepository(conf)
	usecase := jenkinsUsecase.NewJenkinsUsecase(repository)
	delivery := jenkinsDelivery.NewJenkinsDelivery(usecase)

	// EnableTile route to echo
	routeGroup := m.store.MonitorableRouter.Group("/jenkins", variantName)
	route := routeGroup.GET("/build", delivery.GetBuild)

	// EnableTile data for config hydration
	m.buildTileEnabler.Enable(variantName, &jenkinsModels.BuildParams{}, route.Path)
	m.buildGeneratorEnabler.Enable(variantName, &jenkinsModels.BuildGeneratorParams{}, usecase.BuildGenerator)
}
