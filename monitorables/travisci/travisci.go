//+build !faker

package travisci

import (
	"net/url"

	uiConfig "github.com/monitoror/monitoror/api/config/usecase"
	coreConfig "github.com/monitoror/monitoror/config"
	"github.com/monitoror/monitoror/monitorables/travisci/api"
	travisciDelivery "github.com/monitoror/monitoror/monitorables/travisci/api/delivery/http"
	travisciModels "github.com/monitoror/monitoror/monitorables/travisci/api/models"
	travisciRepository "github.com/monitoror/monitoror/monitorables/travisci/api/repository"
	travisciUsecase "github.com/monitoror/monitoror/monitorables/travisci/api/usecase"
	travisciCoreConfig "github.com/monitoror/monitoror/monitorables/travisci/config"
	"github.com/monitoror/monitoror/service/store"
)

type Monitorable struct {
	store *store.Store

	config map[string]*travisciCoreConfig.TravisCI
}

func NewMonitorable(store *store.Store) *Monitorable {
	monitorable := &Monitorable{}
	monitorable.store = store
	monitorable.config = make(map[string]*travisciCoreConfig.TravisCI)

	// Load core config from env
	coreConfig.LoadMonitorableConfig(&monitorable.config, travisciCoreConfig.Default)

	// Register Monitorable Tile in config manager
	store.UIConfigManager.RegisterTile(api.TravisCIBuildTileType, monitorable.GetVariants(), uiConfig.MinimalVersion)

	return monitorable
}

func (m *Monitorable) GetVariants() []string { return coreConfig.GetVariantsFromConfig(m.config) }
func (m *Monitorable) IsValid(variant string) bool {
	conf := m.config[variant]

	if conf.URL == "" {
		return false
	}

	if _, err := url.Parse(conf.URL); err != nil {
		return false
	}

	return true
}

func (m *Monitorable) Enable(variant string) {
	conf := m.config[variant]

	repository := travisciRepository.NewTravisCIRepository(conf)
	usecase := travisciUsecase.NewTravisCIUsecase(repository)
	delivery := travisciDelivery.NewTravisCIDelivery(usecase)

	// EnableTile route to echo
	route := m.store.MonitorableRouter.Group("/travisci", variant).GET("/build", delivery.GetBuild)

	// EnableTile data for config hydration
	m.store.UIConfigManager.EnableTile(api.TravisCIBuildTileType, variant, &travisciModels.BuildParams{}, route.Path, conf.InitialMaxDelay)
}
