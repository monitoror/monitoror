//+build !faker

package monitorable

import (
	"github.com/monitoror/monitoror/config"
	monitorableConfig "github.com/monitoror/monitoror/monitorable/config"
	"github.com/monitoror/monitoror/monitorable/port"
	portDelivery "github.com/monitoror/monitoror/monitorable/port/delivery/http"
	portModels "github.com/monitoror/monitoror/monitorable/port/models"
	portRepository "github.com/monitoror/monitoror/monitorable/port/repository"
	portUsecase "github.com/monitoror/monitoror/monitorable/port/usecase"
	"github.com/monitoror/monitoror/service/router"
)

type portMonitorable struct {
	config map[string]*config.Port
}

func NewPortMonitorable(config map[string]*config.Port) Monitorable {
	return &portMonitorable{config: config}
}

func (m *portMonitorable) GetHelp() string         { return "HEEEELLLPPPP" }
func (m *portMonitorable) GetVariants() []string   { return config.GetVariantsFromConfig(m.config) }
func (m *portMonitorable) isEnabled(_ string) bool { return true }

func (m *portMonitorable) Register(variant string, router router.MonitorableRouter, configManager monitorableConfig.Manager) bool {
	enabled := m.isEnabled(variant)
	if enabled {
		conf := m.config[variant]

		repository := portRepository.NewPortRepository(conf)
		usecase := portUsecase.NewPortUsecase(repository)
		delivery := portDelivery.NewPortDelivery(usecase)

		// RegisterTile route to echo
		route := router.Group("/port", variant).GET("/port", delivery.GetPort)

		// RegisterTile data for config hydration
		configManager.RegisterTile(port.PortTileType, variant, &portModels.PortParams{}, route.Path, conf.InitialMaxDelay)
	} else {
		// RegisterTile data for config verify
		configManager.DisableTile(port.PortTileType, variant)
	}

	return enabled
}
