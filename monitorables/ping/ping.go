//+build !faker

package ping

import (
	uiConfig "github.com/monitoror/monitoror/api/config"
	"github.com/monitoror/monitoror/api/config/versions"
	pkgMonitorable "github.com/monitoror/monitoror/internal/pkg/monitorable"
	coreModels "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorables/ping/api"
	pingDelivery "github.com/monitoror/monitoror/monitorables/ping/api/delivery/http"
	pingModels "github.com/monitoror/monitoror/monitorables/ping/api/models"
	pingRepository "github.com/monitoror/monitoror/monitorables/ping/api/repository"
	pingUsecase "github.com/monitoror/monitoror/monitorables/ping/api/usecase"
	pingConfig "github.com/monitoror/monitoror/monitorables/ping/config"
	"github.com/monitoror/monitoror/pkg/system"
	"github.com/monitoror/monitoror/service/store"
)

type Monitorable struct {
	store *store.Store

	config map[coreModels.VariantName]*pingConfig.Ping

	// Config tile settings
	pingTileSetting uiConfig.TileEnabler
}

func NewMonitorable(store *store.Store) *Monitorable {
	m := &Monitorable{}
	m.store = store
	m.config = make(map[coreModels.VariantName]*pingConfig.Ping)

	// Load core config from env
	pkgMonitorable.LoadConfig(&m.config, pingConfig.Default)

	// Register Monitorable Tile in config manager
	m.pingTileSetting = store.TileSettingManager.Register(api.PingTileType, versions.MinimalVersion, m.GetVariantNames())

	return m
}

func (m *Monitorable) GetDisplayName() string {
	return "Ping"
}

func (m *Monitorable) GetVariantNames() []coreModels.VariantName {
	return pkgMonitorable.GetVariants(m.config)
}

func (m *Monitorable) Validate(_ coreModels.VariantName) (bool, error) {
	return system.IsRawSocketAvailable(), nil
}

func (m *Monitorable) Enable(variantName coreModels.VariantName) {
	conf := m.config[variantName]

	repository := pingRepository.NewPingRepository(conf)
	usecase := pingUsecase.NewPingUsecase(repository)
	delivery := pingDelivery.NewPingDelivery(usecase)

	// EnableTile route to echo
	routeGroup := m.store.MonitorableRouter.Group("/ping", variantName)
	route := routeGroup.GET("/ping", delivery.GetPing)

	// EnableTile data for config hydration
	m.pingTileSetting.Enable(variant, &pingModels.PingParams{}, route.Path, conf.InitialMaxDelay)
}
