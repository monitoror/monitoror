//+build !faker

package ping

import (
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
	"github.com/monitoror/monitoror/registry"
	"github.com/monitoror/monitoror/store"
)

type Monitorable struct {
	store *store.Store

	config map[coreModels.VariantName]*pingConfig.Ping

	// Config tile settings
	pingTileEnabler registry.TileEnabler
}

func NewMonitorable(store *store.Store) *Monitorable {
	m := &Monitorable{}
	m.store = store
	m.config = make(map[coreModels.VariantName]*pingConfig.Ping)

	// Load core config from env
	pkgMonitorable.LoadConfig(&m.config, pingConfig.Default)

	// Register Monitorable Tile in config manager
	m.pingTileEnabler = store.Registry.RegisterTile(api.PingTileType, versions.MinimalVersion, m.GetVariantsNames())

	return m
}

func (m *Monitorable) GetDisplayName() string {
	return "Ping"
}

func (m *Monitorable) GetVariantsNames() []coreModels.VariantName {
	return pkgMonitorable.GetVariantsNames(m.config)
}

func (m *Monitorable) Validate(variantName coreModels.VariantName) (bool, []error) {
	conf := m.config[variantName]

	// Validate Config
	if errors := pkgMonitorable.ValidateConfig(conf, variantName); errors != nil {
		return false, errors
	}

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
	m.pingTileEnabler.Enable(variantName, &pingModels.PingParams{}, route.Path)
}
