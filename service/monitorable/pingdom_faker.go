//+build faker

package monitorable

import (
	"github.com/jsdidierlaurent/echo-middleware/cache"

	"github.com/monitoror/monitoror/config"
	monitorableConfig "github.com/monitoror/monitoror/monitorable/config"
	"github.com/monitoror/monitoror/monitorable/pingdom"
	pingdomDelivery "github.com/monitoror/monitoror/monitorable/pingdom/delivery/http"
	pingdomModels "github.com/monitoror/monitoror/monitorable/pingdom/models"
	pingdomUsecase "github.com/monitoror/monitoror/monitorable/pingdom/usecase"
	"github.com/monitoror/monitoror/service/router"
)

type pingdomMonitorable struct{}

func NewPingdomMonitorable(_ map[string]*config.Pingdom, _ cache.Store) Monitorable {
	return &pingdomMonitorable{}
}

func (m *pingdomMonitorable) GetHelp() string       { return "" }
func (m *pingdomMonitorable) GetVariants() []string { return []string{config.DefaultVariant} }

func (m *pingdomMonitorable) Register(variant string, router router.MonitorableRouter, configManager monitorableConfig.Manager) bool {
	usecase := pingdomUsecase.NewPingdomUsecase()
	delivery := pingdomDelivery.NewPingdomDelivery(usecase)

	// RegisterTile route to echo
	route := router.Group("/pingdom", variant).GET("/check", delivery.GetCheck)

	// RegisterTile data for config hydration
	configManager.RegisterTile(pingdom.PingdomCheckTileType, variant, &pingdomModels.CheckParams{}, route.Path, config.DefaultInitialMaxDelay)
	configManager.DisableTile(pingdom.PingdomChecksTileType, variant)

	return true
}
