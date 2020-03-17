//+build !faker

package monitorable

import (
	"net/url"

	"github.com/monitoror/monitoror/config"
	monitorableConfig "github.com/monitoror/monitoror/monitorable/config"
	"github.com/monitoror/monitoror/monitorable/jenkins"
	jenkinsDelivery "github.com/monitoror/monitoror/monitorable/jenkins/delivery/http"
	jenkinsModels "github.com/monitoror/monitoror/monitorable/jenkins/models"
	jenkinsRepository "github.com/monitoror/monitoror/monitorable/jenkins/repository"
	jenkinsUsecase "github.com/monitoror/monitoror/monitorable/jenkins/usecase"
	"github.com/monitoror/monitoror/service/router"
)

type jenkinsMonitorable struct {
	config map[string]*config.Jenkins
}

func NewJenkinsMonitorable(config map[string]*config.Jenkins) Monitorable {
	return &jenkinsMonitorable{config: config}
}

func (m *jenkinsMonitorable) GetHelp() string       { return "HEEEELLLPPPP" }
func (m *jenkinsMonitorable) GetVariants() []string { return config.GetVariantsFromConfig(m.config) }
func (m *jenkinsMonitorable) isEnabled(variant string) bool {
	conf := m.config[variant]

	if conf.URL == "" {
		return false
	}

	if _, err := url.Parse(conf.URL); err != nil {
		return false
	}

	return true
}

func (m *jenkinsMonitorable) Register(variant string, router router.MonitorableRouter, configManager monitorableConfig.Manager) bool {
	enabled := m.isEnabled(variant)
	if enabled {
		conf := m.config[variant]

		repository := jenkinsRepository.NewJenkinsRepository(conf)
		usecase := jenkinsUsecase.NewJenkinsUsecase(repository)
		delivery := jenkinsDelivery.NewJenkinsDelivery(usecase)

		// RegisterTile route to echo
		route := router.Group("/http", variant).GET("/build", delivery.GetBuild)

		// RegisterTile data for config hydration
		configManager.RegisterTile(jenkins.JenkinsBuildTileType, variant, &jenkinsModels.BuildParams{}, route.Path, conf.InitialMaxDelay)
		configManager.RegisterDynamicTile(jenkins.JenkinsMultiBranchTileType, variant, &jenkinsModels.MultiBranchParams{}, usecase.MultiBranch)
	} else {
		// RegisterTile data for config verify
		configManager.DisableTile(jenkins.JenkinsBuildTileType, variant)
		configManager.DisableTile(jenkins.JenkinsMultiBranchTileType, variant)
	}

	return enabled
}
