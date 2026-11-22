//go:build !faker

package usecase

import (
	"fmt"
	"sort"
	"time"

	uiConfigModels "github.com/monitoror/monitoror/api/config/models"
	"github.com/monitoror/monitoror/internal/pkg/monitorable/cache"
	coreModels "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorables/gitea/api"
	"github.com/monitoror/monitoror/monitorables/gitea/api/models"
	"github.com/monitoror/monitoror/pkg/git"
	"github.com/monitoror/monitoror/pkg/hash"

	"github.com/AlekSi/pointer"
)

type (
	giteaUsecase struct {
		repository api.Repository

		// builds cache. used for save small history of build for stats
		buildsCache *cache.BuildCache
	}
)

var orderedTileStatus = map[coreModels.TileStatus]int{
	coreModels.RunningStatus:        0,
	coreModels.FailedStatus:         1,
	coreModels.WarningStatus:        2,
	coreModels.CanceledStatus:       3,
	coreModels.ActionRequiredStatus: 4,
	coreModels.QueuedStatus:         5,
	coreModels.SuccessStatus:        6,
	coreModels.DisabledStatus:       7,
	coreModels.UnknownStatus:        8,
}

const buildCacheSize = 5

func NewGiteaUsecase(repository api.Repository) api.Usecase {
	return &giteaUsecase{
		repository,
		cache.NewBuildCache(buildCacheSize),
	}
}

func (gu *giteaUsecase) Count(params *models.CountParams) (*coreModels.Tile, error) {
	tile := coreModels.NewTile(api.GiteaCountTileType).WithMetrics(coreModels.NumberUnit)
	tile.Label = "Gitea count"

	count, err := gu.repository.GetCount(params)
	if err != nil {
		return nil, &coreModels.MonitororError{Err: err, Tile: tile, Message: "unable to load count"}
	}

	tile.Status = coreModels.SuccessStatus
	tile.Metrics.Values = append(tile.Metrics.Values, fmt.Sprintf("%d", count))

	return tile, nil
}

func (gu *giteaUsecase) Checks(params *models.ChecksParams) (*coreModels.Tile, error) {
	tile := coreModels.NewTile(api.GiteaChecksTileType).WithBuild()
	tile.Label = params.Repository
	tile.Build.Branch = pointer.ToString(git.HumanizeBranch(params.Ref))

	// Request checks
	checks, err := gu.repository.GetChecks(params.Owner, params.Repository, params.Ref)
	if err != nil {
		return nil, &coreModels.MonitororError{Err: err, Tile: tile, Message: "unable to load ref checks"}
	}
	if len(checks.Statuses) == 0 {
		// Unknown because request was correct but there is no build
		return nil, &coreModels.MonitororError{Tile: tile, Message: "no ref checks found", ErrorStatus: coreModels.UnknownStatus}
	}

	// Compute checks into tile
	gu.computeRefChecks(tile, checks, params.String())

	// Author of last commit
	if tile.Status == coreModels.FailedStatus && checks.HeadCommit != nil {
		commit, err := gu.repository.GetCommit(params.Owner, params.Repository, *checks.HeadCommit)
		if err == nil {
			tile.Build.Author = &coreModels.Author{
				Name:      commit.Author.Name,
				AvatarURL: commit.Author.AvatarURL,
			}
		}
	}

	return tile, nil
}

func (gu *giteaUsecase) PullRequest(params *models.PullRequestParams) (*coreModels.Tile, error) {
	tile := coreModels.NewTile(api.GiteaPullRequestTileType).WithBuild()

	// Request pullRequest
	pullRequest, err := gu.repository.GetPullRequest(params.Owner, params.Repository, *params.ID)
	if err != nil {
		return nil, &coreModels.MonitororError{Err: err, Tile: tile, Message: "unable to load pull request"}
	}

	tile.Label = params.Repository
	tile.Build.Branch = pointer.ToString(git.HumanizeBranch(pullRequest.SourceBranch))
	if pullRequest.SourceOwner != "" && params.Owner != pullRequest.SourceOwner {
		tile.Build.Branch = pointer.ToString(fmt.Sprintf("%s:%s", pullRequest.SourceOwner, *tile.Build.Branch))
	}
	tile.Build.MergeRequest = &coreModels.TileMergeRequest{
		ID:    pullRequest.ID,
		Title: pullRequest.Title,
	}

	// With pull request, we use CommitSHA as Ref on params.Owner/params.Repo to handle correctly Checks from Forks
	checks, err := gu.repository.GetChecks(params.Owner, params.Repository, pullRequest.CommitSHA)
	if err != nil {
		return nil, &coreModels.MonitororError{Err: err, Tile: tile, Message: "unable to load ref checks"}
	}

	// Compute checks into tile
	gu.computeRefChecks(tile, checks, params.String())

	// Author of pull request
	if tile.Status == coreModels.FailedStatus {
		tile.Build.Author = &pullRequest.Author
	}

	return tile, nil
}

func (gu *giteaUsecase) PullRequestsGenerator(params interface{}) ([]uiConfigModels.GeneratedTile, error) {
	prParams := params.(*models.PullRequestGeneratorParams)

	pullRequests, err := gu.repository.GetPullRequests(prParams.Owner, prParams.Repository)
	if err != nil {
		return nil, &coreModels.MonitororError{Err: err, Message: "unable to load pull request"}
	}

	var results []uiConfigModels.GeneratedTile
	for _, pullRequest := range pullRequests {
		p := &models.PullRequestParams{}
		p.Owner = prParams.Owner
		p.Repository = prParams.Repository
		p.ID = pointer.ToInt(pullRequest.ID)

		results = append(results, uiConfigModels.GeneratedTile{
			Params: p,
		})
	}

	return results, nil
}

func (gu *giteaUsecase) computeRefChecks(tile *coreModels.Tile, checks *models.Checks, paramsKey string) {
	// convert checks
	statuses, startedAt, finishedAt, id := convertChecks(checks)

	// if statuses is empty, return SuccessStatus by default
	if len(statuses) == 0 {
		tile.Status = coreModels.SuccessStatus
		tile.Build.PreviousStatus = coreModels.UnknownStatus

		return
	}

	tile.Status = statuses[0]

	// Set Previous Status
	previousStatus := gu.buildsCache.GetPreviousStatus(paramsKey, id)
	if previousStatus != nil {
		tile.Build.PreviousStatus = *previousStatus
	} else {
		tile.Build.PreviousStatus = coreModels.UnknownStatus
	}

	// StartedAt / FinishedAt
	tile.Build.StartedAt = startedAt
	if tile.Status != coreModels.RunningStatus && tile.Status != coreModels.QueuedStatus {
		tile.Build.FinishedAt = finishedAt
	}

	// Duration
	if tile.Status == coreModels.RunningStatus {
		tile.Build.Duration = pointer.ToInt64(int64(time.Since(*tile.Build.StartedAt).Seconds()))

		estimatedDuration := gu.buildsCache.GetEstimatedDuration(paramsKey)
		if estimatedDuration != nil {
			tile.Build.EstimatedDuration = pointer.ToInt64(int64(estimatedDuration.Seconds()))
		} else {
			tile.Build.EstimatedDuration = pointer.ToInt64(int64(0))
		}
	}

	// Cache Duration when success / failed
	if tile.Status == coreModels.SuccessStatus || tile.Status == coreModels.FailedStatus || tile.Status == coreModels.WarningStatus {
		gu.buildsCache.Add(paramsKey, id, tile.Status, tile.Build.FinishedAt.Sub(*tile.Build.StartedAt))
	}
}

// convertChecks transform models.Checks to use it in computeRefChecks
func convertChecks(checks *models.Checks) ([]coreModels.TileStatus, *time.Time, *time.Time, string) {
	var statuses []coreModels.TileStatus
	var startedAt *time.Time = nil
	var finishedAt *time.Time = nil
	var ids = ""

	// Sort statues from newest to oldest and save every title to remove duplicate statues
	// Gitea adds a new status each time a status changes. Timestamps have a one second precision,
	// so the highest ID wins when two statuses share the same date.
	sort.Slice(checks.Statuses, func(i, j int) bool {
		if checks.Statuses[i].CreatedAt.Equal(checks.Statuses[j].CreatedAt) {
			return checks.Statuses[i].ID > checks.Statuses[j].ID
		}
		return checks.Statuses[i].CreatedAt.After(checks.Statuses[j].CreatedAt)
	})

	titles := make(map[string]bool)
	for i := range checks.Statuses {
		status := &checks.Statuses[i]
		if _, ok := titles[status.Title]; !ok {
			statuses = append(statuses, parseStatus(status))
			titles[status.Title] = true
		}

		if startedAt == nil || startedAt.After(status.CreatedAt) {
			startedAt = &status.CreatedAt
		}
		if finishedAt == nil || finishedAt.Before(status.UpdatedAt) {
			finishedAt = &status.UpdatedAt
		}
		ids = fmt.Sprintf("%s-%d", ids, status.ID)
	}

	sort.Slice(statuses, func(i, j int) bool {
		return orderedTileStatus[statuses[i]] < orderedTileStatus[statuses[j]]
	})

	return statuses, startedAt, finishedAt, hash.GetMD5Hash(ids)
}

func parseStatus(status *models.Status) coreModels.TileStatus {
	// Based on Gitea CommitStatusState : pending, success, error, failure, warning, skipped
	switch status.State {
	case "success":
		return coreModels.SuccessStatus
	case "failure", "error":
		return coreModels.FailedStatus
	case "warning":
		return coreModels.WarningStatus
	case "pending":
		return coreModels.RunningStatus
	case "skipped":
		return coreModels.DisabledStatus
	}

	return coreModels.UnknownStatus
}
