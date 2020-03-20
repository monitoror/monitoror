//+build !faker

package ping

import (
	uiConfig "github.com/monitoror/monitoror/api/config/usecase"
	coreConfig "github.com/monitoror/monitoror/config"
	"github.com/monitoror/monitoror/monitorables/ping/api"
	pingDelivery "github.com/monitoror/monitoror/monitorables/ping/api/delivery/http"
	pingModels "github.com/monitoror/monitoror/monitorables/ping/api/models"
	pingRepository "github.com/monitoror/monitoror/monitorables/ping/api/repository"
	pingUsecase "github.com/monitoror/monitoror/monitorables/ping/api/usecase"
	pingCoreConfig "github.com/monitoror/monitoror/monitorables/ping/config"
	"github.com/monitoror/monitoror/pkg/monitoror/utils/system"
	"github.com/monitoror/monitoror/service/store"
)

type Monitorable struct {
	store *store.Store

	config map[string]*pingCoreConfig.Ping
}

func NewMonitorable(store *store.Store) *Monitorable {
	monitorable := &Monitorable{}
	monitorable.store = store
	monitorable.config = make(map[string]*pingCoreConfig.Ping)

	// Load core config from env
	coreConfig.LoadMonitorableConfig(&monitorable.config, pingCoreConfig.Default)

	// Register Monitorable Tile in config manager
	store.UIConfigManager.RegisterTile(api.PingTileType, monitorable.GetVariants(), uiConfig.MinimalVersion)

	return monitorable
}

func (m *Monitorable) GetVariants() []string { return coreConfig.GetVariantsFromConfig(m.config) }
func (m *Monitorable) IsValid(_ string) bool { return system.IsRawSocketAvailable() }

func (m *Monitorable) Enable(variant string) {
	conf := m.config[variant]

	repository := pingRepository.NewPingRepository(conf)
	usecase := pingUsecase.NewPingUsecase(repository)
	delivery := pingDelivery.NewPingDelivery(usecase)

	// EnableTile route to echo
	route := m.store.MonitorableRouter.Group("/ping", variant).GET("/ping", delivery.GetPing)

	// EnableTile data for config hydration
	m.store.UIConfigManager.EnableTile(api.PingTileType, variant, &pingModels.PingParams{}, route.Path, conf.InitialMaxDelay)
}
