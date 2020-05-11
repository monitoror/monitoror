//+build !faker

package port

import (
	"github.com/monitoror/monitoror/api/config/versions"
	pkgMonitorable "github.com/monitoror/monitoror/internal/pkg/monitorable"
	coreModels "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorables/port/api"
	portDelivery "github.com/monitoror/monitoror/monitorables/port/api/delivery/http"
	portModels "github.com/monitoror/monitoror/monitorables/port/api/models"
	portRepository "github.com/monitoror/monitoror/monitorables/port/api/repository"
	portUsecase "github.com/monitoror/monitoror/monitorables/port/api/usecase"
	portConfig "github.com/monitoror/monitoror/monitorables/port/config"
	"github.com/monitoror/monitoror/registry"
	"github.com/monitoror/monitoror/store"
)

type Monitorable struct {
	store *store.Store

	config map[coreModels.VariantName]*portConfig.Port

	// Config tile settings
	portTileEnabler registry.TileEnabler
}

func NewMonitorable(store *store.Store) *Monitorable {
	m := &Monitorable{}
	m.store = store
	m.config = make(map[coreModels.VariantName]*portConfig.Port)

	// Load core config from env
	pkgMonitorable.LoadConfig(&m.config, portConfig.Default)

	// Register Monitorable Tile in config manager
	m.portTileEnabler = store.Registry.RegisterTile(api.PortTileType, versions.MinimalVersion, m.GetVariantsNames())

	return m
}

func (m *Monitorable) GetDisplayName() string {
	return "Port"
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

	return true, nil
}

func (m *Monitorable) Enable(variantName coreModels.VariantName) {
	conf := m.config[variantName]

	repository := portRepository.NewPortRepository(conf)
	usecase := portUsecase.NewPortUsecase(repository)
	delivery := portDelivery.NewPortDelivery(usecase)

	// EnableTile route to echo
	routeGroup := m.store.MonitorableRouter.Group("/port", variantName)
	route := routeGroup.GET("/port", delivery.GetPort)

	// EnableTile data for config hydration
	m.portTileEnabler.Enable(variantName, &portModels.PortParams{}, route.Path)
}
