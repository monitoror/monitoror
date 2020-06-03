//+build faker

package gitlab

import (
	"github.com/monitoror/monitoror/api/config/versions"
	"github.com/monitoror/monitoror/internal/pkg/monitorable"
	coreModels "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorables/gitlab/api"
	gitlabDelivery "github.com/monitoror/monitoror/monitorables/gitlab/api/delivery/http"
	gitlabModels "github.com/monitoror/monitoror/monitorables/gitlab/api/models"
	gitlabUsecase "github.com/monitoror/monitoror/monitorables/gitlab/api/usecase"
	"github.com/monitoror/monitoror/registry"
	"github.com/monitoror/monitoror/store"
)

type Monitorable struct {
	monitorable.DefaultMonitorableFaker

	store *store.Store

	// Config tile settings
	issuesTileEnabler       registry.TileEnabler
	pipelineTileEnabler     registry.TileEnabler
	mergeRequestTileEnabler registry.TileEnabler
}

func NewMonitorable(store *store.Store) *Monitorable {
	m := &Monitorable{}
	m.store = store

	// Register Monitorable Tile in config manager
	m.issuesTileEnabler = store.Registry.RegisterTile(api.GitlabIssuesTileType, versions.MinimalVersion, m.GetVariantsNames())
	m.pipelineTileEnabler = store.Registry.RegisterTile(api.GitlabPipelineTileType, versions.MinimalVersion, m.GetVariantsNames())
	m.mergeRequestTileEnabler = store.Registry.RegisterTile(api.GitlabMergeRequestTileType, versions.MinimalVersion, m.GetVariantsNames())

	return m
}

func (m *Monitorable) GetDisplayName() string {
	return "GitLab"
}

func (m *Monitorable) Enable(variantName coreModels.VariantName) {
	usecase := gitlabUsecase.NewGitlabUsecase()
	delivery := gitlabDelivery.NewGitlabDelivery(usecase)

	// EnableTile route to echo
	routeGroup := m.store.MonitorableRouter.Group("/gitlab", variantName)
	routeIssues := routeGroup.GET("/issues", delivery.GetIssues)
	routePipeline := routeGroup.GET("/pipeline", delivery.GetPipeline)
	routeMergeRequest := routeGroup.GET("/mergerequest", delivery.GetMergeRequest)

	// EnableTile data for config hydration
	m.issuesTileEnabler.Enable(variantName, &gitlabModels.IssuesParams{}, routeIssues.Path)
	m.pipelineTileEnabler.Enable(variantName, &gitlabModels.PipelineParams{}, routePipeline.Path)
	m.mergeRequestTileEnabler.Enable(variantName, &gitlabModels.MergeRequestParams{}, routeMergeRequest.Path)
}
