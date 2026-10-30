//go:build !faker

package oneuptime

import (
	"github.com/monitoror/monitoror/api/config/versions"
	pkgMonitorable "github.com/monitoror/monitoror/internal/pkg/monitorable"
	coreModels "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorables/oneuptime/api"
	oneuptimeDelivery "github.com/monitoror/monitoror/monitorables/oneuptime/api/delivery/http"
	oneuptimeModels "github.com/monitoror/monitoror/monitorables/oneuptime/api/models"
	oneuptimeRepository "github.com/monitoror/monitoror/monitorables/oneuptime/api/repository"
	oneuptimeUsecase "github.com/monitoror/monitoror/monitorables/oneuptime/api/usecase"
	oneuptimeConfig "github.com/monitoror/monitoror/monitorables/oneuptime/config"
	"github.com/monitoror/monitoror/registry"
	"github.com/monitoror/monitoror/store"
)

type Monitorable struct {
	store *store.Store

	config map[coreModels.VariantName]*oneuptimeConfig.Oneuptime

	// Config tile settings
	monitorTileEnabler      registry.TileEnabler
	monitorGeneratorEnabler registry.GeneratorEnabler
}

func NewMonitorable(store *store.Store) *Monitorable {
	m := &Monitorable{}
	m.store = store
	m.config = make(map[coreModels.VariantName]*oneuptimeConfig.Oneuptime)

	// Load core config from env
	pkgMonitorable.LoadConfig(&m.config, oneuptimeConfig.Default)

	// Register Monitorable Tile in config manager
	m.monitorTileEnabler = store.Registry.RegisterTile(api.OneuptimeMonitorTileType, versions.Version2003, m.GetVariantsNames())
	m.monitorGeneratorEnabler = store.Registry.RegisterGenerator(api.OneuptimeMonitorTileType, versions.Version2003, m.GetVariantsNames())

	return m
}

func (m *Monitorable) GetDisplayName() string {
	return "OneUptime"
}

func (m *Monitorable) GetVariantsNames() []coreModels.VariantName {
	return pkgMonitorable.GetVariantsNames(m.config)
}

func (m *Monitorable) Validate(variantName coreModels.VariantName) (bool, []error) {
	conf := m.config[variantName]

	// No configuration set
	if conf.ProjectID == "" && conf.Token == "" {
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

	repository := oneuptimeRepository.NewOneuptimeRepository(conf)
	usecase := oneuptimeUsecase.NewOneuptimeUsecase(repository, m.store.CacheStore, conf.CacheExpiration)
	delivery := oneuptimeDelivery.NewOneuptimeDelivery(usecase)

	// EnableTile route to echo
	routeGroup := m.store.MonitorableRouter.Group("/oneuptime", variantName)
	monitorRoute := routeGroup.GET("/monitor", delivery.GetMonitor)

	// EnableTile data for config hydration
	m.monitorTileEnabler.Enable(variantName, &oneuptimeModels.MonitorParams{}, monitorRoute.Path)
	m.monitorGeneratorEnabler.Enable(variantName, &oneuptimeModels.MonitorGeneratorParams{}, usecase.MonitorGenerator)
}
