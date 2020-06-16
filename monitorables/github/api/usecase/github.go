//+build !faker

package usecase

import (
	"fmt"
	"sort"
	"time"

	uiConfigModels "github.com/monitoror/monitoror/api/config/models"
	"github.com/monitoror/monitoror/internal/pkg/monitorable/cache"
	coreModels "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorables/github/api"
	"github.com/monitoror/monitoror/monitorables/github/api/models"
	"github.com/monitoror/monitoror/pkg/git"
	"github.com/monitoror/monitoror/pkg/hash"

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
	tile.Label = "GitHub count"

	count, err := gu.repository.GetCount(params.Query)
	if err != nil {
		return nil, &coreModels.MonitororError{Err: err, Tile: tile, Message: "unable to load count or wrong query"}
	}

	tile.Status = coreModels.SuccessStatus
	tile.Value.Values = append(tile.Value.Values, fmt.Sprintf("%d", count))

	return tile, nil
}

func (gu *githubUsecase) Checks(params *models.ChecksParams) (*coreModels.Tile, error) {
	tile := coreModels.NewTile(api.GithubChecksTileType).WithBuild()
	tile.Label = params.Repository
	tile.Build.Branch = pointer.ToString(git.HumanizeBranch(params.Ref))

	// Request checks
	checks, err := gu.repository.GetChecks(params.Owner, params.Repository, params.Ref)
	if err != nil {
		return nil, &coreModels.MonitororError{Err: err, Tile: tile, Message: "unable to load ref checks"}
	}
	if len(checks.Statuses) == 0 && len(checks.Runs) == 0 {
		// Warning because request was correct but there is no build
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

func (gu *githubUsecase) PullRequest(params *models.PullRequestParams) (*coreModels.Tile, error) {
	tile := coreModels.NewTile(api.GithubPullRequestTileType).WithBuild()

	// Request pullRequest
	pullRequest, err := gu.repository.GetPullRequest(params.Owner, params.Repository, *params.ID)
	if err != nil {
		return nil, &coreModels.MonitororError{Err: err, Tile: tile, Message: "unable to load pull request"}
	}

	tile.Label = params.Repository
	tile.Build.Branch = pointer.ToString(git.HumanizeBranch(pullRequest.SourceBranch))
	if params.Owner != pullRequest.SourceOwner {
		tile.Build.Branch = pointer.ToString(fmt.Sprintf("%s:%s", pullRequest.SourceOwner, *tile.Build.Branch))
	}
	tile.Build.MergeRequest = &coreModels.TileMergeRequest{
		ID:    pullRequest.ID,
		Title: pullRequest.Title,
	}

	// With pull request, we use CommitSHA as Ref on params.Owner/params.Repo to handle correctly Checks from Forks
	// Request checks
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

func (gu *githubUsecase) PullRequestsGenerator(params interface{}) ([]uiConfigModels.GeneratedTile, error) {
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

func (gu *githubUsecase) computeRefChecks(tile *coreModels.Tile, checks *models.Checks, paramsKey string) {
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

	for _, run := range checks.Runs {
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
	sort.Slice(checks.Statuses, func(i, j int) bool {
		return checks.Statuses[i].CreatedAt.After(checks.Statuses[j].CreatedAt)
	})

	titles := make(map[string]bool)
	for _, status := range checks.Statuses {
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

	return statuses, startedAt, finishedAt, hash.GetMD5Hash(ids)
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
