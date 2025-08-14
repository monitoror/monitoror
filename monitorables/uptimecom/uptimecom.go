//go:build !faker

package uptimecom

import (
	"github.com/monitoror/monitoror/api/config/versions"
	pkgMonitorable "github.com/monitoror/monitoror/internal/pkg/monitorable"
	coreModels "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorables/uptimecom/api"
	uptimecomDelivery "github.com/monitoror/monitoror/monitorables/uptimecom/api/delivery/http"
	uptimecomModels "github.com/monitoror/monitoror/monitorables/uptimecom/api/models"
	uptimecomRepository "github.com/monitoror/monitoror/monitorables/uptimecom/api/repository"
	uptimecomUsecase "github.com/monitoror/monitoror/monitorables/uptimecom/api/usecase"
	uptimecomConfig "github.com/monitoror/monitoror/monitorables/uptimecom/config"
	"github.com/monitoror/monitoror/registry"
	"github.com/monitoror/monitoror/store"
)

type Monitorable struct {
	store *store.Store

	config map[coreModels.VariantName]*uptimecomConfig.Uptimecom

	// Config tile settings
	checkTileEnabler      registry.TileEnabler
	checkGeneratorEnabler registry.GeneratorEnabler
}

func NewMonitorable(store *store.Store) *Monitorable {
	m := &Monitorable{}
	m.store = store
	m.config = make(map[coreModels.VariantName]*uptimecomConfig.Uptimecom)

	// Load core config from env
	pkgMonitorable.LoadConfig(&m.config, uptimecomConfig.Default)

	// Register Monitorable Tile in config manager
	m.checkTileEnabler = store.Registry.RegisterTile(api.UptimecomCheckTileType, versions.Version2001, m.GetVariantsNames())
	m.checkGeneratorEnabler = store.Registry.RegisterGenerator(api.UptimecomCheckTileType, versions.Version2001, m.GetVariantsNames())

	return m
}

func (m *Monitorable) GetDisplayName() string {
	return "Uptimecom"
}

func (m *Monitorable) GetVariantsNames() []coreModels.VariantName {
	return pkgMonitorable.GetVariantsNames(m.config)
}

func (m *Monitorable) Validate(variantName coreModels.VariantName) (bool, []error) {
	conf := m.config[variantName]

	// No configuration set
	if conf.URL == uptimecomConfig.Default.URL && conf.Token == "" {
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

	repository := uptimecomRepository.NewUptimecomRepository(conf)
	usecase := uptimecomUsecase.NewUptimecomUsecase(repository, m.store.CacheStore, conf.CacheExpiration)
	delivery := uptimecomDelivery.NewUptimecomDelivery(usecase)

	// EnableTile route to echo
	routeGroup := m.store.MonitorableRouter.Group("/uptimecom", variantName)
	checkRoute := routeGroup.GET("/check", delivery.GetCheck)

	// EnableTile data for config hydration
	m.checkTileEnabler.Enable(variantName, &uptimecomModels.CheckParams{}, checkRoute.Path)
	m.checkGeneratorEnabler.Enable(variantName, &uptimecomModels.CheckGeneratorParams{}, usecase.CheckGenerator)
}
