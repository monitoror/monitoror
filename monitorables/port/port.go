//+build !faker

package port

import (
	uiConfig "github.com/monitoror/monitoror/api/config/usecase"
	coreConfig "github.com/monitoror/monitoror/config"
	"github.com/monitoror/monitoror/monitorables/port/api"
	portDelivery "github.com/monitoror/monitoror/monitorables/port/api/delivery/http"
	portModels "github.com/monitoror/monitoror/monitorables/port/api/models"
	portRepository "github.com/monitoror/monitoror/monitorables/port/api/repository"
	portUsecase "github.com/monitoror/monitoror/monitorables/port/api/usecase"
	portCoreConfig "github.com/monitoror/monitoror/monitorables/port/config"
	"github.com/monitoror/monitoror/service/store"
)

type Monitorable struct {
	store *store.Store

	config map[string]*portCoreConfig.Port
}

func NewMonitorable(store *store.Store) *Monitorable {
	monitorable := &Monitorable{}
	monitorable.store = store
	monitorable.config = make(map[string]*portCoreConfig.Port)

	// Load core config from env
	coreConfig.LoadMonitorableConfig(&monitorable.config, portCoreConfig.Default)

	// Register Monitorable Tile in config manager
	store.UIConfigManager.RegisterTile(api.PortTileType, monitorable.GetVariants(), uiConfig.MinimalVersion)

	return monitorable
}

func (m *Monitorable) GetVariants() []string { return coreConfig.GetVariantsFromConfig(m.config) }
func (m *Monitorable) IsValid(_ string) bool { return true }

func (m *Monitorable) Enable(variant string) {
	conf := m.config[variant]

	repository := portRepository.NewPortRepository(conf)
	usecase := portUsecase.NewPortUsecase(repository)
	delivery := portDelivery.NewPortDelivery(usecase)

	// EnableTile route to echo
	route := m.store.MonitorableRouter.Group("/port", variant).GET("/port", delivery.GetPort)

	// EnableTile data for config hydration
	m.store.UIConfigManager.EnableTile(api.PortTileType, variant, &portModels.PortParams{}, route.Path, conf.InitialMaxDelay)
}
