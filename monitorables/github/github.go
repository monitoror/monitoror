//+build !faker

package github

import (
	"time"

	uiConfig "github.com/monitoror/monitoror/api/config/usecase"
	coreConfig "github.com/monitoror/monitoror/config"
	coreModels "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorables/github/api"
	githubDelivery "github.com/monitoror/monitoror/monitorables/github/api/delivery/http"
	githubModels "github.com/monitoror/monitoror/monitorables/github/api/models"
	githubRepository "github.com/monitoror/monitoror/monitorables/github/api/repository"
	githubUsecase "github.com/monitoror/monitoror/monitorables/github/api/usecase"
	githubConfig "github.com/monitoror/monitoror/monitorables/github/config"
	"github.com/monitoror/monitoror/service/options"
	"github.com/monitoror/monitoror/service/store"
)

type Monitorable struct {
	store *store.Store

	config map[coreModels.Variant]*githubConfig.Github
}

func NewMonitorable(store *store.Store) *Monitorable {
	monitorable := &Monitorable{}
	monitorable.store = store
	monitorable.config = make(map[coreModels.Variant]*githubConfig.Github)

	// Load core config from env
	coreConfig.LoadMonitorableConfig(&monitorable.config, githubConfig.Default)

	// Register Monitorable Tile in config manager
	store.UIConfigManager.RegisterTile(api.GithubCountTileType, monitorable.GetVariants(), uiConfig.MinimalVersion)
	store.UIConfigManager.RegisterTile(api.GithubChecksTileType, monitorable.GetVariants(), uiConfig.MinimalVersion)
	store.UIConfigManager.RegisterTile(api.GithubPullRequestTileType, monitorable.GetVariants(), uiConfig.MinimalVersion)

	return monitorable
}

func (m *Monitorable) GetDisplayName() string {
	return "GitHub"
}

func (m *Monitorable) GetVariants() []coreModels.Variant {
	return coreConfig.GetVariantsFromConfig(m.config)
}

func (m *Monitorable) Validate(variant coreModels.Variant) (bool, error) {
	conf := m.config[variant]

	// No configuration set
	if conf.Token == "" {
		return false, nil
	}

	return true, nil
}

func (m *Monitorable) Enable(variant coreModels.Variant) {
	conf := m.config[variant]

	// Custom UpstreamCacheExpiration only for count because github has no-cache for this query and the rate limit is 30req/Hour
	countCacheExpiration := time.Millisecond * time.Duration(conf.CountCacheExpiration)

	repository := githubRepository.NewGithubRepository(conf)
	usecase := githubUsecase.NewGithubUsecase(repository)
	delivery := githubDelivery.NewGithubDelivery(usecase)

	// EnableTile route to echo
	githubGroup := m.store.MonitorableRouter.Group("/github", variant)
	routeCount := githubGroup.GET("/count", delivery.GetCount, options.WithCustomCacheExpiration(countCacheExpiration))
	routeChecks := githubGroup.GET("/checks", delivery.GetChecks)

	// EnableTile data for config hydration
	m.store.UIConfigManager.EnableTile(api.GithubCountTileType, variant,
		&githubModels.CountParams{}, routeCount.Path, conf.InitialMaxDelay)
	m.store.UIConfigManager.EnableTile(api.GithubChecksTileType, variant,
		&githubModels.ChecksParams{}, routeChecks.Path, conf.InitialMaxDelay)
	m.store.UIConfigManager.EnableDynamicTile(api.GithubPullRequestTileType, variant,
		&githubModels.PullRequestParams{}, usecase.PullRequests)
}
