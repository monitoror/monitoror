//+build !faker

package usecase

import (
	"fmt"
	"time"

	"github.com/AlekSi/pointer"

	"github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorable/azuredevops"
	azureModels "github.com/monitoror/monitoror/monitorable/azuredevops/models"
	"github.com/monitoror/monitoror/pkg/monitoror/cache"
	"github.com/monitoror/monitoror/pkg/monitoror/utils/git"
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
	tile := models.NewTile(azuredevops.AzureDevOpsBuildTileType)
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
	tile.Label = fmt.Sprintf("%s (%s)\n%s - #%s", params.Project, build.DefinitionName, git.HumanizeBranch(build.Branch), build.BuildNumber)

	// Status
	tile.Status = parseBuildResult(build.Status, build.Result)

	// Previous Status
	previousStatus := au.buildsCache.GetPreviousStatus(params, build.BuildNumber)
	if previousStatus != nil {
		tile.PreviousStatus = *previousStatus
	} else {
		tile.PreviousStatus = models.UnknownStatus
	}

	// Author
	if build.Author != nil {
		tile.Author = &models.Author{
			Name:      build.Author.Name,
			AvatarURL: build.Author.AvatarURL,
		}
	}

	// StartedAt / FinishedAt
	tile.StartedAt = build.StartedAt
	if tile.Status != models.QueuedStatus && tile.Status != models.RunningStatus {
		tile.FinishedAt = build.FinishedAt
	}

	if tile.Status == models.QueuedStatus {
		tile.StartedAt = build.QueuedAt
	}

	// Duration / Previous Duration
	if tile.Status == models.RunningStatus {
		tile.Duration = pointer.ToInt64(int64(time.Since(*tile.StartedAt).Seconds()))

		estimatedDuration := au.buildsCache.GetEstimatedDuration(params)
		if estimatedDuration != nil {
			tile.EstimatedDuration = pointer.ToInt64(int64(estimatedDuration.Seconds()))
		} else {
			tile.EstimatedDuration = pointer.ToInt64(int64(0))
		}
	}

	// Cache Duration when success / failed
	if tile.Status == models.SuccessStatus || tile.Status == models.FailedStatus {
		au.buildsCache.Add(params, build.BuildNumber, tile.Status, tile.FinishedAt.Sub(*tile.StartedAt))
	}

	return tile, nil
}

func (au *azureDevOpsUsecase) Release(params *azureModels.ReleaseParams) (*models.Tile, error) {
	tile := models.NewTile(azuredevops.AzureDevOpsReleaseTileType)
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
	tile.Label = fmt.Sprintf("%s (%s)\n#%s", params.Project, release.DefinitionName, release.ReleaseNumber)

	// Status
	tile.Status = parseReleaseStatus(release.Status)

	// Previous Status
	previousStatus := au.buildsCache.GetPreviousStatus(params, release.ReleaseNumber)
	if previousStatus != nil {
		tile.PreviousStatus = *previousStatus
	} else {
		tile.PreviousStatus = models.UnknownStatus
	}

	// Author
	if release.Author != nil {
		tile.Author = &models.Author{
			Name:      release.Author.Name,
			AvatarURL: release.Author.AvatarURL,
		}
	}

	// StartedAt / FinishedAt
	tile.StartedAt = release.StartedAt
	tile.StartedAt = release.StartedAt
	if tile.Status != models.RunningStatus {
		tile.FinishedAt = release.FinishedAt
	}
	// Duration
	if tile.Status == models.RunningStatus {
		tile.Duration = pointer.ToInt64(int64(time.Since(*tile.StartedAt).Seconds()))

		estimatedDuration := au.buildsCache.GetEstimatedDuration(params)
		if estimatedDuration != nil {
			tile.EstimatedDuration = pointer.ToInt64(int64(estimatedDuration.Seconds()))
		} else {
			tile.EstimatedDuration = pointer.ToInt64(int64(0))
		}
	}

	// Cache Duration when success / failed
	if tile.Status == models.SuccessStatus || tile.Status == models.FailedStatus {
		au.buildsCache.Add(params, release.ReleaseNumber, tile.Status, tile.FinishedAt.Sub(*tile.StartedAt))
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
