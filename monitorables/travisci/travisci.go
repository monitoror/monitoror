//+build !faker

package travisci

import (
	"fmt"
	"net/url"

	pkgMonitorable "github.com/monitoror/monitoror/internal/pkg/monitorable"

	coreModels "github.com/monitoror/monitoror/models"

	uiConfig "github.com/monitoror/monitoror/api/config/usecase"
	"github.com/monitoror/monitoror/monitorables/travisci/api"
	travisciDelivery "github.com/monitoror/monitoror/monitorables/travisci/api/delivery/http"
	travisciModels "github.com/monitoror/monitoror/monitorables/travisci/api/models"
	travisciRepository "github.com/monitoror/monitoror/monitorables/travisci/api/repository"
	travisciUsecase "github.com/monitoror/monitoror/monitorables/travisci/api/usecase"
	travisciConfig "github.com/monitoror/monitoror/monitorables/travisci/config"
	"github.com/monitoror/monitoror/service/store"
)

type Monitorable struct {
	store *store.Store

	config map[coreModels.Variant]*travisciConfig.TravisCI
}

func NewMonitorable(store *store.Store) *Monitorable {
	monitorable := &Monitorable{}
	monitorable.store = store
	monitorable.config = make(map[coreModels.Variant]*travisciConfig.TravisCI)

	// Load core config from env
	pkgMonitorable.LoadConfig(&monitorable.config, travisciConfig.Default)

	// Register Monitorable Tile in config manager
	store.UIConfigManager.RegisterTile(api.TravisCIBuildTileType, monitorable.GetVariants(), uiConfig.MinimalVersion)

	return monitorable
}

func (m *Monitorable) GetDisplayName() string {
	return "Travis CI"
}

func (m *Monitorable) GetVariants() []coreModels.Variant {
	return pkgMonitorable.GetVariants(m.config)
}

func (m *Monitorable) Validate(variant coreModels.Variant) (bool, error) {
	conf := m.config[variant]
	// No configuration set
	if conf.URL == "" {
		return false, nil
	}

	// Error in URL
	if _, err := url.Parse(conf.URL); err != nil {
		return false, fmt.Errorf(`%s contains invalid URL: "%s"`, pkgMonitorable.GetEnvName(conf, variant, "URL"), conf.URL)
	}

	return true, nil
}

func (m *Monitorable) Enable(variant coreModels.Variant) {
	conf := m.config[variant]

	repository := travisciRepository.NewTravisCIRepository(conf)
	usecase := travisciUsecase.NewTravisCIUsecase(repository)
	delivery := travisciDelivery.NewTravisCIDelivery(usecase)

	// EnableTile route to echo
	route := m.store.MonitorableRouter.Group("/travisci", variant).GET("/build", delivery.GetBuild)

	// EnableTile data for config hydration
	m.store.UIConfigManager.EnableTile(api.TravisCIBuildTileType, variant,
		&travisciModels.BuildParams{}, route.Path, conf.InitialMaxDelay)
}
