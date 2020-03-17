//+build !faker

package monitorable

import (
	"github.com/monitoror/monitoror/config"
	monitorableConfig "github.com/monitoror/monitoror/monitorable/config"
	"github.com/monitoror/monitoror/monitorable/ping"
	pingDelivery "github.com/monitoror/monitoror/monitorable/ping/delivery/http"
	pingModels "github.com/monitoror/monitoror/monitorable/ping/models"
	pingRepository "github.com/monitoror/monitoror/monitorable/ping/repository"
	pingUsecase "github.com/monitoror/monitoror/monitorable/ping/usecase"
	"github.com/monitoror/monitoror/pkg/monitoror/utils/system"
	"github.com/monitoror/monitoror/service/router"
)

type pingMonitorable struct {
	config map[string]*config.Ping
}

func NewPingMonitorable(config map[string]*config.Ping) Monitorable {
	return &pingMonitorable{config: config}
}

func (m *pingMonitorable) GetHelp() string         { return "HEEEELLLPPPP" }
func (m *pingMonitorable) GetVariants() []string   { return config.GetVariantsFromConfig(m.config) }
func (m *pingMonitorable) isEnabled(_ string) bool { return system.IsRawSocketAvailable() }

func (m *pingMonitorable) Register(variant string, router router.MonitorableRouter, configManager monitorableConfig.Manager) bool {
	enabled := m.isEnabled(variant)
	if enabled {
		conf := m.config[variant]

		repository := pingRepository.NewPingRepository(conf)
		usecase := pingUsecase.NewPingUsecase(repository)
		delivery := pingDelivery.NewPingDelivery(usecase)

		// RegisterTile route to echo
		route := router.Group("/ping", variant).GET("/ping", delivery.GetPing)

		// RegisterTile data for config hydration
		configManager.RegisterTile(ping.PingTileType, variant, &pingModels.PingParams{}, route.Path, conf.InitialMaxDelay)
	} else {
		// RegisterTile data for config verify
		configManager.DisableTile(ping.PingTileType, variant)
	}

	return enabled
}
