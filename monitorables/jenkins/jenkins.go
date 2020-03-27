//+build !faker

package jenkins

import (
	"fmt"
	"net/url"

	pkgMonitorable "github.com/monitoror/monitoror/internal/pkg/monitorable"

	coreModels "github.com/monitoror/monitoror/models"

	uiConfig "github.com/monitoror/monitoror/api/config/usecase"
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
}

func NewMonitorable(store *store.Store) *Monitorable {
	monitorable := &Monitorable{}
	monitorable.store = store
	monitorable.config = make(map[coreModels.VariantName]*jenkinsConfig.Jenkins)

	// Load core config from env
	pkgMonitorable.LoadConfig(&monitorable.config, jenkinsConfig.Default)

	// Register Monitorable Tile in config manager
	store.UIConfigManager.RegisterTile(api.JenkinsBuildTileType, monitorable.GetVariants(), uiConfig.MinimalVersion)
	store.UIConfigManager.RegisterTile(api.JenkinsMultiBranchTileType, monitorable.GetVariants(), uiConfig.MinimalVersion)

	return monitorable
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
	m.store.UIConfigManager.EnableTile(api.JenkinsBuildTileType, variant,
		&jenkinsModels.BuildParams{}, route.Path, conf.InitialMaxDelay)
	m.store.UIConfigManager.EnableDynamicTile(api.JenkinsMultiBranchTileType, variant,
		&jenkinsModels.MultiBranchParams{}, usecase.MultiBranch)
}
