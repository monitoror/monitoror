//+build !faker

package usecase

import (
	"fmt"
	"time"

	coreModels "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorables/azuredevops/api"
	"github.com/monitoror/monitoror/monitorables/azuredevops/api/models"
	"github.com/monitoror/monitoror/pkg/monitoror/cache"
	"github.com/monitoror/monitoror/pkg/monitoror/utils/git"

	"github.com/AlekSi/pointer"
)

type (
	azureDevOpsUsecase struct {
		repository api.Repository

		// builds cache. used for save small history of build for stats
		buildsCache *cache.BuildCache
	}
)

const buildCacheSize = 5

func NewAzureDevOpsUsecase(repository api.Repository) api.Usecase {
	return &azureDevOpsUsecase{
		repository,
		cache.NewBuildCache(buildCacheSize),
	}
}

func (au *azureDevOpsUsecase) Build(params *models.BuildParams) (*coreModels.Tile, error) {
	tile := coreModels.NewTile(api.AzureDevOpsBuildTileType).WithBuild()
	// Default label if build not found
	tile.Label = params.Project

	// Lookup for build
	build, err := au.repository.GetBuild(params.Project, *params.Definition, params.Branch)
	if err != nil {
		return nil, &coreModels.MonitororError{Err: err, Tile: tile, Message: "unable to find build"}
	}
	if build == nil {
		return nil, &coreModels.MonitororError{Tile: tile, Message: "no build found", ErrorStatus: coreModels.UnknownStatus}
	}

	// Title and Message
	tile.Label = fmt.Sprintf("%s (%s)", params.Project, build.DefinitionName)
	tile.Build.ID = &build.BuildNumber
	tile.Build.Branch = pointer.ToString(git.HumanizeBranch(build.Branch))

	// Status
	tile.Status = parseBuildResult(build.Status, build.Result)

	// Previous Status
	previousStatus := au.buildsCache.GetPreviousStatus(params, *tile.Build.ID)
	if previousStatus != nil {
		tile.Build.PreviousStatus = *previousStatus
	} else {
		tile.Build.PreviousStatus = coreModels.UnknownStatus
	}

	// Author
	if tile.Status == coreModels.FailedStatus && build.Author != nil {
		tile.Build.Author = &coreModels.Author{
			Name:      build.Author.Name,
			AvatarURL: build.Author.AvatarURL,
		}
	}

	// StartedAt / FinishedAt
	tile.Build.StartedAt = build.StartedAt
	if tile.Status != coreModels.QueuedStatus && tile.Status != coreModels.RunningStatus {
		tile.Build.FinishedAt = build.FinishedAt
	}

	if tile.Status == coreModels.QueuedStatus {
		tile.Build.StartedAt = build.QueuedAt
	}

	// Duration / Previous Duration
	if tile.Status == coreModels.RunningStatus {
		tile.Build.Duration = pointer.ToInt64(int64(time.Since(*tile.Build.StartedAt).Seconds()))

		estimatedDuration := au.buildsCache.GetEstimatedDuration(params)
		if estimatedDuration != nil {
			tile.Build.EstimatedDuration = pointer.ToInt64(int64(estimatedDuration.Seconds()))
		} else {
			tile.Build.EstimatedDuration = pointer.ToInt64(int64(0))
		}
	}

	// Cache Duration when success / failed
	if tile.Status == coreModels.SuccessStatus || tile.Status == coreModels.FailedStatus {
		au.buildsCache.Add(params, build.BuildNumber, tile.Status, tile.Build.FinishedAt.Sub(*tile.Build.StartedAt))
	}

	return tile, nil
}

func (au *azureDevOpsUsecase) Release(params *models.ReleaseParams) (*coreModels.Tile, error) {
	tile := coreModels.NewTile(api.AzureDevOpsReleaseTileType).WithBuild()
	// Default label if build not found
	tile.Label = params.Project

	// Lookup for release
	release, err := au.repository.GetRelease(params.Project, *params.Definition)
	if err != nil {
		return nil, &coreModels.MonitororError{Err: err, Tile: tile, Message: "unable to find release"}
	}
	if release == nil {
		// Warning because request was correct but there is no release
		return nil, &coreModels.MonitororError{Tile: tile, Message: "no release found", ErrorStatus: coreModels.UnknownStatus}
	}

	// Label
	tile.Label = fmt.Sprintf("%s (%s)", params.Project, release.DefinitionName)
	tile.Build.ID = &release.ReleaseNumber

	// Status
	tile.Status = parseReleaseStatus(release.Status)

	// Previous Status
	previousStatus := au.buildsCache.GetPreviousStatus(params, *tile.Build.ID)
	if previousStatus != nil {
		tile.Build.PreviousStatus = *previousStatus
	} else {
		tile.Build.PreviousStatus = coreModels.UnknownStatus
	}

	// Author
	if tile.Status == coreModels.FailedStatus && release.Author != nil {
		tile.Build.Author = &coreModels.Author{
			Name:      release.Author.Name,
			AvatarURL: release.Author.AvatarURL,
		}
	}

	// StartedAt / FinishedAt
	tile.Build.StartedAt = release.StartedAt
	if tile.Status != coreModels.RunningStatus && tile.Status != coreModels.QueuedStatus {
		tile.Build.FinishedAt = release.FinishedAt
	}
	// Duration
	if tile.Status == coreModels.RunningStatus {
		tile.Build.Duration = pointer.ToInt64(int64(time.Since(*tile.Build.StartedAt).Seconds()))

		estimatedDuration := au.buildsCache.GetEstimatedDuration(params)
		if estimatedDuration != nil {
			tile.Build.EstimatedDuration = pointer.ToInt64(int64(estimatedDuration.Seconds()))
		} else {
			tile.Build.EstimatedDuration = pointer.ToInt64(int64(0))
		}
	}

	// Cache Duration when success / failed
	if tile.Status == coreModels.SuccessStatus || tile.Status == coreModels.FailedStatus || tile.Status == coreModels.WarningStatus {
		au.buildsCache.Add(params, *tile.Build.ID, tile.Status, tile.Build.FinishedAt.Sub(*tile.Build.StartedAt))
	}

	return tile, nil
}

func parseBuildResult(status, result string) coreModels.TileStatus {
	switch status {
	case "inProgress":
		return coreModels.RunningStatus
	case "cancelling":
		return coreModels.RunningStatus
	case "notStarted":
		return coreModels.QueuedStatus
	case "completed":
		switch result {
		case "succeeded":
			return coreModels.SuccessStatus
		case "partiallySucceeded":
			return coreModels.WarningStatus
		case "failed":
			return coreModels.FailedStatus
		case "canceled":
			return coreModels.CanceledStatus
		}
	}

	return coreModels.UnknownStatus
}

func parseReleaseStatus(status string) coreModels.TileStatus {
	switch status {
	case "failed":
		return coreModels.FailedStatus
	case "succeeded":
		return coreModels.SuccessStatus
	case "partiallySucceeded":
		return coreModels.WarningStatus
	case "inProgress":
		return coreModels.RunningStatus
	}

	// all / notDeployed
	return coreModels.UnknownStatus
}
