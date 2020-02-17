//+build !faker

package usecase

import (
	"net/url"
	"regexp"
	"time"

	"github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorable/jenkins"
	jenkinsModels "github.com/monitoror/monitoror/monitorable/jenkins/models"
	"github.com/monitoror/monitoror/pkg/monitoror/builder"
	"github.com/monitoror/monitoror/pkg/monitoror/cache"
	"github.com/monitoror/monitoror/pkg/monitoror/utils/git"

	"github.com/AlekSi/pointer"
)

type (
	jenkinsUsecase struct {
		repository jenkins.Repository

		// builds cache. used for save small history of build for stats
		buildsCache *cache.BuildCache
	}
)

const buildCacheSize = 5

func NewJenkinsUsecase(repository jenkins.Repository) jenkins.Usecase {
	return &jenkinsUsecase{
		repository,
		cache.NewBuildCache(buildCacheSize),
	}
}

func (tu *jenkinsUsecase) Build(params *jenkinsModels.BuildParams) (*models.Tile, error) {
	tile := models.NewTile(jenkins.JenkinsBuildTileType).WithBuild()

	tile.Label, _ = url.QueryUnescape(params.Job)
	if params.Branch != "" {
		branchLabel, _ := url.QueryUnescape(params.Branch)
		tile.Build.Branch = pointer.ToString(git.HumanizeBranch(branchLabel))
	}

	job, err := tu.repository.GetJob(params.Job, params.Branch)
	if err != nil {
		return nil, &models.MonitororError{Err: err, Tile: tile, Message: "unable to find job"}
	}

	// Is Buildable
	if !job.Buildable {
		tile.Status = models.DisabledStatus
		return tile, nil
	}

	// Set Previous Status
	previousStatus := tu.buildsCache.GetPreviousStatus(params, "null") // null because we don't have build number yet, but it's not important in jenkins
	if previousStatus != nil {
		tile.Build.PreviousStatus = *previousStatus
	} else {
		tile.Build.PreviousStatus = models.UnknownStatus
	}

	// Queued build
	if job.InQueue {
		tile.Status = models.QueuedStatus
		tile.Build.StartedAt = job.QueuedAt
		return tile, nil
	}

	// Get Last Build
	build, err := tu.repository.GetLastBuildStatus(job)
	if err != nil || build == nil {
		return nil, &models.MonitororError{Err: err, Tile: tile, Message: "no build found", ErrorStatus: models.UnknownStatus}
	}

	// Build ID
	tile.Build.ID = pointer.ToString(build.Number)

	// Set Status
	if build.Building {
		tile.Status = models.RunningStatus
	} else {
		tile.Status = parseResult(build.Result)
	}

	// Set StartedAt
	tile.Build.StartedAt = pointer.ToTime(build.StartedAt)

	// Set FinishedAt Or Duration
	if tile.Status != models.RunningStatus {
		tile.Build.FinishedAt = pointer.ToTime(build.StartedAt.Add(build.Duration))
	} else {
		tile.Build.Duration = pointer.ToInt64(int64(time.Since(build.StartedAt).Seconds()))

		estimatedDuration := tu.buildsCache.GetEstimatedDuration(params)
		if estimatedDuration != nil {
			tile.Build.EstimatedDuration = pointer.ToInt64(int64(estimatedDuration.Seconds()))
		} else {
			tile.Build.EstimatedDuration = pointer.ToInt64(int64(0))
		}
	}

	// Set Author
	if tile.Status == models.FailedStatus && build.Author != nil {
		tile.Build.Author = &models.Author{
			Name:      build.Author.Name,
			AvatarURL: build.Author.AvatarURL,
		}
	}

	// Cache Duration when success / failed / warning
	if tile.Status == models.SuccessStatus || tile.Status == models.FailedStatus || tile.Status == models.WarningStatus {
		tu.buildsCache.Add(params, *tile.Build.ID, tile.Status, build.Duration)
	}

	return tile, nil
}

func (tu *jenkinsUsecase) ListDynamicTile(params interface{}) ([]builder.Result, error) {
	mbParams := params.(*jenkinsModels.MultiBranchParams)

	job, err := tu.repository.GetJob(mbParams.Job, "")
	if err != nil {
		return nil, &models.MonitororError{Err: err, Message: "unable to find job"}
	}

	matcher, err := regexp.Compile(mbParams.Match)
	if err != nil {
		return nil, err
	}

	unmatcher, err := regexp.Compile(mbParams.Unmatch)
	if err != nil {
		return nil, err
	}

	var results []builder.Result
	for _, branch := range job.Branches {
		branchToFilter, _ := url.QueryUnescape(branch)
		if !matcher.MatchString(branchToFilter) ||
			(mbParams.Unmatch != "" && unmatcher.MatchString(branchToFilter)) {
			continue
		}

		p := make(map[string]interface{})
		p["job"] = mbParams.Job
		p["branch"] = branch

		results = append(results, builder.Result{
			TileType: jenkins.JenkinsBuildTileType,
			Params:   p,
		})
	}

	return results, nil
}

func parseResult(result string) models.TileStatus {
	switch result {
	case "SUCCESS":
		return models.SuccessStatus
	case "UNSTABLE":
		return models.WarningStatus
	case "FAILURE":
		return models.FailedStatus
	case "ABORTED":
		return models.CanceledStatus
	default:
		return models.UnknownStatus
	}
}
