//+build !faker

package monitorable

import (
	"net/url"

	"github.com/monitoror/monitoror/config"
	monitorableConfig "github.com/monitoror/monitoror/monitorable/config"
	"github.com/monitoror/monitoror/monitorable/travisci"
	travisciDelivery "github.com/monitoror/monitoror/monitorable/travisci/delivery/http"
	travisciModels "github.com/monitoror/monitoror/monitorable/travisci/models"
	travisciRepository "github.com/monitoror/monitoror/monitorable/travisci/repository"
	travisciUsecase "github.com/monitoror/monitoror/monitorable/travisci/usecase"
	"github.com/monitoror/monitoror/service/router"
)

type travisciMonitorable struct {
	config map[string]*config.TravisCI
}

func NewTravisCIMonitorable(config map[string]*config.TravisCI) Monitorable {
	return &travisciMonitorable{config: config}
}

func (m *travisciMonitorable) GetHelp() string       { return "HEEEELLLPPPP" }
func (m *travisciMonitorable) GetVariants() []string { return config.GetVariantsFromConfig(m.config) }
func (m *travisciMonitorable) isEnabled(variant string) bool {
	conf := m.config[variant]

	if conf.URL == "" {
		return false
	}

	if _, err := url.Parse(conf.URL); err != nil {
		return false
	}

	return true
}

func (m *travisciMonitorable) Register(variant string, router router.MonitorableRouter, configManager monitorableConfig.Manager) bool {
	enabled := m.isEnabled(variant)
	if enabled {
		conf := m.config[variant]

		repository := travisciRepository.NewTravisCIRepository(conf)
		usecase := travisciUsecase.NewTravisCIUsecase(repository)
		delivery := travisciDelivery.NewTravisCIDelivery(usecase)

		// RegisterTile route to echo
		route := router.Group("/travisci", variant).GET("/build", delivery.GetBuild)

		// RegisterTile data for config hydration
		configManager.RegisterTile(travisci.TravisCIBuildTileType, variant, &travisciModels.BuildParams{}, route.Path, conf.InitialMaxDelay)
	} else {
		// RegisterTile data for config verify
		configManager.DisableTile(travisci.TravisCIBuildTileType, variant)
	}
	return enabled
}
