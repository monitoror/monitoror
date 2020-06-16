//+build !faker

package usecase

import (
	"net/url"
	"regexp"
	"time"

	uiConfigModels "github.com/monitoror/monitoror/api/config/models"
	"github.com/monitoror/monitoror/internal/pkg/monitorable/cache"
	coreModels "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorables/jenkins/api"
	"github.com/monitoror/monitoror/monitorables/jenkins/api/models"
	"github.com/monitoror/monitoror/pkg/git"

	"github.com/AlekSi/pointer"
)

type (
	jenkinsUsecase struct {
		repository api.Repository

		// builds cache. used for save small history of build for stats
		buildsCache *cache.BuildCache
	}
)

const buildCacheSize = 5

func NewJenkinsUsecase(repository api.Repository) api.Usecase {
	return &jenkinsUsecase{
		repository,
		cache.NewBuildCache(buildCacheSize),
	}
}

func (tu *jenkinsUsecase) Build(params *models.BuildParams) (*coreModels.Tile, error) {
	tile := coreModels.NewTile(api.JenkinsBuildTileType).WithBuild()

	tile.Label, _ = url.QueryUnescape(params.Job)
	if params.Branch != "" {
		branchLabel, _ := url.QueryUnescape(params.Branch)
		tile.Build.Branch = pointer.ToString(git.HumanizeBranch(branchLabel))
	}

	job, err := tu.repository.GetJob(params.Job, params.Branch)
	if err != nil {
		return nil, &coreModels.MonitororError{Err: err, Tile: tile, Message: "unable to find job"}
	}

	// Is Buildable
	if !job.Buildable {
		tile.Status = coreModels.DisabledStatus
		return tile, nil
	}

	// Set Previous Status
	previousStatus := tu.buildsCache.GetPreviousStatus(params, "null") // null because we don't have build number yet, but it's not important in jenkins
	if previousStatus != nil {
		tile.Build.PreviousStatus = *previousStatus
	} else {
		tile.Build.PreviousStatus = coreModels.UnknownStatus
	}

	// Queued build
	if job.InQueue {
		tile.Status = coreModels.QueuedStatus
		tile.Build.StartedAt = job.QueuedAt
		return tile, nil
	}

	// Get Last Build
	build, err := tu.repository.GetLastBuildStatus(job)
	if err != nil || build == nil {
		return nil, &coreModels.MonitororError{Err: err, Tile: tile, Message: "no build found", ErrorStatus: coreModels.UnknownStatus}
	}

	// Build ID
	tile.Build.ID = pointer.ToString(build.Number)

	// Set Status
	if build.Building {
		tile.Status = coreModels.RunningStatus
	} else {
		tile.Status = parseResult(build.Result)
	}

	// Set StartedAt
	tile.Build.StartedAt = pointer.ToTime(build.StartedAt)

	// Set FinishedAt Or Duration
	if tile.Status != coreModels.RunningStatus {
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
	if tile.Status == coreModels.FailedStatus && build.Author != nil {
		tile.Build.Author = &coreModels.Author{
			Name:      build.Author.Name,
			AvatarURL: build.Author.AvatarURL,
		}
	}

	// Cache Duration when success / failed / warning
	if tile.Status == coreModels.SuccessStatus || tile.Status == coreModels.FailedStatus || tile.Status == coreModels.WarningStatus {
		tu.buildsCache.Add(params, *tile.Build.ID, tile.Status, build.Duration)
	}

	return tile, nil
}

func (tu *jenkinsUsecase) BuildGenerator(params interface{}) ([]uiConfigModels.GeneratedTile, error) {
	mbParams := params.(*models.BuildGeneratorParams)

	job, err := tu.repository.GetJob(mbParams.Job, "")
	if err != nil {
		return nil, &coreModels.MonitororError{Err: err, Message: "unable to find job"}
	}

	matcher, err := regexp.Compile(mbParams.Match)
	if err != nil {
		return nil, err
	}

	unmatcher, err := regexp.Compile(mbParams.Unmatch)
	if err != nil {
		return nil, err
	}

	var results []uiConfigModels.GeneratedTile
	for _, branch := range job.Branches {
		branchToFilter, _ := url.QueryUnescape(branch)
		if !matcher.MatchString(branchToFilter) ||
			(mbParams.Unmatch != "" && unmatcher.MatchString(branchToFilter)) {
			continue
		}

		p := &models.BuildParams{}
		p.Job = mbParams.Job
		p.Branch = branch

		results = append(results, uiConfigModels.GeneratedTile{
			Params: p,
		})
	}

	return results, nil
}

func parseResult(result string) coreModels.TileStatus {
	switch result {
	case "SUCCESS":
		return coreModels.SuccessStatus
	case "UNSTABLE":
		return coreModels.WarningStatus
	case "FAILURE":
		return coreModels.FailedStatus
	case "ABORTED":
		return coreModels.CanceledStatus
	default:
		return coreModels.UnknownStatus
	}
}
