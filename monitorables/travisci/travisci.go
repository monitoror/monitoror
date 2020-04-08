//+build !faker

package travisci

import (
	"fmt"
	"net/url"

	uiConfig "github.com/monitoror/monitoror/api/config"
	"github.com/monitoror/monitoror/api/config/versions"
	pkgMonitorable "github.com/monitoror/monitoror/internal/pkg/monitorable"

	coreModels "github.com/monitoror/monitoror/models"

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

	config map[coreModels.VariantName]*travisciConfig.TravisCI

	// Config tile settings
	buildTileSetting uiConfig.TileEnabler
}

func NewMonitorable(store *store.Store) *Monitorable {
	m := &Monitorable{}
	m.store = store
	m.config = make(map[coreModels.VariantName]*travisciConfig.TravisCI)

	// Load core config from env
	pkgMonitorable.LoadConfig(&m.config, travisciConfig.Default)

	// Register Monitorable Tile in config manager
	m.buildTileSetting = store.TileSettingManager.Register(api.TravisCIBuildTileType, versions.MinimalVersion, m.GetVariants())

	return m
}

func (m *Monitorable) GetDisplayName() string {
	return "Travis CI"
}

func (m *Monitorable) GetVariants() []coreModels.VariantName {
	return pkgMonitorable.GetVariants(m.config)
}

func (m *Monitorable) Validate(variant coreModels.VariantName) (bool, error) {
	conf := m.config[variant]
	// Error in URL
	if _, err := url.Parse(conf.URL); err != nil {
		return false, fmt.Errorf(`%s contains invalid URL: "%s"`, pkgMonitorable.BuildMonitorableEnvKey(conf, variant, "URL"), conf.URL)
	}

	return true, nil
}

func (m *Monitorable) Enable(variant coreModels.VariantName) {
	conf := m.config[variant]

	repository := travisciRepository.NewTravisCIRepository(conf)
	usecase := travisciUsecase.NewTravisCIUsecase(repository)
	delivery := travisciDelivery.NewTravisCIDelivery(usecase)

	// EnableTile route to echo
	routeGroup := m.store.MonitorableRouter.Group("/travisci", variant)
	route := routeGroup.GET("/build", delivery.GetBuild)

	// EnableTile data for config hydration
	m.buildTileSetting.Enable(variant, &travisciModels.BuildParams{}, route.Path, conf.InitialMaxDelay)
}
