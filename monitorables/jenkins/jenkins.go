//+build !faker

package jenkins

import (
	"fmt"
	"net/url"

	uiConfig "github.com/monitoror/monitoror/api/config"
	"github.com/monitoror/monitoror/api/config/versions"
	pkgMonitorable "github.com/monitoror/monitoror/internal/pkg/monitorable"
	coreModels "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorables/jenkins/api"
	jenkinsDelivery "github.com/monitoror/monitoror/monitorables/jenkins/api/delivery/http"
	jenkinsModels "github.com/monitoror/monitoror/monitorables/jenkins/api/models"
	jenkinsRepository "github.com/monitoror/monitoror/monitorables/jenkins/api/repository"
	jenkinsUsecase "github.com/monitoror/monitoror/monitorables/jenkins/api/usecase"
	jenkinsConfig "github.com/monitoror/monitoror/monitorables/jenkins/config"
	"github.com/monitoror/monitoror/service/store"
)

type Monitorable struct {
	store *store.Store

	config map[coreModels.VariantName]*jenkinsConfig.Jenkins

	// Config tile settings
	buildTileSetting          uiConfig.TileEnabler
	buildGeneratorTileSetting uiConfig.TileGeneratorEnabler
}

func NewMonitorable(store *store.Store) *Monitorable {
	m := &Monitorable{}
	m.store = store
	m.config = make(map[coreModels.VariantName]*jenkinsConfig.Jenkins)

	// Load core config from env
	pkgMonitorable.LoadConfig(&m.config, jenkinsConfig.Default)

	// Register Monitorable Tile in config manager
	m.buildTileSetting = store.TileSettingManager.Register(api.JenkinsBuildTileType, versions.MinimalVersion, m.GetVariants())
	m.buildGeneratorTileSetting = store.TileSettingManager.RegisterGenerator(api.JenkinsBuildTileType, versions.MinimalVersion, m.GetVariants())

	return m
}

func (m *Monitorable) GetDisplayName() string {
	return "Jenkins"
}

func (m *Monitorable) GetVariants() []coreModels.VariantName {
	return pkgMonitorable.GetVariants(m.config)
}

func (m *Monitorable) Validate(variant coreModels.VariantName) (bool, error) {
	conf := m.config[variant]

	// No configuration set
	if conf.URL == "" {
		return false, nil
	}

	// Error in URL
	if _, err := url.Parse(conf.URL); err != nil {
		return false, fmt.Errorf(`%s contains invalid URL: "%s"`, pkgMonitorable.BuildMonitorableEnvKey(conf, variant, "URL"), conf.URL)
	}

	return true, nil
}

func (m *Monitorable) Enable(variant coreModels.VariantName) {
	conf := m.config[variant]

	repository := jenkinsRepository.NewJenkinsRepository(conf)
	usecase := jenkinsUsecase.NewJenkinsUsecase(repository)
	delivery := jenkinsDelivery.NewJenkinsDelivery(usecase)

	// EnableTile route to echo
	routeGroup := m.store.MonitorableRouter.Group("/jenkins", variant)
	route := routeGroup.GET("/build", delivery.GetBuild)

	// EnableTile data for config hydration
	m.buildTileSetting.Enable(variant, &jenkinsModels.BuildParams{}, route.Path, conf.InitialMaxDelay)
	m.buildGeneratorTileSetting.Enable(variant, &jenkinsModels.BuildGeneratorParams{}, usecase.BuildGenerator)
}
