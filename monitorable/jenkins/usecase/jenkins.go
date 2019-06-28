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

func (tu *jenkinsUsecase) Build(params *models.JobParams) (tile *BuildTile, err error) {
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
		if strings.Contains(err.Error(), "request canceled while waiting for connection") ||
			strings.Contains(err.Error(), "unsupported protocol scheme") {
			err = errors.NewTimeoutError(tile.Tile, "Timeout/Host Unreachable")
		} else {
			err = errors.NewSystemError("unable to get jenkins job", nil)
		}
		return nil, err
	}

	if !job.Buildable {
		tile.Status = DisabledStatus
		return
	}

	build, err := tu.repository.GetLastBuildStatus(job)
	if err != nil || build == nil {
		err = errors.NewNoBuildError(tile)
		return nil, err
	}

	// Set Status
	if job.InQueue {
		tile.Status = QueuedStatus
	} else if build.Building {
		tile.Status = RunningStatus
	} else {
		tile.Status = parseResult(build.Result)
	}

	// Set Previous Status
	if tile.Status == RunningStatus || tile.Status == QueuedStatus || tile.Status == AbortedStatus {
		previousStatus := tu.buildsCache.GetPreviousStatus(tile.Label)
		if previousStatus != nil {
			tile.PreviousStatus = *previousStatus
		} else {
			tile.PreviousStatus = UnknownStatus
		}
	}

	// Set StartedAt
	if !build.StartedAt.IsZero() {
		tile.StartedAt = ToInt64(build.StartedAt.Unix())
	}
	// Set FinishedAt
	if tile.Status != RunningStatus {
		tile.FinishedAt = ToInt64(build.StartedAt.Add(build.Duration).Unix())
	}

	// Set Duration / EstimatedDuration
	if tile.Status == RunningStatus {
		tile.Duration = ToInt64(int64(build.Duration.Seconds()))

		estimatedDuration := tu.buildsCache.GetEstimatedDuration(tile.Label)
		if estimatedDuration != nil {
			tile.EstimatedDuration = ToInt64(int64(*estimatedDuration / time.Second))
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
		tu.buildsCache.Add(tile.Label, tile.Status, build.Duration)
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
