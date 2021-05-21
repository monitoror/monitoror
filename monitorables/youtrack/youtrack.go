//+build !faker

package youtrack

import (
	"github.com/monitoror/monitoror/api/config/versions"
	pkgMonitorable "github.com/monitoror/monitoror/internal/pkg/monitorable"
	coreModels "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorables/youtrack/api"
	youtrackDelivery "github.com/monitoror/monitoror/monitorables/youtrack/api/delivery/http"
	youtrackModels "github.com/monitoror/monitoror/monitorables/youtrack/api/models"
	youtrackRepository "github.com/monitoror/monitoror/monitorables/youtrack/api/repository"
	youtrackUsecase "github.com/monitoror/monitoror/monitorables/youtrack/api/usecase"
	youtrackConfig "github.com/monitoror/monitoror/monitorables/youtrack/config"
	"github.com/monitoror/monitoror/registry"
	"github.com/monitoror/monitoror/store"
)

type Monitorable struct {
	store *store.Store

	config map[coreModels.VariantName]*youtrackConfig.Youtrack

	// Config tile settings
	countIssuesTileEnabler registry.TileEnabler
}

func NewMonitorable(store *store.Store) *Monitorable {
	m := &Monitorable{}
	m.store = store
	m.config = make(map[coreModels.VariantName]*youtrackConfig.Youtrack)

	// Load core config from env
	pkgMonitorable.LoadConfig(&m.config, youtrackConfig.Default)

	// Register Monitorable Tile in config manager
	m.countIssuesTileEnabler = store.Registry.RegisterTile(api.YoutrackCountIssuesTileType, versions.Version2001, m.GetVariantsNames())

	return m
}

func (m *Monitorable) GetDisplayName() string {
	return "Youtrack"
}

func (m *Monitorable) GetVariantsNames() []coreModels.VariantName {
	return pkgMonitorable.GetVariantsNames(m.config)
}

func (m *Monitorable) Validate(variantName coreModels.VariantName) (bool, []error) {
	conf := m.config[variantName]

	// No configuration set
	if conf.URL == youtrackConfig.Default.URL && conf.Token == "" {
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

	repository := youtrackRepository.NewYoutrackRepository(conf)
	usecase := youtrackUsecase.NewYoutrackUsecase(repository)
	delivery := youtrackDelivery.NewYoutrackDelivery(usecase)

	// EnableTile route to echo
	routeGroup := m.store.MonitorableRouter.Group("/youtrack", variantName)
	routeCountIssues := routeGroup.GET("/count-issues", delivery.GetCountIssues)

	// EnableTile data for config hydration
	m.countIssuesTileEnabler.Enable(variantName, &youtrackModels.IssuesCountParams{}, routeCountIssues.Path)
}
