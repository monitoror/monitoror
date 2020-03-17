//+build faker

package monitorable

import (
	"github.com/monitoror/monitoror/config"
	monitorableConfig "github.com/monitoror/monitoror/monitorable/config"
	"github.com/monitoror/monitoror/monitorable/travisci"
	travisciDelivery "github.com/monitoror/monitoror/monitorable/travisci/delivery/http"
	travisciModels "github.com/monitoror/monitoror/monitorable/travisci/models"
	travisciUsecase "github.com/monitoror/monitoror/monitorable/travisci/usecase"
	"github.com/monitoror/monitoror/service/router"
)

type travisciMonitorable struct{}

func NewTravisCIMonitorable(_ map[string]*config.TravisCI) Monitorable {
	return &travisciMonitorable{}
}

func (m *travisciMonitorable) GetHelp() string       { return "" }
func (m *travisciMonitorable) GetVariants() []string { return []string{config.DefaultVariant} }

func (m *travisciMonitorable) Register(variant string, router router.MonitorableRouter, configManager monitorableConfig.Manager) bool {
	usecase := travisciUsecase.NewTravisCIUsecase()
	delivery := travisciDelivery.NewTravisCIDelivery(usecase)

	// RegisterTile route to echo
	route := router.Group("/travisci", variant).GET("/build", delivery.GetBuild)

	// RegisterTile data for config hydration
	configManager.RegisterTile(travisci.TravisCIBuildTileType, variant, &travisciModels.BuildParams{}, route.Path, config.DefaultInitialMaxDelay)

	return true
}
