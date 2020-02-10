//+build !faker

package usecase

import (
	"fmt"
	"sort"
	"time"

	"github.com/monitoror/monitoror/pkg/monitoror/utils/hash"

	"github.com/AlekSi/pointer"

	"github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorable/github"
	githubModels "github.com/monitoror/monitoror/monitorable/github/models"
	"github.com/monitoror/monitoror/pkg/monitoror/builder"
	"github.com/monitoror/monitoror/pkg/monitoror/cache"
	"github.com/monitoror/monitoror/pkg/monitoror/utils/git"
)

type (
	githubUsecase struct {
		repository github.Repository

		// builds cache. used for save small history of build for stats
		buildsCache *cache.BuildCache
	}
)

var orderedTileStatus = map[models.TileStatus]int{
	models.RunningStatus:        0,
	models.FailedStatus:         1,
	models.WarningStatus:        2,
	models.CanceledStatus:       3,
	models.ActionRequiredStatus: 4,
	models.QueuedStatus:         5,
	models.SuccessStatus:        6,
	models.DisabledStatus:       7,
	models.UnknownStatus:        8,
}

const buildCacheSize = 5

func NewGithubUsecase(repository github.Repository) github.Usecase {
	return &githubUsecase{
		repository,
		cache.NewBuildCache(buildCacheSize),
	}
}

func (gu *githubUsecase) Count(params *githubModels.CountParams) (*models.Tile, error) {
	tile := models.NewTile(github.GithubCountTileType)
	tile.Label = params.Query

	count, err := gu.repository.GetCount(params.Query)
	if err != nil {
		return nil, &models.MonitororError{Err: err, Tile: tile, Message: "unable to find count or wrong query"}
	}

	tile.Status = models.SuccessStatus
	tile.Values = []float64{float64(count)}

	return tile, nil
}

func (gu *githubUsecase) Checks(params *githubModels.ChecksParams) (*models.Tile, error) {
	tile := models.NewTile(github.GithubChecksTileType)
	tile.Label = fmt.Sprintf("%s\n%s", params.Repository, git.HumanizeBranch(params.Ref))

	// Request
	checks, err := gu.repository.GetChecks(params.Owner, params.Repository, params.Ref)
	if err != nil {
		return nil, &models.MonitororError{Err: err, Tile: tile, Message: "unable to find ref checks"}
	}
	if len(checks.Statuses) == 0 && len(checks.Runs) == 0 {
		// Warning because request was correct but there is no build
		return nil, &models.MonitororError{Tile: tile, Message: "no ref checks found", ErrorStatus: models.UnknownStatus}
	}

	var startedAt, finishedAt *time.Time
	var id string
	tile.Status, startedAt, finishedAt, id = computeChecks(checks)

	// Set Previous Status
	previousStatus := gu.buildsCache.GetPreviousStatus(params, id)
	if previousStatus != nil {
		tile.PreviousStatus = *previousStatus
	} else {
		tile.PreviousStatus = models.UnknownStatus
	}

	// Author
	if tile.Status == models.FailedStatus && checks.HeadCommit != nil {
		commit, err := gu.repository.GetCommit(params.Owner, params.Repository, *checks.HeadCommit)
		if err == nil {
			tile.Author = &models.Author{
				Name:      commit.Author.Name,
				AvatarURL: commit.Author.AvatarURL,
			}
		}
	}

	// StartedAt / FinishedAt
	tile.StartedAt = startedAt
	if tile.Status != models.RunningStatus && tile.Status != models.QueuedStatus {
		tile.FinishedAt = finishedAt
	}

	// Duration
	if tile.Status == models.RunningStatus {
		tile.Duration = pointer.ToInt64(int64(time.Since(*tile.StartedAt).Seconds()))

		estimatedDuration := gu.buildsCache.GetEstimatedDuration(params)
		if estimatedDuration != nil {
			tile.EstimatedDuration = pointer.ToInt64(int64(estimatedDuration.Seconds()))
		} else {
			tile.EstimatedDuration = pointer.ToInt64(int64(0))
		}
	}

	// Cache Duration when success / failed
	if tile.Status == models.SuccessStatus || tile.Status == models.FailedStatus || tile.Status == models.WarningStatus {
		gu.buildsCache.Add(params, id, tile.Status, tile.FinishedAt.Sub(*tile.StartedAt))
	}

	return tile, nil
}

func (gu *githubUsecase) ListDynamicTile(params interface{}) ([]builder.Result, error) {
	prParams := params.(*githubModels.PullRequestParams)

	pullRequests, err := gu.repository.GetPullRequests(prParams.Owner, prParams.Repository)
	if err != nil {
		return nil, &models.MonitororError{Err: err, Message: "unable to find pull request"}
	}

	var results []builder.Result
	for _, pullRequest := range pullRequests {
		p := make(map[string]interface{})
		p["owner"] = pullRequest.Owner
		p["repository"] = pullRequest.Repository
		p["ref"] = pullRequest.Ref

		results = append(results, builder.Result{
			TileType: github.GithubChecksTileType,
			Label:    fmt.Sprintf("%s\n%s", pullRequest.Repository, pullRequest.Title),
			Params:   p,
		})
	}

	return results, nil
}

func computeChecks(refStatus *githubModels.Checks) (models.TileStatus, *time.Time, *time.Time, string) {
	var statuses []models.TileStatus
	var startedAt *time.Time = nil
	var finishedAt *time.Time = nil
	var ids = ""

	for _, run := range refStatus.Runs {
		statuses = append(statuses, parseRun(&run))
		if startedAt == nil || (run.StartedAt != nil && startedAt.After(*run.StartedAt)) {
			startedAt = run.StartedAt
		}
		if finishedAt == nil || (run.CompletedAt != nil && finishedAt.Before(*run.CompletedAt)) {
			finishedAt = run.CompletedAt
		}
		ids = fmt.Sprintf("%s-%d", ids, run.ID)
	}

	// Sort statues by created date and save every title to remove duplicate statues
	// Some app add new status with the same name each time status change
	sort.Slice(refStatus.Statuses, func(i, j int) bool {
		return refStatus.Statuses[i].CreatedAt.After(refStatus.Statuses[j].CreatedAt)
	})

	titles := make(map[string]bool)
	for _, status := range refStatus.Statuses {
		if _, ok := titles[status.Title]; !ok {
			statuses = append(statuses, parseStatus(&status))
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

	ids = hash.GetMD5Hash(ids)
	if len(statuses) == 0 {
		return models.UnknownStatus, nil, nil, ids
	}

	return statuses[0], startedAt, finishedAt, ids
}

func parseRun(run *githubModels.Run) models.TileStatus {
	// Based on : https://developer.github.com/v3/checks/runs/
	switch run.Status {
	case "in_progress":
		return models.RunningStatus
	case "queued":
		return models.QueuedStatus
	case "completed":
		switch run.Conclusion {
		case "success":
			return models.SuccessStatus
		case "failure":
			return models.FailedStatus
		case "timed_out":
			return models.FailedStatus
		case "neutral":
			return models.WarningStatus
		case "cancelled":
			return models.CanceledStatus
		case "action_required":
			return models.ActionRequiredStatus
		}
	}

	return models.UnknownStatus
}

func parseStatus(status *githubModels.Status) models.TileStatus {
	// Based on : https://developer.github.com/v3/repos/statuses/
	switch status.State {
	case "success":
		return models.SuccessStatus
	case "failure":
		return models.FailedStatus
	case "error":
		return models.FailedStatus
	case "pending":
		return models.RunningStatus
	}

	return models.UnknownStatus
}
