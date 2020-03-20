//+build !faker

package usecase

import (
	"fmt"
	"sort"
	"time"

	coreModels "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorables/github/api"
	"github.com/monitoror/monitoror/monitorables/github/api/models"
	"github.com/monitoror/monitoror/pkg/monitoror/builder"
	"github.com/monitoror/monitoror/pkg/monitoror/cache"
	"github.com/monitoror/monitoror/pkg/monitoror/utils/git"
	"github.com/monitoror/monitoror/pkg/monitoror/utils/hash"

	"github.com/AlekSi/pointer"
)

type (
	githubUsecase struct {
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

func NewGithubUsecase(repository api.Repository) api.Usecase {
	return &githubUsecase{
		repository,
		cache.NewBuildCache(buildCacheSize),
	}
}

func (gu *githubUsecase) Count(params *models.CountParams) (*coreModels.Tile, error) {
	tile := coreModels.NewTile(api.GithubCountTileType).WithValue(coreModels.NumberUnit)
	tile.Label = params.Query

	count, err := gu.repository.GetCount(params.Query)
	if err != nil {
		return nil, &coreModels.MonitororError{Err: err, Tile: tile, Message: "unable to find count or wrong query"}
	}

	tile.Status = coreModels.SuccessStatus
	tile.Value.Values = append(tile.Value.Values, fmt.Sprintf("%d", count))

	return tile, nil
}

func (gu *githubUsecase) Checks(params *models.ChecksParams) (*coreModels.Tile, error) {
	tile := coreModels.NewTile(api.GithubChecksTileType).WithBuild()
	tile.Label = params.Repository
	tile.Build.Branch = pointer.ToString(git.HumanizeBranch(params.Ref))

	// Request
	checks, err := gu.repository.GetChecks(params.Owner, params.Repository, params.Ref)
	if err != nil {
		return nil, &coreModels.MonitororError{Err: err, Tile: tile, Message: "unable to find ref checks"}
	}
	if len(checks.Statuses) == 0 && len(checks.Runs) == 0 {
		// Warning because request was correct but there is no build
		return nil, &coreModels.MonitororError{Tile: tile, Message: "no ref checks found", ErrorStatus: coreModels.UnknownStatus}
	}

	var startedAt, finishedAt *time.Time
	var id string
	tile.Status, startedAt, finishedAt, id = computeChecks(checks)

	// Set Previous Status
	previousStatus := gu.buildsCache.GetPreviousStatus(params, id)
	if previousStatus != nil {
		tile.Build.PreviousStatus = *previousStatus
	} else {
		tile.Build.PreviousStatus = coreModels.UnknownStatus
	}

	// Author
	if tile.Status == coreModels.FailedStatus && checks.HeadCommit != nil {
		commit, err := gu.repository.GetCommit(params.Owner, params.Repository, *checks.HeadCommit)
		if err == nil {
			tile.Build.Author = &coreModels.Author{
				Name:      commit.Author.Name,
				AvatarURL: commit.Author.AvatarURL,
			}
		}
	}

	// StartedAt / FinishedAt
	tile.Build.StartedAt = startedAt
	if tile.Status != coreModels.RunningStatus && tile.Status != coreModels.QueuedStatus {
		tile.Build.FinishedAt = finishedAt
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
	if tile.Status == coreModels.SuccessStatus || tile.Status == coreModels.FailedStatus || tile.Status == coreModels.WarningStatus {
		gu.buildsCache.Add(params, id, tile.Status, tile.Build.FinishedAt.Sub(*tile.Build.StartedAt))
	}

	return tile, nil
}

func (gu *githubUsecase) PullRequests(params interface{}) ([]builder.Result, error) {
	prParams := params.(*models.PullRequestParams)

	pullRequests, err := gu.repository.GetPullRequests(prParams.Owner, prParams.Repository)
	if err != nil {
		return nil, &coreModels.MonitororError{Err: err, Message: "unable to find pull request"}
	}

	var results []builder.Result
	for _, pullRequest := range pullRequests {
		p := make(map[string]interface{})
		p["owner"] = pullRequest.Owner
		p["repository"] = pullRequest.Repository
		p["ref"] = pullRequest.Ref

		results = append(results, builder.Result{
			TileType: api.GithubChecksTileType,
			Label:    fmt.Sprintf("PR#%d @ %s", pullRequest.ID, pullRequest.Repository),
			Params:   p,
		})
	}

	return results, nil
}

func computeChecks(refStatus *models.Checks) (coreModels.TileStatus, *time.Time, *time.Time, string) {
	var statuses []coreModels.TileStatus
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
		return coreModels.UnknownStatus, nil, nil, ids
	}

	return statuses[0], startedAt, finishedAt, ids
}

func parseRun(run *models.Run) coreModels.TileStatus {
	// Based on : https://developer.github.com/v3/checks/runs/
	switch run.Status {
	case "in_progress":
		return coreModels.RunningStatus
	case "queued":
		return coreModels.QueuedStatus
	case "completed":
		switch run.Conclusion {
		case "success":
			return coreModels.SuccessStatus
		case "failure":
			return coreModels.FailedStatus
		case "timed_out":
			return coreModels.FailedStatus
		case "neutral":
			return coreModels.WarningStatus
		case "cancelled":
			return coreModels.CanceledStatus
		case "action_required":
			return coreModels.ActionRequiredStatus
		}
	}

	return coreModels.UnknownStatus
}

func parseStatus(status *models.Status) coreModels.TileStatus {
	// Based on : https://developer.github.com/v3/repos/statuses/
	switch status.State {
	case "success":
		return coreModels.SuccessStatus
	case "failure":
		return coreModels.FailedStatus
	case "error":
		return coreModels.FailedStatus
	case "pending":
		return coreModels.RunningStatus
	}

	return coreModels.UnknownStatus
}
