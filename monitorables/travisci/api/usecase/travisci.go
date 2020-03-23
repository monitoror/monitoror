//+build !faker

package usecase

import (
	"fmt"
	"time"

	"github.com/monitoror/monitoror/internal/pkg/monitorable/cache"
	coreModels "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorables/travisci/api"
	"github.com/monitoror/monitoror/monitorables/travisci/api/models"
	"github.com/monitoror/monitoror/pkg/git"

	"github.com/AlekSi/pointer"
)

type (
	travisCIUsecase struct {
		repository api.Repository

		// builds cache
		buildsCache *cache.BuildCache
	}
)

const cacheSize = 5

func NewTravisCIUsecase(repository api.Repository) api.Usecase {
	return &travisCIUsecase{repository, cache.NewBuildCache(cacheSize)}
}

func (tu *travisCIUsecase) Build(params *models.BuildParams) (*coreModels.Tile, error) {
	tile := coreModels.NewTile(api.TravisCIBuildTileType).WithBuild()
	tile.Label = params.Repository
	tile.Build.Branch = pointer.ToString(git.HumanizeBranch(params.Branch))

	// Request
	build, err := tu.repository.GetLastBuildStatus(params.Owner, params.Repository, params.Branch)
	if err != nil {
		return nil, &coreModels.MonitororError{Err: err, Tile: tile, Message: "unable to find build"}
	}
	if build == nil {
		// Warning because request was correct but there is no build
		return nil, &coreModels.MonitororError{Tile: tile, Message: "no build found", ErrorStatus: coreModels.UnknownStatus}
	}

	tile.Build.ID = pointer.ToString(fmt.Sprintf("%d", build.ID))

	// Set Status
	tile.Status = parseState(build.State)

	// Set Previous Status
	previousStatus := tu.buildsCache.GetPreviousStatus(params, *tile.Build.ID)
	if previousStatus != nil {
		tile.Build.PreviousStatus = *previousStatus
	} else {
		tile.Build.PreviousStatus = coreModels.UnknownStatus
	}

	// Set StartedAt
	if !build.StartedAt.IsZero() {
		tile.Build.StartedAt = pointer.ToTime(build.StartedAt)
	}
	// Set FinishedAt
	if !build.FinishedAt.IsZero() {
		tile.Build.FinishedAt = pointer.ToTime(build.FinishedAt)
	}

	if tile.Status == coreModels.RunningStatus {
		tile.Build.Duration = pointer.ToInt64(int64(time.Since(build.StartedAt).Seconds()))

		estimatedDuration := tu.buildsCache.GetEstimatedDuration(params)
		if estimatedDuration != nil {
			tile.Build.EstimatedDuration = pointer.ToInt64(int64(estimatedDuration.Seconds()))
		} else {
			tile.Build.EstimatedDuration = pointer.ToInt64(int64(0))
		}
	}

	// Set Author
	if tile.Status == coreModels.FailedStatus && (build.Author.Name != "" || build.Author.AvatarURL != "") {
		tile.Build.Author = &coreModels.Author{
			Name:      build.Author.Name,
			AvatarURL: build.Author.AvatarURL,
		}
	}

	// Cache Duration when success / failed
	if tile.Status == coreModels.SuccessStatus || tile.Status == coreModels.FailedStatus {
		tu.buildsCache.Add(params, *tile.Build.ID, tile.Status, build.Duration)
	}

	return tile, nil
}

func parseState(state string) coreModels.TileStatus {
	switch state {
	case "created":
		return coreModels.QueuedStatus
	case "received":
		return coreModels.QueuedStatus
	case "started":
		return coreModels.RunningStatus
	case "passed":
		return coreModels.SuccessStatus
	case "failed":
		return coreModels.FailedStatus
	case "errored":
		return coreModels.FailedStatus
	case "canceled":
		return coreModels.CanceledStatus
	default:
		return coreModels.UnknownStatus
	}
}
