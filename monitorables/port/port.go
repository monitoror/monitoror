//+build !faker

package port

import (
	uiConfig "github.com/monitoror/monitoror/api/config"
	"github.com/monitoror/monitoror/api/config/versions"
	pkgMonitorable "github.com/monitoror/monitoror/internal/pkg/monitorable"
	coreModels "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorables/port/api"
	portDelivery "github.com/monitoror/monitoror/monitorables/port/api/delivery/http"
	portModels "github.com/monitoror/monitoror/monitorables/port/api/models"
	portRepository "github.com/monitoror/monitoror/monitorables/port/api/repository"
	portUsecase "github.com/monitoror/monitoror/monitorables/port/api/usecase"
	portConfig "github.com/monitoror/monitoror/monitorables/port/config"
	"github.com/monitoror/monitoror/service/store"
)

type Monitorable struct {
	store *store.Store

	config map[coreModels.VariantName]*portConfig.Port

	// Config tile settings
	portTileSetting uiConfig.TileEnabler
}

func NewMonitorable(store *store.Store) *Monitorable {
	m := &Monitorable{}
	m.store = store
	m.config = make(map[coreModels.VariantName]*portConfig.Port)

	// Load core config from env
	pkgMonitorable.LoadConfig(&m.config, portConfig.Default)

	// Register Monitorable Tile in config manager
	m.portTileSetting = store.TileSettingManager.Register(api.PortTileType, versions.MinimalVersion, m.GetVariants())

	return m
}

func (m *Monitorable) GetDisplayName() string {
	return "Port"
}

func (m *Monitorable) GetVariants() []coreModels.VariantName {
	return pkgMonitorable.GetVariants(m.config)
}

func (m *Monitorable) Validate(_ coreModels.VariantName) (bool, error) {
	return true, nil
}

func (m *Monitorable) Enable(variant coreModels.VariantName) {
	conf := m.config[variant]

	repository := portRepository.NewPortRepository(conf)
	usecase := portUsecase.NewPortUsecase(repository)
	delivery := portDelivery.NewPortDelivery(usecase)

	// EnableTile route to echo
	routeGroup := m.store.MonitorableRouter.Group("/port", variant)
	route := routeGroup.GET("/port", delivery.GetPort)

	// EnableTile data for config hydration
	m.portTileSetting.Enable(variant, &portModels.PortParams{}, route.Path, conf.InitialMaxDelay)
}
