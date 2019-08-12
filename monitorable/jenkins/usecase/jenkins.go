//+build !faker

package usecase

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	. "github.com/AlekSi/pointer"

	"github.com/monitoror/monitoror/pkg/monitoror/cache"

	"github.com/monitoror/monitoror/models/errors"

	. "github.com/monitoror/monitoror/models/tiles"
	"github.com/monitoror/monitoror/monitorable/jenkins"
	"github.com/monitoror/monitoror/monitorable/jenkins/models"
)

type (
	jenkinsUsecase struct {
		repository jenkins.Repository

		// builds cache
		buildsCache *cache.BuildCache
	}
)

const cacheSize = 5

func NewJenkinsUsecase(repository jenkins.Repository) jenkins.Usecase {
	return &jenkinsUsecase{repository, cache.NewBuildCache(cacheSize)}
}

func (tu *jenkinsUsecase) Build(params *models.BuildParams) (tile *BuildTile, err error) {
	tile = NewBuildTile(jenkins.JenkinsBuildTileType)

	jobLabel, _ := url.QueryUnescape(params.Job)
	parentLabel, _ := url.QueryUnescape(params.Parent)
	if params.Parent == "" {
		tile.Label = jobLabel
	} else {
		tile.Label = fmt.Sprintf("%s : #%s", parentLabel, jobLabel)
	}

	job, err := tu.repository.GetJob(params.Job, params.Parent)
	if err != nil {
		// TODO : Replace that by errors.Is when go 1.13 will be released
		if strings.Contains(err.Error(), "no such host") ||
			strings.Contains(err.Error(), "dial tcp: lookup") ||
			strings.Contains(err.Error(), "request canceled") ||
			strings.Contains(err.Error(), "unsupported protocol scheme") {
			err = errors.NewTimeoutError(tile.Tile)
		} else {
			err = errors.NewSystemError("unable to found job", nil)
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
		if err != nil && (strings.Contains(err.Error(), "no such host") ||
			strings.Contains(err.Error(), "dial tcp: lookup") ||
			strings.Contains(err.Error(), "request canceled") ||
			strings.Contains(err.Error(), "unsupported protocol scheme")) {
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
