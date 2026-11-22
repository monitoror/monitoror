//go:build !faker

package gitea

import (
	"github.com/monitoror/monitoror/api/config/versions"
	pkgMonitorable "github.com/monitoror/monitoror/internal/pkg/monitorable"
	coreModels "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorables/gitea/api"
	giteaDelivery "github.com/monitoror/monitoror/monitorables/gitea/api/delivery/http"
	giteaModels "github.com/monitoror/monitoror/monitorables/gitea/api/models"
	giteaRepository "github.com/monitoror/monitoror/monitorables/gitea/api/repository"
	giteaUsecase "github.com/monitoror/monitoror/monitorables/gitea/api/usecase"
	giteaConfig "github.com/monitoror/monitoror/monitorables/gitea/config"
	"github.com/monitoror/monitoror/registry"
	"github.com/monitoror/monitoror/store"
)

type Monitorable struct {
	store *store.Store

	config map[coreModels.VariantName]*giteaConfig.Gitea

	// Config tile settings
	countTileEnabler            registry.TileEnabler
	checksTileEnabler           registry.TileEnabler
	pullRequestTileEnabler      registry.TileEnabler
	pullRequestGeneratorEnabler registry.GeneratorEnabler
}

func NewMonitorable(store *store.Store) *Monitorable {
	m := &Monitorable{}
	m.store = store
	m.config = make(map[coreModels.VariantName]*giteaConfig.Gitea)

	// Load core config from env
	pkgMonitorable.LoadConfig(&m.config, giteaConfig.Default)

	// Register Monitorable Tile in config manager
	m.countTileEnabler = store.Registry.RegisterTile(api.GiteaCountTileType, versions.Version2001, m.GetVariantsNames())
	m.checksTileEnabler = store.Registry.RegisterTile(api.GiteaChecksTileType, versions.Version2001, m.GetVariantsNames())
	m.pullRequestTileEnabler = store.Registry.RegisterTile(api.GiteaPullRequestTileType, versions.Version2001, m.GetVariantsNames())
	m.pullRequestGeneratorEnabler = store.Registry.RegisterGenerator(api.GiteaPullRequestTileType, versions.Version2001, m.GetVariantsNames())

	return m
}

func (m *Monitorable) GetDisplayName() string {
	return "Gitea"
}

func (m *Monitorable) GetVariantsNames() []coreModels.VariantName {
	return pkgMonitorable.GetVariantsNames(m.config)
}

func (m *Monitorable) Validate(variantName coreModels.VariantName) (bool, []error) {
	conf := m.config[variantName]

	// No configuration set
	if conf.URL == giteaConfig.Default.URL && conf.Token == "" {
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

	repository := giteaRepository.NewGiteaRepository(conf)
	usecase := giteaUsecase.NewGiteaUsecase(repository)
	delivery := giteaDelivery.NewGiteaDelivery(usecase)

	// EnableTile route to echo
	routeGroup := m.store.MonitorableRouter.Group("/gitea", variantName)
	routeCount := routeGroup.GET("/count", delivery.GetCount)
	routeChecks := routeGroup.GET("/checks", delivery.GetChecks)
	routePullRequest := routeGroup.GET("/pullrequest", delivery.GetPullRequest)

	// EnableTile data for config hydration
	m.countTileEnabler.Enable(variantName, &giteaModels.CountParams{}, routeCount.Path)
	m.checksTileEnabler.Enable(variantName, &giteaModels.ChecksParams{}, routeChecks.Path)
	m.pullRequestTileEnabler.Enable(variantName, &giteaModels.PullRequestParams{}, routePullRequest.Path)
	m.pullRequestGeneratorEnabler.Enable(variantName, &giteaModels.PullRequestGeneratorParams{}, usecase.PullRequestsGenerator)
}
