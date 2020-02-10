//+build !faker

package usecase

import (
	"fmt"
	"time"

	"github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorable/azuredevops"
	azureModels "github.com/monitoror/monitoror/monitorable/azuredevops/models"
	"github.com/monitoror/monitoror/pkg/monitoror/cache"
	"github.com/monitoror/monitoror/pkg/monitoror/utils/git"

	"github.com/AlekSi/pointer"
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

func (au *azureDevOpsUsecase) Build(params *azureModels.BuildParams) (*models.Tile, error) {
	tile := models.NewTile(azuredevops.AzureDevOpsBuildTileType).WithBuild()
	// Default label if build not found
	tile.Label = params.Project

	// Lookup for build
	build, err := au.repository.GetBuild(params.Project, *params.Definition, params.Branch)
	if err != nil {
		return nil, &models.MonitororError{Err: err, Tile: tile, Message: "unable to find build"}
	}
	if build == nil {
		return nil, &models.MonitororError{Tile: tile, Message: "no build found", ErrorStatus: models.UnknownStatus}
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
		tile.Build.PreviousStatus = models.UnknownStatus
	}

	// Author
	if tile.Status == models.FailedStatus && build.Author != nil {
		tile.Build.Author = &models.Author{
			Name:      build.Author.Name,
			AvatarURL: build.Author.AvatarURL,
		}
	}

	// StartedAt / FinishedAt
	tile.Build.StartedAt = build.StartedAt
	if tile.Status != models.QueuedStatus && tile.Status != models.RunningStatus {
		tile.Build.FinishedAt = build.FinishedAt
	}

	if tile.Status == models.QueuedStatus {
		tile.Build.StartedAt = build.QueuedAt
	}

	// Duration / Previous Duration
	if tile.Status == models.RunningStatus {
		tile.Build.Duration = pointer.ToInt64(int64(time.Since(*tile.Build.StartedAt).Seconds()))

		estimatedDuration := au.buildsCache.GetEstimatedDuration(params)
		if estimatedDuration != nil {
			tile.Build.EstimatedDuration = pointer.ToInt64(int64(estimatedDuration.Seconds()))
		} else {
			tile.Build.EstimatedDuration = pointer.ToInt64(int64(0))
		}
	}

	// Cache Duration when success / failed
	if tile.Status == models.SuccessStatus || tile.Status == models.FailedStatus {
		au.buildsCache.Add(params, build.BuildNumber, tile.Status, tile.Build.FinishedAt.Sub(*tile.Build.StartedAt))
	}

	return tile, nil
}

func (au *azureDevOpsUsecase) Release(params *azureModels.ReleaseParams) (*models.Tile, error) {
	tile := models.NewTile(azuredevops.AzureDevOpsReleaseTileType).WithBuild()
	// Default label if build not found
	tile.Label = params.Project

	// Lookup for release
	release, err := au.repository.GetRelease(params.Project, *params.Definition)
	if err != nil {
		return nil, &models.MonitororError{Err: err, Tile: tile, Message: "unable to find release"}
	}
	if release == nil {
		// Warning because request was correct but there is no release
		return nil, &models.MonitororError{Tile: tile, Message: "no release found", ErrorStatus: models.UnknownStatus}
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
		tile.Build.PreviousStatus = models.UnknownStatus
	}

	// Author
	if tile.Status == models.FailedStatus && release.Author != nil {
		tile.Build.Author = &models.Author{
			Name:      release.Author.Name,
			AvatarURL: release.Author.AvatarURL,
		}
	}

	// StartedAt / FinishedAt
	tile.Build.StartedAt = release.StartedAt
	if tile.Status != models.RunningStatus && tile.Status != models.QueuedStatus {
		tile.Build.FinishedAt = release.FinishedAt
	}
	// Duration
	if tile.Status == models.RunningStatus {
		tile.Build.Duration = pointer.ToInt64(int64(time.Since(*tile.Build.StartedAt).Seconds()))

		estimatedDuration := au.buildsCache.GetEstimatedDuration(params)
		if estimatedDuration != nil {
			tile.Build.EstimatedDuration = pointer.ToInt64(int64(estimatedDuration.Seconds()))
		} else {
			tile.Build.EstimatedDuration = pointer.ToInt64(int64(0))
		}
	}

	// Cache Duration when success / failed
	if tile.Status == models.SuccessStatus || tile.Status == models.FailedStatus || tile.Status == models.WarningStatus {
		au.buildsCache.Add(params, *tile.Build.ID, tile.Status, tile.Build.FinishedAt.Sub(*tile.Build.StartedAt))
	}

	return tile, nil
}

func parseBuildResult(status, result string) models.TileStatus {
	switch status {
	case "inProgress":
		return models.RunningStatus
	case "cancelling":
		return models.RunningStatus
	case "notStarted":
		return models.QueuedStatus
	case "completed":
		switch result {
		case "succeeded":
			return models.SuccessStatus
		case "partiallySucceeded":
			return models.WarningStatus
		case "failed":
			return models.FailedStatus
		case "canceled":
			return models.CanceledStatus
		}
	}

	return models.UnknownStatus
}

func parseReleaseStatus(status string) models.TileStatus {
	switch status {
	case "failed":
		return models.FailedStatus
	case "succeeded":
		return models.SuccessStatus
	case "partiallySucceeded":
		return models.WarningStatus
	case "inProgress":
		return models.RunningStatus
	}

	// all / notDeployed
	return models.UnknownStatus
}
