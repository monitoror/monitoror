//+build !faker

package usecase

import (
	"fmt"
	"time"

	"github.com/monitoror/monitoror/pkg/monitoror/utils/git"

	. "github.com/AlekSi/pointer"

	. "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorable/azuredevops"
	"github.com/monitoror/monitoror/monitorable/azuredevops/models"
	"github.com/monitoror/monitoror/pkg/monitoror/cache"
)

type (
	azureDevOpsUsecase struct {
		repository azuredevops.Repository

		// builds cache. used for save small history of build for stats
		buildsCache *cache.BuildCache
	}
)

const buildCacheSize = 5

func NewAzureDevOpsUsecase(repository azuredevops.Repository) azuredevops.Usecase {
	return &azureDevOpsUsecase{
		repository,
		cache.NewBuildCache(buildCacheSize),
	}
}

func (au *azureDevOpsUsecase) Build(params *models.BuildParams) (tile *Tile, err error) {
	tile = NewTile(azuredevops.AzureDevOpsBuildTileType)
	// Default label if build not found
	tile.Label = fmt.Sprintf("%s", params.Project)

	// Lookup for build
	build, err := au.repository.GetBuild(params.Project, *params.Definition, params.Branch)
	if err != nil {
		return nil, &MonitororError{Err: err, Tile: tile, Message: "unable to find build"}
	}
	if build == nil {
		return nil, &MonitororError{Tile: tile, Message: "no build found", ErrorStatus: UnknownStatus}
	}

	// Title and Message
	tile.Label = fmt.Sprintf("%s | %s", params.Project, build.DefinitionName)
	tile.Message = fmt.Sprintf("%s - %s", git.HumanizeBranch(build.Branch), build.BuildNumber)

	// Status
	tile.Status = parseBuildResult(build.Status, build.Result)

	// Previous Status
	previousStatus := au.buildsCache.GetPreviousStatus(params, build.BuildNumber)
	if previousStatus != nil {
		tile.PreviousStatus = *previousStatus
	} else {
		tile.PreviousStatus = UnknownStatus
	}

	// Author
	if build.Author != nil {
		tile.Author = &Author{
			Name:      build.Author.Name,
			AvatarUrl: build.Author.AvatarUrl,
		}
	}

	// StartedAt / FinishedAt
	tile.StartedAt = build.StartedAt
	if tile.Status != QueuedStatus && tile.Status != RunningStatus {
		tile.FinishedAt = build.FinishedAt
	}

	if tile.Status == QueuedStatus {
		tile.StartedAt = build.QueuedAt
	}

	// Duration / Previous Duration
	if tile.Status == RunningStatus {
		tile.Duration = ToInt64(int64(time.Now().Sub(*tile.StartedAt).Seconds()))

		estimatedDuration := au.buildsCache.GetEstimatedDuration(params)
		if estimatedDuration != nil {
			tile.EstimatedDuration = ToInt64(int64(estimatedDuration.Seconds()))
		} else {
			tile.EstimatedDuration = ToInt64(int64(0))
		}
	}

	// Cache Duration when success / failed
	if tile.Status == SuccessStatus || tile.Status == FailedStatus {
		au.buildsCache.Add(params, build.BuildNumber, tile.Status, tile.FinishedAt.Sub(*tile.StartedAt))
	}

	return
}

func (au *azureDevOpsUsecase) Release(params *models.ReleaseParams) (tile *Tile, err error) {
	tile = NewTile(azuredevops.AzureDevOpsReleaseTileType)
	// Default label if build not found
	tile.Label = fmt.Sprintf("%s", params.Project)

	// Lookup for release
	release, err := au.repository.GetRelease(params.Project, *params.Definition)
	if err != nil {
		return nil, &MonitororError{Err: err, Tile: tile, Message: "unable to find release"}
	}
	if release == nil {
		// Warning because request was correct but there is no release
		return nil, &MonitororError{Tile: tile, Message: "no release found", ErrorStatus: UnknownStatus}
	}

	// Label
	tile.Label = fmt.Sprintf("%s | %s", params.Project, release.DefinitionName)
	tile.Message = release.ReleaseNumber

	// Status
	tile.Status = parseReleaseStatus(release.Status)

	// Previous Status
	previousStatus := au.buildsCache.GetPreviousStatus(params, release.ReleaseNumber)
	if previousStatus != nil {
		tile.PreviousStatus = *previousStatus
	} else {
		tile.PreviousStatus = UnknownStatus
	}

	// Author
	if release.Author != nil {
		tile.Author = &Author{
			Name:      release.Author.Name,
			AvatarUrl: release.Author.AvatarUrl,
		}
	}

	// StartedAt / FinishedAt
	tile.StartedAt = release.StartedAt
	tile.StartedAt = release.StartedAt
	if tile.Status != RunningStatus {
		tile.FinishedAt = release.FinishedAt
	}
	// Duration
	if tile.Status == RunningStatus {
		tile.Duration = ToInt64(int64(time.Now().Sub(*tile.StartedAt).Seconds()))

		estimatedDuration := au.buildsCache.GetEstimatedDuration(params)
		if estimatedDuration != nil {
			tile.EstimatedDuration = ToInt64(int64(estimatedDuration.Seconds()))
		} else {
			tile.EstimatedDuration = ToInt64(int64(0))
		}
	}

	// Cache Duration when success / failed
	if tile.Status == SuccessStatus || tile.Status == FailedStatus {
		au.buildsCache.Add(params, release.ReleaseNumber, tile.Status, tile.FinishedAt.Sub(*tile.StartedAt))
	}

	return
}

func parseBuildResult(status, result string) TileStatus {
	switch status {
	case "inProgress":
		return RunningStatus
	case "cancelling":
		return RunningStatus
	case "notStarted":
		return QueuedStatus
	case "completed":
		switch result {
		case "succeeded":
			return SuccessStatus
		case "partiallySucceeded":
			return WarningStatus
		case "failed":
			return FailedStatus
		case "canceled":
			return AbortedStatus
		default:
			return UnknownStatus
		}
	default:
		return UnknownStatus
	}
}

func parseReleaseStatus(status string) TileStatus {
	switch status {
	case "failed":
		return FailedStatus
	case "succeeded":
		return SuccessStatus
	case "partiallySucceeded":
		return WarningStatus
	case "inProgress":
		return RunningStatus
	default: // all / notDeployed
		return UnknownStatus
	}
}
