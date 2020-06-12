//+build !faker

package gitlab

import (
	"github.com/monitoror/monitoror/api/config/versions"
	pkgMonitorable "github.com/monitoror/monitoror/internal/pkg/monitorable"
	coreModels "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorables/gitlab/api"
	gitlabDelivery "github.com/monitoror/monitoror/monitorables/gitlab/api/delivery/http"
	gitlabModels "github.com/monitoror/monitoror/monitorables/gitlab/api/models"
	gitlabRepository "github.com/monitoror/monitoror/monitorables/gitlab/api/repository"
	gitlabUsecase "github.com/monitoror/monitoror/monitorables/gitlab/api/usecase"
	gitlabConfig "github.com/monitoror/monitoror/monitorables/gitlab/config"
	"github.com/monitoror/monitoror/registry"
	"github.com/monitoror/monitoror/store"
)

type Monitorable struct {
	store *store.Store

	config map[coreModels.VariantName]*gitlabConfig.Gitlab

	// Config tile settings
	countIssuesTileEnabler       registry.TileEnabler
	pipelineTileEnabler          registry.TileEnabler
	mergeRequestTileEnabler      registry.TileEnabler
	mergeRequestGeneratorEnabler registry.GeneratorEnabler
}

func NewMonitorable(store *store.Store) *Monitorable {
	m := &Monitorable{}
	m.store = store
	m.config = make(map[coreModels.VariantName]*gitlabConfig.Gitlab)

	// Load core config from env
	pkgMonitorable.LoadConfig(&m.config, gitlabConfig.Default)

	// Register Monitorable Tile in config manager
	m.countIssuesTileEnabler = store.Registry.RegisterTile(api.GitlabCountIssuesTileType, versions.MinimalVersion, m.GetVariantsNames())
	m.pipelineTileEnabler = store.Registry.RegisterTile(api.GitlabPipelineTileType, versions.MinimalVersion, m.GetVariantsNames())
	m.mergeRequestTileEnabler = store.Registry.RegisterTile(api.GitlabMergeRequestTileType, versions.MinimalVersion, m.GetVariantsNames())
	m.mergeRequestGeneratorEnabler = store.Registry.RegisterGenerator(api.GitlabMergeRequestTileType, versions.MinimalVersion, m.GetVariantsNames())

	return m
}

func (m *Monitorable) GetDisplayName() string {
	return "GitLab"
}

func (m *Monitorable) GetVariantsNames() []coreModels.VariantName {
	return pkgMonitorable.GetVariantsNames(m.config)
}

func (m *Monitorable) Validate(variantName coreModels.VariantName) (bool, []error) {
	conf := m.config[variantName]

	// No configuration set
	if conf.URL == gitlabConfig.Default.URL && conf.Token == "" {
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

	repository := gitlabRepository.NewGitlabRepository(conf)
	usecase := gitlabUsecase.NewGitlabUsecase(repository, m.store.CacheStore)
	delivery := gitlabDelivery.NewGitlabDelivery(usecase)

	// EnableTile route to echo
	routeGroup := m.store.MonitorableRouter.Group("/gitlab", variantName)
	routeCountIssues := routeGroup.GET("/count-issues", delivery.GetCountIssues)
	routePipeline := routeGroup.GET("/pipeline", delivery.GetPipeline)
	routeMergeRequest := routeGroup.GET("/mergerequest", delivery.GetMergeRequest)

	// EnableTile data for config hydration
	m.countIssuesTileEnabler.Enable(variantName, &gitlabModels.IssuesParams{}, routeCountIssues.Path)
	m.pipelineTileEnabler.Enable(variantName, &gitlabModels.PipelineParams{}, routePipeline.Path)
	m.mergeRequestTileEnabler.Enable(variantName, &gitlabModels.MergeRequestParams{}, routeMergeRequest.Path)
	m.mergeRequestGeneratorEnabler.Enable(variantName, &gitlabModels.MergeRequestGeneratorParams{}, usecase.MergeRequestsGenerator)
}
