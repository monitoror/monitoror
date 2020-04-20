//+build !faker

package pingdom

import (
	"github.com/monitoror/monitoror/api/config/versions"
	pkgMonitorable "github.com/monitoror/monitoror/internal/pkg/monitorable"
	"github.com/monitoror/monitoror/service/registry"

	coreModels "github.com/monitoror/monitoror/models"

	"github.com/monitoror/monitoror/monitorables/pingdom/api"
	pingdomDelivery "github.com/monitoror/monitoror/monitorables/pingdom/api/delivery/http"
	pingdomModels "github.com/monitoror/monitoror/monitorables/pingdom/api/models"
	pingdomRepository "github.com/monitoror/monitoror/monitorables/pingdom/api/repository"
	pingdomUsecase "github.com/monitoror/monitoror/monitorables/pingdom/api/usecase"
	pingdomConfig "github.com/monitoror/monitoror/monitorables/pingdom/config"
	"github.com/monitoror/monitoror/service/store"
)

type Monitorable struct {
	store *store.Store

	config map[coreModels.VariantName]*pingdomConfig.Pingdom

	// Config tile settings
	checkTileEnabler      registry.TileEnabler
	checkGeneratorEnabler registry.GeneratorEnabler
}

func NewMonitorable(store *store.Store) *Monitorable {
	m := &Monitorable{}
	m.store = store
	m.config = make(map[coreModels.VariantName]*pingdomConfig.Pingdom)

	// Load core config from env
	pkgMonitorable.LoadConfig(&m.config, pingdomConfig.Default)

	// Register Monitorable Tile in config manager
	m.checkTileEnabler = store.Registry.RegisterTile(api.PingdomCheckTileType, versions.MinimalVersion, m.GetVariantsNames())
	m.checkGeneratorEnabler = store.Registry.RegisterGenerator(api.PingdomCheckTileType, versions.MinimalVersion, m.GetVariantsNames())

	return m
}

func (m *Monitorable) GetDisplayName() string {
	return "Pingdom"
}

func (m *Monitorable) GetVariantsNames() []coreModels.VariantName {
	return pkgMonitorable.GetVariantsNames(m.config)
}

func (m *Monitorable) Validate(variantName coreModels.VariantName) (bool, []error) {
	conf := m.config[variantName]

	// No configuration set
	if conf.URL == pingdomConfig.Default.URL && conf.Token == "" {
		return false, nil
	}

	// Validate Config
	if errors := pkgMonitorable.ValidateConfig(conf, variantName); errors != nil {
		return false, errors
	}

	return true, nil
}

func (m *Monitorable) Enable(variantName coreModels.VariantName) {
	conf := m.config[variantName]

	repository := pingdomRepository.NewPingdomRepository(conf)
	usecase := pingdomUsecase.NewPingdomUsecase(repository, m.store.CacheStore, conf.CacheExpiration)
	delivery := pingdomDelivery.NewPingdomDelivery(usecase)

	// EnableTile route to echo
	routeGroup := m.store.MonitorableRouter.Group("/pingdom", variantName)
	route := routeGroup.GET("/pingdom", delivery.GetCheck)

	// EnableTile data for config hydration
	m.checkTileEnabler.Enable(variantName, &pingdomModels.CheckParams{}, route.Path)
	m.checkGeneratorEnabler.Enable(variantName, &pingdomModels.CheckGeneratorParams{}, usecase.CheckGenerator)
}
