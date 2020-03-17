//+build !faker

package monitorable

import (
	"time"

	"github.com/monitoror/monitoror/config"
	monitorableConfig "github.com/monitoror/monitoror/monitorable/config"
	"github.com/monitoror/monitoror/monitorable/github"
	githubDelivery "github.com/monitoror/monitoror/monitorable/github/delivery/http"
	githubModels "github.com/monitoror/monitoror/monitorable/github/models"
	githubRepository "github.com/monitoror/monitoror/monitorable/github/repository"
	githubUsecase "github.com/monitoror/monitoror/monitorable/github/usecase"
	"github.com/monitoror/monitoror/service/options"
	"github.com/monitoror/monitoror/service/router"
)

type githubMonitorable struct {
	config map[string]*config.Github
}

func NewGithubMonitorable(config map[string]*config.Github) Monitorable {
	return &githubMonitorable{config: config}
}

func (m *githubMonitorable) GetHelp() string       { return "HEEEELLLPPPP" }
func (m *githubMonitorable) GetVariants() []string { return config.GetVariantsFromConfig(m.config) }
func (m *githubMonitorable) isEnabled(variant string) bool {
	conf := m.config[variant]
	return conf.Token != ""
}

func (m *githubMonitorable) Register(variant string, router router.MonitorableRouter, configManager monitorableConfig.Manager) bool {
	enabled := m.isEnabled(variant)
	if enabled {
		conf := m.config[variant]

		// Custom UpstreamCacheExpiration only for count because github has no-cache for this query and the rate limit is 30req/Hour
		countCacheExpiration := time.Millisecond * time.Duration(conf.CountCacheExpiration)

		repository := githubRepository.NewGithubRepository(conf)
		usecase := githubUsecase.NewGithubUsecase(repository)
		delivery := githubDelivery.NewGithubDelivery(usecase)

		// RegisterTile route to echo
		azureGroup := router.Group("/github", variant)
		routeCount := azureGroup.GET("/count", delivery.GetCount, options.WithCustomCacheExpiration(countCacheExpiration))
		routeChecks := azureGroup.GET("/checks", delivery.GetChecks)

		// RegisterTile data for config hydration
		configManager.RegisterTile(github.GithubCountTileType, variant, &githubModels.CountParams{}, routeCount.Path, conf.InitialMaxDelay)
		configManager.RegisterTile(github.GithubChecksTileType, variant, &githubModels.ChecksParams{}, routeChecks.Path, conf.InitialMaxDelay)
		configManager.RegisterDynamicTile(github.GithubPullRequestTileType, variant, &githubModels.PullRequestParams{}, usecase.PullRequests)
	} else {
		// RegisterTile data for config verify
		configManager.DisableTile(github.GithubCountTileType, variant)
		configManager.DisableTile(github.GithubChecksTileType, variant)
		configManager.DisableTile(github.GithubPullRequestTileType, variant)
	}

	return enabled
}
