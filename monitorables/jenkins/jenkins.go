//+build !faker

package jenkins

import (
	"net/url"

	uiConfig "github.com/monitoror/monitoror/api/config/usecase"
	coreConfig "github.com/monitoror/monitoror/config"
	"github.com/monitoror/monitoror/monitorables/jenkins/api"
	jenkinsDelivery "github.com/monitoror/monitoror/monitorables/jenkins/api/delivery/http"
	jenkinsModels "github.com/monitoror/monitoror/monitorables/jenkins/api/models"
	jenkinsRepository "github.com/monitoror/monitoror/monitorables/jenkins/api/repository"
	jenkinsUsecase "github.com/monitoror/monitoror/monitorables/jenkins/api/usecase"
	jenkinsCoreConfig "github.com/monitoror/monitoror/monitorables/jenkins/config"
	"github.com/monitoror/monitoror/service/store"
)

type Monitorable struct {
	store *store.Store

	config map[string]*jenkinsCoreConfig.Jenkins
}

func NewMonitorable(store *store.Store) *Monitorable {
	monitorable := &Monitorable{}
	monitorable.store = store
	monitorable.config = make(map[string]*jenkinsCoreConfig.Jenkins)

	// Load core config from env
	coreConfig.LoadMonitorableConfig(&monitorable.config, jenkinsCoreConfig.Default)

	// Register Monitorable Tile in config manager
	store.UIConfigManager.RegisterTile(api.JenkinsBuildTileType, monitorable.GetVariants(), uiConfig.MinimalVersion)
	store.UIConfigManager.RegisterTile(api.JenkinsMultiBranchTileType, monitorable.GetVariants(), uiConfig.MinimalVersion)

	return monitorable
}

func (m *Monitorable) GetVariants() []string { return coreConfig.GetVariantsFromConfig(m.config) }
func (m *Monitorable) IsValid(variant string) bool {
	conf := m.config[variant]

	if conf.URL == "" {
		return false
	}

	if _, err := url.Parse(conf.URL); err != nil {
		return false
	}

	return true
}

func (m *Monitorable) Enable(variant string) {
	conf := m.config[variant]

	repository := jenkinsRepository.NewJenkinsRepository(conf)
	usecase := jenkinsUsecase.NewJenkinsUsecase(repository)
	delivery := jenkinsDelivery.NewJenkinsDelivery(usecase)

	// EnableTile route to echo
	jenkinsGroup := m.store.MonitorableRouter.Group("/jenkins", variant)
	route := jenkinsGroup.GET("/build", delivery.GetBuild)

	// EnableTile data for config hydration
	m.store.UIConfigManager.EnableTile(api.JenkinsBuildTileType, variant, &jenkinsModels.BuildParams{}, route.Path, conf.InitialMaxDelay)
	m.store.UIConfigManager.EnableDynamicTile(api.JenkinsMultiBranchTileType, variant, &jenkinsModels.MultiBranchParams{}, usecase.MultiBranch)
}
