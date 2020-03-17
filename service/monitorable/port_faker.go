//+build faker

package monitorable

import (
	"github.com/monitoror/monitoror/config"
	monitorableConfig "github.com/monitoror/monitoror/monitorable/config"
	"github.com/monitoror/monitoror/monitorable/port"
	portDelivery "github.com/monitoror/monitoror/monitorable/port/delivery/http"
	portModels "github.com/monitoror/monitoror/monitorable/port/models"
	portUsecase "github.com/monitoror/monitoror/monitorable/port/usecase"
	"github.com/monitoror/monitoror/service/router"
)

type portMonitorable struct{}

func NewPortMonitorable(_ map[string]*config.Port) Monitorable {
	return &portMonitorable{}
}

func (m *portMonitorable) GetHelp() string       { return "" }
func (m *portMonitorable) GetVariants() []string { return []string{config.DefaultVariant} }

func (m *portMonitorable) Register(variant string, router router.MonitorableRouter, configManager monitorableConfig.Manager) bool {
	usecase := portUsecase.NewPortUsecase()
	delivery := portDelivery.NewPortDelivery(usecase)

	// RegisterTile route to echo
	route := router.Group("/port", variant).GET("/port", delivery.GetPort)

	// RegisterTile data for config hydration
	configManager.RegisterTile(port.PortTileType, variant, &portModels.PortParams{}, route.Path, config.DefaultInitialMaxDelay)

	return true
}
