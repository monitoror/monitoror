//+build !faker

package usecase

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/monitoror/monitoror/config"

	. "github.com/AlekSi/pointer"
	gocache "github.com/robfig/go-cache"

	"github.com/monitoror/monitoror/models/errors"
	. "github.com/monitoror/monitoror/models/tiles"
	"github.com/monitoror/monitoror/monitorable/jenkins"
	"github.com/monitoror/monitoror/monitorable/jenkins/models"
	"github.com/monitoror/monitoror/pkg/monitoror/builder"
	"github.com/monitoror/monitoror/pkg/monitoror/cache"
)

type (
	jenkinsUsecase struct {
		repository jenkins.Repository

		// builds cache
		buildsCache *cache.BuildCache

		// jobs cache
		jobsCache *gocache.Cache
	}
)

const buildCacheSize = 5

func NewJenkinsUsecase(repository jenkins.Repository, downstreamCache config.Cache) jenkins.Usecase {
	return &jenkinsUsecase{
		repository,
		cache.NewBuildCache(buildCacheSize),
		gocache.New(
			time.Millisecond*time.Duration(downstreamCache.Expire),
			time.Millisecond*time.Duration(downstreamCache.CleanupInterval),
		),
	}
}

func (tu *jenkinsUsecase) Build(params *models.BuildParams) (tile *BuildTile, err error) {
	tile = NewBuildTile(jenkins.JenkinsBuildTileType)

	jobLabel, _ := url.QueryUnescape(params.Job)
	branchLabel, _ := url.QueryUnescape(params.Branch)
	if params.Branch == "" {
		tile.Label = jobLabel
	} else {
		tile.Label = fmt.Sprintf("%s : #%s", jobLabel, branchLabel)
	}

	job, err := tu.repository.GetJob(params.Job, params.Branch)
	if err != nil {
		// TODO : Replace that by errors.Is/As when go 1.13 will be released
		if strings.Contains(err.Error(), "no such host") ||
			strings.Contains(err.Error(), "dial tcp: lookup") ||
			strings.Contains(err.Error(), "request canceled") {
			err = errors.NewTimeoutError(nil)
		} else {
			err = errors.NewSystemError("unable to found job", err)
		}
		return nil, err
	}

	// Is Buildable
	if !job.Buildable {
		tile.Status = DisabledStatus
		return
	}

	// Set Previous Status
	previousStatus := tu.buildsCache.GetPreviousStatus(tile.Label, "null")
	if previousStatus != nil {
		tile.PreviousStatus = *previousStatus
	} else {
		tile.PreviousStatus = UnknownStatus
	}

	// Queued build
	if job.InQueue {
		tile.Status = QueuedStatus
		tile.StartedAt = job.QueuedAt
		return
	}

	// Get Last Build
	build, err := tu.repository.GetLastBuildStatus(job)
	if err != nil || build == nil {
		// TODO : Replace that by errors.Is/As when go 1.13 will be released
		if err != nil && (strings.Contains(err.Error(), "no such host") ||
			strings.Contains(err.Error(), "dial tcp: lookup") ||
			strings.Contains(err.Error(), "request canceled")) {
			err = errors.NewTimeoutError(tile.Tile)
		} else {
			err = errors.NewNoBuildError(tile)
		}
		return nil, err
	}

	// Set Status
	if build.Building {
		tile.Status = RunningStatus
	} else {
		tile.Status = parseResult(build.Result)
	}

	// Set StartedAt
	tile.StartedAt = ToTime(build.StartedAt)

	// Set FinishedAt Or Duration
	if tile.Status != RunningStatus {
		tile.FinishedAt = ToTime(build.StartedAt.Add(build.Duration))
	} else {
		tile.Duration = ToInt64(int64(time.Now().Sub(build.StartedAt).Seconds()))

		estimatedDuration := tu.buildsCache.GetEstimatedDuration(tile.Label)
		if estimatedDuration != nil {
			tile.EstimatedDuration = ToInt64(int64(estimatedDuration.Seconds()))
		} else {
			tile.EstimatedDuration = ToInt64(int64(0))
		}
	}

	// Set Author
	if build.Author != nil {
		tile.Author = &Author{
			Name:      build.Author.Name,
			AvatarUrl: build.Author.AvatarUrl,
		}
	}

	// Cache Duration when success / failed / warning
	if tile.Status == SuccessStatus || tile.Status == FailedStatus || tile.Status == WarningStatus {
		tu.buildsCache.Add(tile.Label, build.Number, tile.Status, build.Duration)
	}

	return
}

func (tu *jenkinsUsecase) ListDynamicTile(params interface{}) (results []builder.Result, err error) {
	mbParams := params.(*models.MultiBranchParams)

	job, err := tu.repository.GetJob(mbParams.Job, "")
	if err != nil {
		// TODO : Replace that by errors.Is/As when go 1.13 will be released
		if strings.Contains(err.Error(), "no such host") ||
			strings.Contains(err.Error(), "dial tcp: lookup") ||
			strings.Contains(err.Error(), "request canceled") {

			// Get previous value in cache
			j, exist := tu.jobsCache.Get(mbParams.Job)
			if !exist {
				err = errors.NewTimeoutError(nil)
				return
			}
			job = j.(*models.Job)
		} else {
			err = errors.NewSystemError("unable to found job", err)
			return
		}
	} else {
		tu.jobsCache.Set(mbParams.Job, job, 0)
	}

	matcher, err := regexp.Compile(mbParams.Match)
	if err != nil {
		return
	}

	unmatcher, err := regexp.Compile(mbParams.Unmatch)
	if err != nil {
		return
	}

	results = []builder.Result{}
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

	return
}

func parseResult(result string) TileStatus {
	switch result {
	case "SUCCESS":
		return SuccessStatus
	case "UNSTABLE":
		return WarningStatus
	case "FAILURE":
		return FailedStatus
	case "ABORTED":
		return AbortedStatus
	default:
		return UnknownStatus
	}
}
