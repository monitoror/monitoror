//+build !faker

package usecase

import (
	"fmt"
	"time"

	"github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorable/travisci"
	travisModels "github.com/monitoror/monitoror/monitorable/travisci/models"
	"github.com/monitoror/monitoror/pkg/monitoror/cache"
	"github.com/monitoror/monitoror/pkg/monitoror/utils/git"

	"github.com/AlekSi/pointer"
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

func (tu *travisCIUsecase) Build(params *travisModels.BuildParams) (*models.Tile, error) {
	tile := models.NewTile(travisci.TravisCIBuildTileType).WithBuild()
	tile.Label = params.Repository
	tile.Build.Branch = pointer.ToString(git.HumanizeBranch(params.Branch))

	// Request
	build, err := tu.repository.GetLastBuildStatus(params.Owner, params.Repository, params.Branch)
	if err != nil {
		return nil, &models.MonitororError{Err: err, Tile: tile, Message: "unable to find build"}
	}
	if build == nil {
		// Warning because request was correct but there is no build
		return nil, &models.MonitororError{Tile: tile, Message: "no build found", ErrorStatus: models.UnknownStatus}
	}

	tile.Build.ID = pointer.ToString(fmt.Sprintf("%d", build.ID))

	// Set Status
	tile.Status = parseState(build.State)

	// Set Previous Status
	previousStatus := tu.buildsCache.GetPreviousStatus(params, *tile.Build.ID)
	if previousStatus != nil {
		tile.Build.PreviousStatus = *previousStatus
	} else {
		tile.Build.PreviousStatus = models.UnknownStatus
	}

	// Set StartedAt
	if !build.StartedAt.IsZero() {
		tile.Build.StartedAt = pointer.ToTime(build.StartedAt)
	}
	// Set FinishedAt
	if !build.FinishedAt.IsZero() {
		tile.Build.FinishedAt = pointer.ToTime(build.FinishedAt)
	}

	if tile.Status == models.RunningStatus {
		tile.Build.Duration = pointer.ToInt64(int64(time.Since(build.StartedAt).Seconds()))

		estimatedDuration := tu.buildsCache.GetEstimatedDuration(params)
		if estimatedDuration != nil {
			tile.Build.EstimatedDuration = pointer.ToInt64(int64(estimatedDuration.Seconds()))
		} else {
			tile.Build.EstimatedDuration = pointer.ToInt64(int64(0))
		}
	}

	// Set Author
	if tile.Status == models.FailedStatus && (build.Author.Name != "" || build.Author.AvatarURL != "") {
		tile.Build.Author = &models.Author{
			Name:      build.Author.Name,
			AvatarURL: build.Author.AvatarURL,
		}
	}

	// Cache Duration when success / failed
	if tile.Status == models.SuccessStatus || tile.Status == models.FailedStatus {
		tu.buildsCache.Add(params, *tile.Build.ID, tile.Status, build.Duration)
	}

	return tile, nil
}

func parseState(state string) models.TileStatus {
	switch state {
	case "created":
		return models.QueuedStatus
	case "received":
		return models.QueuedStatus
	case "started":
		return models.RunningStatus
	case "passed":
		return models.SuccessStatus
	case "failed":
		return models.FailedStatus
	case "errored":
		return models.FailedStatus
	case "canceled":
		return models.CanceledStatus
	default:
		return models.UnknownStatus
	}
}
