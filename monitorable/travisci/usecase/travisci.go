//+build !faker

package usecase

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/monitoror/monitoror/pkg/monitoror/cache"

	. "github.com/AlekSi/pointer"

	"github.com/monitoror/monitoror/models/errors"

	. "github.com/monitoror/monitoror/models/tiles"
	"github.com/monitoror/monitoror/monitorable/travisci"
	"github.com/monitoror/monitoror/monitorable/travisci/models"
)

type (
	travisCIUsecase struct {
		repository travisci.Repository

		// builds cache
		buildsCache *cache.BuildCache
	}
)

const cacheSize = 5

func NewTravisCIUsecase(repository travisci.Repository) travisci.Usecase {
	return &travisCIUsecase{repository, cache.NewBuildCache(cacheSize)}
}

func (tu *travisCIUsecase) Build(params *models.BuildParams) (tile *BuildTile, err error) {
	tile = NewBuildTile(travisci.TravisCIBuildTileType)
	tile.Label = fmt.Sprintf("%s : #%s", params.Repository, params.Branch)

	// Request
	build, err := tu.repository.GetLastBuildStatus(params.Group, params.Repository, params.Branch)
	if err != nil {
		// TODO : Replace that by errors.Is/As when go 1.13 will be released
		if err == context.DeadlineExceeded ||
			strings.Contains(err.Error(), "no such host") ||
			strings.Contains(err.Error(), "dial tcp: lookup") {
			err = errors.NewTimeoutError(tile.Tile)
		} else {
			err = errors.NewSystemError("unable to get travisci build", nil)
		}
		return nil, err
	}
	if build == nil {
		err = errors.NewNoBuildError(tile)
		return nil, err
	}

	// Set Status
	tile.Status = parseState(build.State)

	// Set Previous Status
	previousStatus := tu.buildsCache.GetPreviousStatus(tile.Label, fmt.Sprintf("%d", build.Id))
	if previousStatus != nil {
		tile.PreviousStatus = *previousStatus
	} else {
		tile.PreviousStatus = UnknownStatus
	}

	// Set StartedAt
	if !build.StartedAt.IsZero() {
		tile.StartedAt = ToTime(build.StartedAt)
	}
	// Set FinishedAt
	if !build.FinishedAt.IsZero() {
		tile.FinishedAt = ToTime(build.FinishedAt)
	}

	if tile.Status == RunningStatus {
		tile.Duration = ToInt64(int64(time.Now().Sub(build.StartedAt).Seconds()))

		estimatedDuration := tu.buildsCache.GetEstimatedDuration(tile.Label)
		if estimatedDuration != nil {
			tile.EstimatedDuration = ToInt64(int64(estimatedDuration.Seconds()))
		} else {
			tile.EstimatedDuration = ToInt64(int64(0))
		}
	}

	// Set Author
	if build.Author.Name != "" || build.Author.AvatarUrl != "" {
		tile.Author = &Author{
			Name:      build.Author.Name,
			AvatarUrl: build.Author.AvatarUrl,
		}
	}

	// Cache Duration when success / failed
	if tile.Status == SuccessStatus || tile.Status == FailedStatus {
		tu.buildsCache.Add(tile.Label, fmt.Sprintf("%d", build.Id), tile.Status, build.Duration)
	}

	return
}

func parseState(state string) TileStatus {
	switch state {
	case "created":
		return QueuedStatus
	case "received":
		return QueuedStatus
	case "started":
		return RunningStatus
	case "passed":
		return SuccessStatus
	case "failed":
		return FailedStatus
	case "errored":
		return FailedStatus
	case "canceled":
		return AbortedStatus
	default:
		return UnknownStatus
	}
}
