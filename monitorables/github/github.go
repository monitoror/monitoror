//+build !faker

package github

import (
	"fmt"
	"net/url"
	"time"

	uiConfig "github.com/monitoror/monitoror/api/config"
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
	"github.com/monitoror/monitoror/service/store"
)

type Monitorable struct {
	store *store.Store

	config map[coreModels.VariantName]*githubConfig.Github

	// Config tile settings
	countTypeSetting       uiConfig.TileEnabler
	checksTypeSetting      uiConfig.TileEnabler
	pullrequestTypeSetting uiConfig.TileGeneratorEnabler
}

func NewMonitorable(store *store.Store) *Monitorable {
	m := &Monitorable{}
	m.store = store
	m.config = make(map[coreModels.VariantName]*githubConfig.Github)

	// Load core config from env
	pkgMonitorable.LoadConfig(&m.config, githubConfig.Default)

	// Register Monitorable Tile in config manager
	m.countTypeSetting = store.TileSettingManager.Register(api.GithubCountTileType, versions.MinimalVersion, m.GetVariantNames())
	m.checksTypeSetting = store.TileSettingManager.Register(api.GithubChecksTileType, versions.MinimalVersion, m.GetVariantNames())
	m.pullrequestTypeSetting = store.TileSettingManager.RegisterGenerator(api.GithubChecksTileType, versions.MinimalVersion, m.GetVariantNames())

	return m
}

func (m *Monitorable) GetDisplayName() string {
	return "GitHub"
}

func (m *Monitorable) GetVariantNames() []coreModels.VariantName {
	return pkgMonitorable.GetVariants(m.config)
}

func (m *Monitorable) Validate(variantName coreModels.VariantName) (bool, error) {
	conf := m.config[variantName]

	// No configuration set
	if conf.URL == githubConfig.Default.URL && conf.Token == "" {
		return false, nil
	}

	// Error in URL
	if _, err := url.Parse(conf.URL); err != nil {
		return false, fmt.Errorf(`%s contains invalid URL: "%s"`, pkgMonitorable.BuildMonitorableEnvKey(conf, variantName, "URL"), conf.URL)
	}

	// Error in Token
	if conf.Token == "" {
		return false, fmt.Errorf(`%s is required, no value founds`, pkgMonitorable.BuildMonitorableEnvKey(conf, variantName, "TOKEN"))
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
	m.countTypeSetting.Enable(variant, &githubModels.CountParams{}, routeCount.Path, conf.InitialMaxDelay)
	m.checksTypeSetting.Enable(variant, &githubModels.ChecksParams{}, routeChecks.Path, conf.InitialMaxDelay)
	m.pullrequestTypeSetting.Enable(variant, &githubModels.PullRequestGeneratorParams{}, usecase.PullRequestsGenerator)
}
