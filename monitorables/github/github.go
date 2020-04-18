//+build !faker

package github

import (
	"time"

	"github.com/monitoror/monitoror/api/config/versions"
	pkgMonitorable "github.com/monitoror/monitoror/internal/pkg/monitorable"
	coreModels "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorables/github/api"
	githubDelivery "github.com/monitoror/monitoror/monitorables/github/api/delivery/http"
	githubModels "github.com/monitoror/monitoror/monitorables/github/api/models"
	githubRepository "github.com/monitoror/monitoror/monitorables/github/api/repository"
	githubUsecase "github.com/monitoror/monitoror/monitorables/github/api/usecase"
	githubConfig "github.com/monitoror/monitoror/monitorables/github/config"
	"github.com/monitoror/monitoror/service/options"
	"github.com/monitoror/monitoror/service/registry"
	"github.com/monitoror/monitoror/service/store"
)

type Monitorable struct {
	store *store.Store

	config map[coreModels.VariantName]*githubConfig.Github

	// Config tile settings
	countTileEnabler            registry.TileEnabler
	checksTileEnabler           registry.TileEnabler
	pullrequestGeneratorEnabler registry.GeneratorEnabler
}

func NewMonitorable(store *store.Store) *Monitorable {
	m := &Monitorable{}
	m.store = store
	m.config = make(map[coreModels.VariantName]*githubConfig.Github)

	// Load core config from env
	pkgMonitorable.LoadConfig(&m.config, githubConfig.Default)

	// Register Monitorable Tile in config manager
	m.countTileEnabler = store.Registry.RegisterTile(api.GithubCountTileType, versions.MinimalVersion, m.GetVariantsNames())
	m.checksTileEnabler = store.Registry.RegisterTile(api.GithubChecksTileType, versions.MinimalVersion, m.GetVariantsNames())
	m.pullrequestGeneratorEnabler = store.Registry.RegisterGenerator(api.GithubChecksTileType, versions.MinimalVersion, m.GetVariantsNames())

	return m
}

func (m *Monitorable) GetDisplayName() string {
	return "GitHub"
}

func (m *Monitorable) GetVariantsNames() []coreModels.VariantName {
	return pkgMonitorable.GetVariantsNames(m.config)
}

func (m *Monitorable) Validate(variantName coreModels.VariantName) (bool, []error) {
	conf := m.config[variantName]

	// No configuration set
	if conf.URL == githubConfig.Default.URL && conf.Token == "" {
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

	// Custom UpstreamCacheExpiration only for count because github has no-cache for this query and the rate limit is 30req/Hour
	countCacheExpiration := time.Millisecond * time.Duration(conf.CountCacheExpiration)

	repository := githubRepository.NewGithubRepository(conf)
	usecase := githubUsecase.NewGithubUsecase(repository)
	delivery := githubDelivery.NewGithubDelivery(usecase)

	// EnableTile route to echo
	routeGroup := m.store.MonitorableRouter.Group("/github", variantName)
	routeCount := routeGroup.GET("/count", delivery.GetCount, options.WithCustomCacheExpiration(countCacheExpiration))
	routeChecks := routeGroup.GET("/checks", delivery.GetChecks)

	// EnableTile data for config hydration
	m.countTileEnabler.Enable(variantName, &githubModels.CountParams{}, routeCount.Path)
	m.checksTileEnabler.Enable(variantName, &githubModels.ChecksParams{}, routeChecks.Path)
	m.pullrequestGeneratorEnabler.Enable(variantName, &githubModels.PullRequestGeneratorParams{}, usecase.PullRequestsGenerator)
}
