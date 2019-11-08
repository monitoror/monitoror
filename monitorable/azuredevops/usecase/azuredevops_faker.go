//+build faker

package usecase

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/AlekSi/pointer"

	"github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorable/azuredevops"
	azureModels "github.com/monitoror/monitoror/monitorable/azuredevops/models"
	"github.com/monitoror/monitoror/pkg/monitoror/utils/git"
	"github.com/monitoror/monitoror/pkg/monitoror/utils/nonempty"
)

var AvailableBuildStatus = []models.TileStatus{
	models.SuccessStatus,
	models.FailedStatus,
	models.AbortedStatus,
	models.RunningStatus,
	models.QueuedStatus,
	models.WarningStatus,
}
var AvailableReleaseStatus = []models.TileStatus{
	models.SuccessStatus,
	models.FailedStatus,
	models.RunningStatus,
	models.WarningStatus,
}
var AvailablePreviousStatus = []models.TileStatus{
	models.SuccessStatus,
	models.FailedStatus,
	models.WarningStatus,
	models.UnknownStatus,
}

type (
	azureDevOpsUsecase struct {
		cachedRunningValue map[string]*durations
	}

	durations struct {
		duration          int64
		estimatedDuration int64
	}
)

func NewAzureDevOpsUsecase() azuredevops.Usecase {
	return &azureDevOpsUsecase{make(map[string]*durations)}
}

func (tu *azureDevOpsUsecase) Build(params *azureModels.BuildParams) (tile *models.Tile, err error) {
	tile = models.NewTile(azuredevops.AzureDevOpsBuildTileType)
	tile.Label = fmt.Sprintf("%s | %d", params.Project, *params.Definition)

	// Init random generator
	rand.Seed(time.Now().UnixNano())

	branch := "master"
	if params.Branch != nil {
		branch = *params.Branch
	}
	tile.Message = fmt.Sprintf("%s - %d", git.HumanizeBranch(branch), rand.Intn(100))
	tile.Status = nonempty.Struct(params.Status, randomStatus(AvailableBuildStatus)).(models.TileStatus)

	if tile.Status == models.WarningStatus {
		// Warning can be Unstable Build
		if rand.Intn(2) == 0 {
			tile.Message = "random error message"
			return
		}
	}

	tile.PreviousStatus = nonempty.Struct(params.PreviousStatus, randomStatus(AvailablePreviousStatus)).(models.TileStatus)

	// Author
	if tile.Status != models.QueuedStatus {
		tile.Author = &models.Author{}
		tile.Author.Name = nonempty.String(params.AuthorName, "Faker")
		tile.Author.AvatarURL = nonempty.String(params.AuthorAvatarURL, "https://www.gravatar.com/avatar/00000000000000000000000000000000")
	}

	// StartedAt / FinishedAt
	if tile.Status == models.SuccessStatus || tile.Status == models.FailedStatus || tile.Status == models.WarningStatus || tile.Status == models.AbortedStatus {
		min := time.Now().Unix() - int64(time.Hour.Seconds()*24*30) - 3600
		max := time.Now().Unix() - 3600
		delta := max - min

		tile.StartedAt = pointer.ToTime(nonempty.Time(params.StartedAt, time.Unix(rand.Int63n(delta)+min, 0)))
		tile.FinishedAt = pointer.ToTime(nonempty.Time(params.FinishedAt, tile.StartedAt.Add(time.Second*time.Duration(rand.Int63n(3600)))))
	}
	if tile.Status == models.QueuedStatus || tile.Status == models.RunningStatus {
		tile.StartedAt = pointer.ToTime(nonempty.Time(params.StartedAt, time.Now().Add(-time.Second*time.Duration(rand.Int63n(3600)))))
	}

	// Duration / EstimatedDuration
	if tile.Status == models.RunningStatus {
		// Creating cache for duration
		dur, ok := tu.cachedRunningValue[tile.Label]
		if !ok {
			dur = &durations{}
			tu.cachedRunningValue[tile.Label] = dur
		}

		// Test if there is cached value or if user force value with param
		if dur.estimatedDuration == 0 || params.EstimatedDuration != 0 {
			dur.estimatedDuration = nonempty.Int64(params.EstimatedDuration, rand.Int63n(340)+10)
		}

		// Increment cached Duration
		dur.duration += 10
		if dur.duration > dur.estimatedDuration {
			dur.duration = 0
		}

		tile.Duration = pointer.ToInt64(nonempty.Int64(params.Duration, dur.duration))

		if tile.PreviousStatus != models.UnknownStatus {
			tile.EstimatedDuration = pointer.ToInt64(dur.estimatedDuration)
		} else {
			tile.EstimatedDuration = pointer.ToInt64(0)
		}
	}

	return
}

func (tu *azureDevOpsUsecase) Release(params *azureModels.ReleaseParams) (tile *models.Tile, err error) {
	tile = models.NewTile(azuredevops.AzureDevOpsReleaseTileType)
	tile.Label = fmt.Sprintf("%s | %d", params.Project, *params.Definition)

	// Init random generator
	rand.Seed(time.Now().UnixNano())

	tile.Status = nonempty.Struct(params.Status, randomStatus(AvailableReleaseStatus)).(models.TileStatus)

	if tile.Status == models.WarningStatus {
		// Warning can be Unstable Build
		if rand.Intn(2) == 0 {
			tile.Message = "random error message"
			return
		}
	}

	tile.PreviousStatus = nonempty.Struct(params.PreviousStatus, randomStatus(AvailablePreviousStatus)).(models.TileStatus)

	// Author
	tile.Author = &models.Author{}
	tile.Author.Name = nonempty.String(params.AuthorName, "Faker")
	tile.Author.AvatarURL = nonempty.String(params.AuthorAvatarURL, "https://www.gravatar.com/avatar/00000000000000000000000000000000")

	// StartedAt / FinishedAt
	if tile.Status == models.SuccessStatus || tile.Status == models.FailedStatus || tile.Status == models.WarningStatus || tile.Status == models.AbortedStatus {
		min := time.Now().Unix() - int64(time.Hour.Seconds()*24*30) - 3600
		max := time.Now().Unix() - 3600
		delta := max - min

		tile.StartedAt = pointer.ToTime(nonempty.Time(params.StartedAt, time.Unix(rand.Int63n(delta)+min, 0)))
		tile.FinishedAt = pointer.ToTime(nonempty.Time(params.FinishedAt, tile.StartedAt.Add(time.Second*time.Duration(rand.Int63n(3600)))))
	}
	if tile.Status == models.RunningStatus {
		tile.StartedAt = pointer.ToTime(nonempty.Time(params.StartedAt, time.Now().Add(-time.Second*time.Duration(rand.Int63n(3600)))))
	}

	// Duration / EstimatedDuration
	if tile.Status == models.RunningStatus {
		// Creating cache for duration
		dur, ok := tu.cachedRunningValue[tile.Label]
		if !ok {
			dur = &durations{}
			tu.cachedRunningValue[tile.Label] = dur
		}

		// Test if there is cached value or if user force value with param
		if dur.estimatedDuration == 0 || params.EstimatedDuration != 0 {
			dur.estimatedDuration = nonempty.Int64(params.EstimatedDuration, rand.Int63n(340)+10)
		}

		// Increment cached Duration
		dur.duration += 10
		if dur.duration > dur.estimatedDuration {
			dur.duration = 0
		}

		tile.Duration = pointer.ToInt64(nonempty.Int64(params.Duration, dur.duration))

		if tile.PreviousStatus != models.UnknownStatus {
			tile.EstimatedDuration = pointer.ToInt64(dur.estimatedDuration)
		} else {
			tile.EstimatedDuration = pointer.ToInt64(0)
		}
	}

	return
}

func randomStatus(status []models.TileStatus) models.TileStatus {
	return status[rand.Intn(len(status))]
}
