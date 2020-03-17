//+build faker

package monitorable

import (
	"github.com/monitoror/monitoror/config"
	monitorableConfig "github.com/monitoror/monitoror/monitorable/config"
	"github.com/monitoror/monitoror/monitorable/ping"
	pingDelivery "github.com/monitoror/monitoror/monitorable/ping/delivery/http"
	pingModels "github.com/monitoror/monitoror/monitorable/ping/models"
	pingUsecase "github.com/monitoror/monitoror/monitorable/ping/usecase"
	"github.com/monitoror/monitoror/service/router"
)

type pingMonitorable struct{}

func NewPingMonitorable(_ map[string]*config.Ping) Monitorable {
	return &pingMonitorable{}
}

func (m *pingMonitorable) GetHelp() string       { return "" }
func (m *pingMonitorable) GetVariants() []string { return []string{config.DefaultVariant} }

func (m *pingMonitorable) Register(variant string, router router.MonitorableRouter, configManager monitorableConfig.Manager) bool {
	usecase := pingUsecase.NewPingUsecase()
	delivery := pingDelivery.NewPingDelivery(usecase)

	// RegisterTile route to echo
	route := router.Group("/ping", variant).GET("/ping", delivery.GetPing)

	// RegisterTile data for config hydration
	configManager.RegisterTile(ping.PingTileType, variant, &pingModels.PingParams{}, route.Path, config.DefaultInitialMaxDelay)

	return true
}
