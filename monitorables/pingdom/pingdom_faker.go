//+build faker

package pingdom

import (
	"github.com/monitoror/monitoror/api/config/versions"
	"github.com/monitoror/monitoror/internal/pkg/monitorable"
	coreModels "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorables/pingdom/api"
	pingdomDelivery "github.com/monitoror/monitoror/monitorables/pingdom/api/delivery/http"
	pingdomModels "github.com/monitoror/monitoror/monitorables/pingdom/api/models"
	pingdomUsecase "github.com/monitoror/monitoror/monitorables/pingdom/api/usecase"
	"github.com/monitoror/monitoror/registry"
	"github.com/monitoror/monitoror/store"
)

type Monitorable struct {
	monitorable.DefaultMonitorableFaker

	store *store.Store

	// Config tile settings
	checkTileEnabler            registry.TileEnabler
	transactionCheckTileEnabler registry.TileEnabler
}

func NewMonitorable(store *store.Store) *Monitorable {
	m := &Monitorable{}
	m.store = store

	// Register Monitorable Tile in config manager
	m.checkTileEnabler = store.Registry.RegisterTile(api.PingdomCheckTileType, versions.MinimalVersion, m.GetVariantsNames())
	m.transactionCheckTileEnabler = store.Registry.RegisterTile(api.PingdomTransactionCheckTileType, versions.MinimalVersion, m.GetVariantsNames())

	return m
}

func (m *Monitorable) GetDisplayName() string { return "Pingdom" }

func (m *Monitorable) Enable(variantName coreModels.VariantName) {
	usecase := pingdomUsecase.NewPingdomUsecase()
	delivery := pingdomDelivery.NewPingdomDelivery(usecase)

	// EnableTile route to echo
	routeGroup := m.store.MonitorableRouter.Group("/pingdom", variantName)
	checkRoute := routeGroup.GET("/check", delivery.GetCheck)
	transactionCheckRoute := routeGroup.GET("/transaction-check", delivery.GetTransactionCheck)

	// EnableTile data for config hydration
	m.checkTileEnabler.Enable(variantName, &pingdomModels.CheckParams{}, checkRoute.Path)
	m.transactionCheckTileEnabler.Enable(variantName, &pingdomModels.TransactionCheckParams{}, transactionCheckRoute.Path)
}
