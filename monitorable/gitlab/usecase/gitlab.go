//+build !faker

package usecase

import (
	"fmt"
	"sort"
	"time"

	"github.com/monitoror/monitoror/pkg/monitoror/utils/hash"

	"github.com/AlekSi/pointer"

	"github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorable/gitlab"
	gitlabModels "github.com/monitoror/monitoror/monitorable/gitlab/models"
	"github.com/monitoror/monitoror/pkg/monitoror/builder"
	"github.com/monitoror/monitoror/pkg/monitoror/cache"
	"github.com/monitoror/monitoror/pkg/monitoror/utils/git"
)

type (
	gitlabUsecase struct {
		repository gitlab.Repository

		// builds cache. used for save small history of build for stats
		buildsCache *cache.BuildCache
	}
)

const buildCacheSize = 5

func NewGitlabUsecase(repository gitlab.Repository) gitlab.Usecase {
	return &gitlabUsecase{
		repository,
		cache.NewBuildCache(buildCacheSize),
	}
}

func (gu *gitlabUsecase) Count(params *gitlabModels.CountParams) (*models.Tile, error) {
	tile := models.NewTile(gitlab.GitlabCountTileType).WithValue(models.NumberUnit)
	tile.Label = params.Query

	count, err := gu.repository.GetCount(params.Query)
	if err != nil {
		return nil, &models.MonitororError{Err: err, Tile: tile, Message: "unable to find count or wrong query"}
	}

	tile.Status = models.SuccessStatus
	tile.Value.Values = append(tile.Value.Values, fmt.Sprintf("%d", count))

	return tile, nil
}

func (gu *gitlabUsecase) Pipelines(params *gitlabModels.PipelinesParams) (*models.Tile, error) {
	tile := models.NewTile(gitlab.GitlabPipelinesTileType).WithBuild()
	tile.Label = params.Repository
	tile.Build.Branch = pointer.ToString(git.HumanizeBranch(params.Ref))

	// Request
	pipelines, err := gu.repository.GetPipelines(params.Repository, params.Ref)
	if err != nil {
		return nil, &models.MonitororError{Err: err, Tile: tile, Message: "unable to find ref pipelines"}
	}
	if len(pipelines.Runs) == 0 {
		// Warning because request was correct but there is no build
		return nil, &models.MonitororError{Tile: tile, Message: "no ref pipelines found", ErrorStatus: models.UnknownStatus}
	}

	var startedAt, finishedAt *time.Time
	var duration int
	var id string
	tile.Status, startedAt, finishedAt, duration, id = computePipelines(pipelines)

	// Set Previous Status
	previousStatus := gu.buildsCache.GetPreviousStatus(params, id)
	if previousStatus != nil {
		tile.Build.PreviousStatus = *previousStatus
	} else {
		tile.Build.PreviousStatus = models.UnknownStatus
	}

	// Author
	if tile.Status == models.FailedStatus && pipelines.HeadCommit != "" {
		commit, err := gu.repository.GetCommit(params.Repository, pipelines.HeadCommit)
		if err == nil {
			tile.Build.Author = &models.Author{
				Name:      commit.Author.Name,
				AvatarURL: commit.Author.AvatarURL,
			}
		}
	}

	// StartedAt / FinishedAt
	tile.Build.StartedAt = startedAt
	if tile.Status != models.RunningStatus && tile.Status != models.QueuedStatus {
		tile.Build.FinishedAt = finishedAt
	}

	// Duration
	if tile.Status == models.RunningStatus {
		tile.Build.Duration = pointer.ToInt64(int64(duration))

		estimatedDuration := gu.buildsCache.GetEstimatedDuration(params)
		if estimatedDuration != nil {
			tile.Build.EstimatedDuration = pointer.ToInt64(int64(estimatedDuration.Seconds()))
		} else {
			tile.Build.EstimatedDuration = pointer.ToInt64(int64(0))
		}
	}

	// Cache Duration when success / failed
	if tile.Status == models.SuccessStatus || tile.Status == models.FailedStatus && duration > 0 {
		gu.buildsCache.Add(params, id, tile.Status, time.Second*time.Duration(duration))
	}

	return tile, nil
}

func (gu *gitlabUsecase) ListDynamicTile(params interface{}) ([]builder.Result, error) {
	prParams := params.(*gitlabModels.MergeRequestParams)

	mergeRequests, err := gu.repository.GetMergeRequests(prParams.Repository)
	if err != nil {
		return nil, &models.MonitororError{Err: err, Message: "unable to find merge request"}
	}

	var results []builder.Result
	for _, mergeRequest := range mergeRequests {
		p := make(map[string]interface{})
		p["repository"] = mergeRequest.Repository
		p["ref"] = mergeRequest.Ref

		results = append(results, builder.Result{
			TileType: gitlab.GitlabPipelinesTileType,
			Label:    fmt.Sprintf("MR#%d @ %s", mergeRequest.ID, mergeRequest.Repository),
			Params:   p,
		})
	}

	return results, nil
}

func computePipelines(pipelines *gitlabModels.Pipelines) (models.TileStatus, *time.Time, *time.Time, int, string) {
	var statuses []models.TileStatus
	var startedAt *time.Time = nil
	var finishedAt *time.Time = nil
	var duration int = 0
	var ids = ""

	sort.Slice(pipelines.Runs, func(i, j int) bool {
		return pipelines.Runs[i].CreatedAt.After(pipelines.Runs[j].CreatedAt)
	})

	for _, run := range pipelines.Runs {
		statuses = append(statuses, parseRun(&run))
		if startedAt == nil || (run.StartedAt != nil && startedAt.After(*run.StartedAt)) {
			startedAt = run.StartedAt
		}
		if finishedAt == nil || (run.FinishedAt != nil && finishedAt.Before(*run.FinishedAt)) {
			finishedAt = run.FinishedAt
		}

		duration = run.Duration
		ids = fmt.Sprintf("%s-%d", ids, run.ID)
	}

	ids = hash.GetMD5Hash(ids)
	if len(statuses) == 0 {
		return models.UnknownStatus, nil, nil, 0, ids
	}

	return statuses[0], startedAt, finishedAt, duration, ids
}

func parseRun(run *gitlabModels.Run) models.TileStatus {
	// Based on : https://developer.gitlab.com/v3/checks/runs/
	switch run.Status {
	case "running":
		return models.RunningStatus
	case "pending":
		return models.QueuedStatus
	case "success":
		return models.SuccessStatus
	case "failed":
		return models.FailedStatus
	case "canceled":
		return models.CanceledStatus
	case "skipped":
		return models.DisabledStatus
	}

	return models.UnknownStatus
}
