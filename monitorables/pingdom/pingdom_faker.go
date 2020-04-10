//+build faker

package pingdom

import (
	uiConfig "github.com/monitoror/monitoror/api/config"
	"github.com/monitoror/monitoror/api/config/versions"
	coreConfig "github.com/monitoror/monitoror/config"
	"github.com/monitoror/monitoror/internal/pkg/monitorable"
	coreModels "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorables/pingdom/api"
	pingdomDelivery "github.com/monitoror/monitoror/monitorables/pingdom/api/delivery/http"
	pingdomModels "github.com/monitoror/monitoror/monitorables/pingdom/api/models"
	pingdomUsecase "github.com/monitoror/monitoror/monitorables/pingdom/api/usecase"
	"github.com/monitoror/monitoror/service/store"
)

type Monitorable struct {
	monitorable.DefaultMonitorableFaker

	store *store.Store

	// Config tile settings
	checkTileSetting uiConfig.TileEnabler
}

func NewMonitorable(store *store.Store) *Monitorable {
	m := &Monitorable{}
	m.store = store

	// Register Monitorable Tile in config manager
	m.checkTileSetting = store.TileSettingManager.Register(api.PingdomCheckTileType, versions.MinimalVersion, m.GetVariantNames())

	return m
}

func (m *Monitorable) GetDisplayName() string { return "Pingdom (faker)" }

func (m *Monitorable) Enable(variantName coreModels.VariantName) {
	usecase := pingdomUsecase.NewPingdomUsecase()
	delivery := pingdomDelivery.NewPingdomDelivery(usecase)

	// EnableTile route to echo
	routeGroup := m.store.MonitorableRouter.Group("/pingdom", variantName)
	route := routeGroup.GET("/pingdom", delivery.GetCheck)

	// EnableTile data for config hydration
	m.checkTileSetting.Enable(variant, &pingdomModels.CheckParams{}, route.Path, coreConfig.DefaultInitialMaxDelay)
}
