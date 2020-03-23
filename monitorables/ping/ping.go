//+build !faker

package ping

import (
	uiConfig "github.com/monitoror/monitoror/api/config/usecase"
	coreConfig "github.com/monitoror/monitoror/config"
	coreModels "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorables/ping/api"
	pingDelivery "github.com/monitoror/monitoror/monitorables/ping/api/delivery/http"
	pingModels "github.com/monitoror/monitoror/monitorables/ping/api/models"
	pingRepository "github.com/monitoror/monitoror/monitorables/ping/api/repository"
	pingUsecase "github.com/monitoror/monitoror/monitorables/ping/api/usecase"
	pingConfig "github.com/monitoror/monitoror/monitorables/ping/config"
	"github.com/monitoror/monitoror/pkg/monitoror/utils/system"
	"github.com/monitoror/monitoror/service/store"
)

type Monitorable struct {
	store *store.Store

	config map[coreModels.Variant]*pingConfig.Ping
}

func NewMonitorable(store *store.Store) *Monitorable {
	monitorable := &Monitorable{}
	monitorable.store = store
	monitorable.config = make(map[coreModels.Variant]*pingConfig.Ping)

	// Load core config from env
	coreConfig.LoadMonitorableConfig(&monitorable.config, pingConfig.Default)

	// Register Monitorable Tile in config manager
	store.UIConfigManager.RegisterTile(api.PingTileType, monitorable.GetVariants(), uiConfig.MinimalVersion)

	return monitorable
}

func (m *Monitorable) GetDisplayName() string {
	return "Ping"
}

func (m *Monitorable) GetVariants() []coreModels.Variant {
	return coreConfig.GetVariantsFromConfig(m.config)
}

func (m *Monitorable) Validate(_ coreModels.Variant) (bool, error) {
	return system.IsRawSocketAvailable(), nil
}

func (m *Monitorable) Enable(variant coreModels.Variant) {
	conf := m.config[variant]

	repository := pingRepository.NewPingRepository(conf)
	usecase := pingUsecase.NewPingUsecase(repository)
	delivery := pingDelivery.NewPingDelivery(usecase)

	// EnableTile route to echo
	route := m.store.MonitorableRouter.Group("/ping", variant).GET("/ping", delivery.GetPing)

	// EnableTile data for config hydration
	m.store.UIConfigManager.EnableTile(api.PingTileType, variant,
		&pingModels.PingParams{}, route.Path, conf.InitialMaxDelay)
}
