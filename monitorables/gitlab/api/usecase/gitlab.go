//+build !faker

package usecase

import (
	"fmt"
	"time"

	uiConfigModels "github.com/monitoror/monitoror/api/config/models"
	monitorableCache "github.com/monitoror/monitoror/internal/pkg/monitorable/cache"
	coreModels "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorables/gitlab/api"
	"github.com/monitoror/monitoror/monitorables/gitlab/api/models"
	"github.com/monitoror/monitoror/pkg/git"

	"github.com/AlekSi/pointer"
	"github.com/jsdidierlaurent/echo-middleware/cache"
	uuid "github.com/satori/go.uuid"
)

type (
	gitlabUsecase struct {
		repository api.Repository
		// Used to generate store key by repository
		repositoryUID string

		// store is used to store persistent data (project, merge requests)
		store cache.Store

		// builds cache. used for save small history of build for stats
		buildsCache *monitorableCache.BuildCache
	}
)

const (
	buildCacheSize = 5

	projectCacheExpiration      = cache.NEVER
	mergeRequestCacheExpiration = time.Second * 30

	GitlabProjectStoreKeyPrefix      = "monitoror.gitlab.project.store"
	GitlabMergeRequestStoreKeyPrefix = "monitoror.gitlab.mergeRequest.store"
)

func NewGitlabUsecase(repository api.Repository, store cache.Store) api.Usecase {
	return &gitlabUsecase{
		repository:    repository,
		repositoryUID: uuid.NewV4().String(),
		store:         store,
		buildsCache:   monitorableCache.NewBuildCache(buildCacheSize),
	}
}

func (gu *gitlabUsecase) CountIssues(params *models.IssuesParams) (*coreModels.Tile, error) {
	tile := coreModels.NewTile(api.GitlabCountIssuesTileType).WithMetrics(coreModels.NumberUnit)
	tile.Label = "GitLab count"

	count, err := gu.repository.GetCountIssues(params)
	if err != nil {
		return nil, &coreModels.MonitororError{Err: err, Tile: tile, Message: "unable to load issues"}
	}

	tile.Status = coreModels.SuccessStatus
	tile.Metrics.Values = append(tile.Metrics.Values, fmt.Sprintf("%d", count))

	return tile, nil
}

func (gu *gitlabUsecase) Pipeline(params *models.PipelineParams) (*coreModels.Tile, error) {
	tile := coreModels.NewTile(api.GitlabPipelineTileType).WithBuild()
	tile.Label = fmt.Sprintf("%d", params.ProjectID)
	tile.Build.Branch = pointer.ToString(git.HumanizeBranch(params.Ref))

	// Load Project and cache it
	project, err := gu.getProject(*params.ProjectID)
	if err != nil {
		return nil, &coreModels.MonitororError{Err: err, Tile: tile, Message: "unable to load project"}
	}
	tile.Label = project.Repository

	// Load pipelines for given ref
	pipelines, err := gu.repository.GetPipelines(*params.ProjectID, params.Ref)
	if err != nil {
		return nil, &coreModels.MonitororError{Err: err, Tile: tile, Message: "unable to load pipelines"}
	}
	if len(pipelines) == 0 {
		// Warning because request was correct but there is no build
		return nil, &coreModels.MonitororError{Tile: tile, Message: "no pipelines found", ErrorStatus: coreModels.UnknownStatus}
	}

	// Load pipeline detail
	pipeline, err := gu.repository.GetPipeline(*params.ProjectID, pipelines[0])
	if err != nil {
		return nil, &coreModels.MonitororError{Err: err, Tile: tile, Message: "unable to load pipeline"}
	}

	gu.computePipeline(params, tile, pipeline)

	// Author
	if tile.Status == coreModels.FailedStatus {
		tile.Build.Author = &pipeline.Author
	}

	return tile, nil
}

func (gu *gitlabUsecase) MergeRequest(params *models.MergeRequestParams) (*coreModels.Tile, error) {
	tile := coreModels.NewTile(api.GitlabMergeRequestTileType).WithBuild()
	tile.Label = fmt.Sprintf("%d", params.ProjectID)

	// Load Project and cache it
	project, err := gu.getProject(*params.ProjectID)
	if err != nil {
		return nil, &coreModels.MonitororError{Err: err, Tile: tile, Message: "unable to load project"}
	}
	tile.Label = project.Repository

	// Load MergeRequest
	mergeRequest, err := gu.getMergeRequest(*params.ProjectID, *params.ID)
	if err != nil {
		return nil, &coreModels.MonitororError{Err: err, Tile: tile, Message: "unable to load merge request"}
	}

	// Load MergeRequest project
	mergeRequestProject, err := gu.getProject(mergeRequest.SourceProjectID)
	if err != nil {
		return nil, &coreModels.MonitororError{Err: err, Tile: tile, Message: "unable to load project"}
	}

	tile.Build.Branch = pointer.ToString(git.HumanizeBranch(mergeRequest.SourceBranch))
	if project.Owner != mergeRequestProject.Owner {
		tile.Build.Branch = pointer.ToString(fmt.Sprintf("%s:%s", mergeRequestProject.Owner, *tile.Build.Branch))
	}
	tile.Build.MergeRequest = &coreModels.TileMergeRequest{
		ID:    mergeRequest.ID,
		Title: mergeRequest.Title,
	}

	// Load merge request pipelines
	pipelines, err := gu.repository.GetMergeRequestPipelines(*params.ProjectID, mergeRequest.ID)
	if err != nil {
		return nil, &coreModels.MonitororError{Err: err, Tile: tile, Message: "unable to load pipelines"}
	}
	if len(pipelines) == 0 {
		// Warning because request was correct but there is no build
		return nil, &coreModels.MonitororError{Tile: tile, Message: "no pipelines found", ErrorStatus: coreModels.UnknownStatus}
	}

	// Load pipeline detail
	pipeline, err := gu.repository.GetPipeline(*params.ProjectID, pipelines[0])
	if err != nil {
		return nil, &coreModels.MonitororError{Err: err, Tile: tile, Message: "unable to load pipeline"}
	}

	gu.computePipeline(params, tile, pipeline)

	// Author
	if tile.Status == coreModels.FailedStatus {
		tile.Build.Author = &mergeRequest.Author
	}

	return tile, nil
}

func (gu *gitlabUsecase) computePipeline(params interface{}, tile *coreModels.Tile, pipeline *models.Pipeline) {
	tile.Status = parseStatus(pipeline.Status)

	// Set Previous Status
	strPipelineID := fmt.Sprintf("%d", pipeline.ID)
	previousStatus := gu.buildsCache.GetPreviousStatus(params, strPipelineID)
	if previousStatus != nil {
		tile.Build.PreviousStatus = *previousStatus
	} else {
		tile.Build.PreviousStatus = coreModels.UnknownStatus
	}

	// StartedAt / FinishedAt
	tile.Build.StartedAt = pipeline.StartedAt
	if tile.Status != coreModels.RunningStatus && tile.Status != coreModels.QueuedStatus {
		tile.Build.FinishedAt = pipeline.FinishedAt
	}

	// Duration
	if tile.Status == coreModels.RunningStatus {
		tile.Build.Duration = pointer.ToInt64(int64(time.Since(*tile.Build.StartedAt).Seconds()))

		estimatedDuration := gu.buildsCache.GetEstimatedDuration(params)
		if estimatedDuration != nil {
			tile.Build.EstimatedDuration = pointer.ToInt64(int64(estimatedDuration.Seconds()))
		} else {
			tile.Build.EstimatedDuration = pointer.ToInt64(int64(0))
		}
	}

	// Cache Duration when success / failed
	if tile.Status == coreModels.SuccessStatus || tile.Status == coreModels.FailedStatus {
		// In case of build without StartedAt ...
		if tile.Build.StartedAt != nil && tile.Build.FinishedAt != nil {
			gu.buildsCache.Add(params, strPipelineID, tile.Status, tile.Build.FinishedAt.Sub(*tile.Build.StartedAt))
		}
	}
}

func (gu *gitlabUsecase) MergeRequestsGenerator(params interface{}) ([]uiConfigModels.GeneratedTile, error) {
	prParams := params.(*models.MergeRequestGeneratorParams)

	mergeRequests, err := gu.repository.GetMergeRequests(*prParams.ProjectID)
	if err != nil {
		return nil, &coreModels.MonitororError{Err: err, Message: "unable to load merge requests"}
	}

	var results []uiConfigModels.GeneratedTile
	for _, mergeRequest := range mergeRequests {
		p := &models.MergeRequestParams{}
		p.ProjectID = prParams.ProjectID
		p.ID = pointer.ToInt(mergeRequest.ID)

		results = append(results, uiConfigModels.GeneratedTile{
			Params: p,
		})

		// Add merge request into store
		_ = gu.store.Set(gu.getMergeRequestStoreKey(*prParams.ProjectID, mergeRequest.ID), mergeRequest, mergeRequestCacheExpiration)
	}

	return results, nil
}

func (gu *gitlabUsecase) getProjectStoreKey(projectID int) string {
	return fmt.Sprintf("%s:%s-%d", GitlabProjectStoreKeyPrefix, gu.repositoryUID, projectID)
}

// getProject load project information (from cache or api) and add result in cache
func (gu *gitlabUsecase) getProject(projectID int) (*models.Project, error) {
	project := &models.Project{}

	storeKey := gu.getProjectStoreKey(projectID)
	if err := gu.store.Get(storeKey, project); err != nil {
		if project, err = gu.repository.GetProject(projectID); err != nil {
			return nil, err
		}

		_ = gu.store.Set(storeKey, *project, projectCacheExpiration)
	}

	return project, nil
}

func (gu *gitlabUsecase) getMergeRequestStoreKey(projectID int, mergeRequestID int) string {
	return fmt.Sprintf("%s:%s-%d-%d", GitlabMergeRequestStoreKeyPrefix, gu.repositoryUID, projectID, mergeRequestID)
}

// getMergeRequest load merge request information (from cache or api) and add result in cache
func (gu *gitlabUsecase) getMergeRequest(projectID int, mergeRequestID int) (*models.MergeRequest, error) {
	mergeRequest := &models.MergeRequest{}

	storeKey := gu.getMergeRequestStoreKey(projectID, mergeRequestID)
	if err := gu.store.Get(storeKey, mergeRequest); err != nil {
		if mergeRequest, err = gu.repository.GetMergeRequest(projectID, mergeRequestID); err != nil {
			return nil, err
		}

		_ = gu.store.Set(storeKey, *mergeRequest, mergeRequestCacheExpiration)
	}

	return mergeRequest, nil
}

func parseStatus(status string) coreModels.TileStatus {
	// See: https://docs.gitlab.com/ee/api/pipelines.html#list-project-pipelines
	switch status {
	case "running":
		return coreModels.RunningStatus
	case "pending":
		return coreModels.QueuedStatus
	case "success":
		return coreModels.SuccessStatus
	case "failed":
		return coreModels.FailedStatus
	case "canceled":
		return coreModels.CanceledStatus
	case "skipped":
		return coreModels.CanceledStatus
	case "created":
		return coreModels.QueuedStatus
	case "manual":
		return coreModels.ActionRequiredStatus
	default:
		return coreModels.UnknownStatus
	}
}
