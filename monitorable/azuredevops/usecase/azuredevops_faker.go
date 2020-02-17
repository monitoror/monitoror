//+build faker

package usecase

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorable/azuredevops"
	azureModels "github.com/monitoror/monitoror/monitorable/azuredevops/models"
	"github.com/monitoror/monitoror/pkg/monitoror/faker"
	"github.com/monitoror/monitoror/pkg/monitoror/utils/git"
	"github.com/monitoror/monitoror/pkg/monitoror/utils/nonempty"

	"github.com/AlekSi/pointer"
)

type (
	azureDevOpsUsecase struct {
		timeRefByProject map[string]time.Time
	}
)

var availableBuildStatus = faker.Statuses{
	{models.SuccessStatus, time.Second * 30},
	{models.FailedStatus, time.Second * 30},
	{models.CanceledStatus, time.Second * 20},
	{models.RunningStatus, time.Second * 60},
	{models.QueuedStatus, time.Second * 30},
	{models.WarningStatus, time.Second * 20},
}

var availableReleaseStatus = faker.Statuses{
	{models.SuccessStatus, time.Second * 30},
	{models.FailedStatus, time.Second * 30},
	{models.RunningStatus, time.Second * 60},
	{models.WarningStatus, time.Second * 20},
}

func NewAzureDevOpsUsecase() azuredevops.Usecase {
	return &azureDevOpsUsecase{make(map[string]time.Time)}
}

func (tu *azureDevOpsUsecase) Build(params *azureModels.BuildParams) (tile *models.Tile, err error) {
	tile = models.NewTile(azuredevops.AzureDevOpsBuildTileType).WithBuild()
	tile.Label = fmt.Sprintf("%s (build-qa-%d)", params.Project, *params.Definition)
	tile.Build.ID = pointer.ToString("12")
	tile.Build.Branch = pointer.ToString(git.HumanizeBranch(nonempty.String(*params.Branch, "master")))

	tile.Status = nonempty.Struct(params.Status, tu.computeStatus(params.Project, params.Definition, availableBuildStatus)).(models.TileStatus)

	if tile.Status == models.WarningStatus {
		// Warning can be Unstable Build
		if rand.Intn(2) == 0 {
			tile.Message = "Fake error message"
			return
		}
	}

	tile.Build.PreviousStatus = nonempty.Struct(params.PreviousStatus, models.SuccessStatus).(models.TileStatus)

	// Author
	if tile.Status != models.QueuedStatus {
		tile.Build.Author = &models.Author{}
		tile.Build.Author.Name = nonempty.String(params.AuthorName, "John Doe")
		tile.Build.Author.AvatarURL = nonempty.String(params.AuthorAvatarURL, "https://monitoror.com/assets/images/avatar.png")
	}

	// StartedAt / FinishedAt
	if tile.Status == models.SuccessStatus || tile.Status == models.FailedStatus || tile.Status == models.WarningStatus || tile.Status == models.CanceledStatus {
		min := time.Now().Unix() - int64(time.Hour.Seconds()*24*30) - 3600
		max := time.Now().Unix() - 3600
		delta := max - min

		tile.Build.StartedAt = pointer.ToTime(nonempty.Time(params.StartedAt, time.Unix(rand.Int63n(delta)+min, 0)))
		tile.Build.FinishedAt = pointer.ToTime(nonempty.Time(params.FinishedAt, tile.Build.StartedAt.Add(time.Second*time.Duration(rand.Int63n(3600)))))
	}
	if tile.Status == models.QueuedStatus || tile.Status == models.RunningStatus {
		tile.Build.StartedAt = pointer.ToTime(nonempty.Time(params.StartedAt, time.Now().Add(-time.Second*time.Duration(rand.Int63n(3600)))))
	}

	// Duration / EstimatedDuration
	if tile.Status == models.RunningStatus {
		estimatedDuration := nonempty.Duration(time.Duration(params.EstimatedDuration), time.Second*300)
		tile.Build.Duration = pointer.ToInt64(nonempty.Int64(params.Duration, int64(tu.computeDuration(params.Project, params.Definition, estimatedDuration).Seconds())))

		if tile.Build.PreviousStatus != models.UnknownStatus {
			tile.Build.EstimatedDuration = pointer.ToInt64(int64(estimatedDuration.Seconds()))
		} else {
			tile.Build.EstimatedDuration = pointer.ToInt64(0)
		}
	}

	return
}

func (tu *azureDevOpsUsecase) Release(params *azureModels.ReleaseParams) (tile *models.Tile, err error) {
	tile = models.NewTile(azuredevops.AzureDevOpsReleaseTileType).WithBuild()
	tile.Label = fmt.Sprintf("%s (release-%d)", params.Project, *params.Definition)
	tile.Build.ID = pointer.ToString("12")

	tile.Status = nonempty.Struct(params.Status, tu.computeStatus(params.Project, params.Definition, availableReleaseStatus)).(models.TileStatus)

	if tile.Status == models.WarningStatus {
		// Warning can be Unstable Build
		if rand.Intn(2) == 0 {
			tile.Message = "Fake error message"
			return
		}
	}

	tile.Build.PreviousStatus = nonempty.Struct(params.PreviousStatus, models.SuccessStatus).(models.TileStatus)

	// Author
	if tile.Status == models.FailedStatus {
		tile.Build.Author = &models.Author{}
		tile.Build.Author.Name = nonempty.String(params.AuthorName, "John Doe")
		tile.Build.Author.AvatarURL = nonempty.String(params.AuthorAvatarURL, "https://monitoror.com/assets/images/avatar.png")
	}

	// Duration / EstimatedDuration
	if tile.Status == models.RunningStatus {
		estimatedDuration := nonempty.Duration(time.Duration(params.EstimatedDuration), time.Second*300)
		tile.Build.Duration = pointer.ToInt64(nonempty.Int64(params.Duration, int64(tu.computeDuration(params.Project, params.Definition, estimatedDuration).Seconds())))

		if tile.Build.PreviousStatus != models.UnknownStatus {
			tile.Build.EstimatedDuration = pointer.ToInt64(int64(estimatedDuration.Seconds()))
		} else {
			tile.Build.EstimatedDuration = pointer.ToInt64(0)
		}
	}

	// StartedAt / FinishedAt
	if tile.Build.Duration == nil {
		tile.Build.StartedAt = pointer.ToTime(nonempty.Time(params.StartedAt, time.Now().Add(-time.Minute*10)))
	} else {
		tile.Build.StartedAt = pointer.ToTime(nonempty.Time(params.StartedAt, time.Now().Add(-time.Second*time.Duration(*tile.Build.Duration))))
	}

	if tile.Status != models.QueuedStatus && tile.Status != models.RunningStatus {
		tile.Build.FinishedAt = pointer.ToTime(nonempty.Time(params.FinishedAt, tile.Build.StartedAt.Add(time.Minute*5)))
	}

	return
}

func (tu *azureDevOpsUsecase) computeStatus(project string, definition *int, statuses faker.Statuses) models.TileStatus {
	projectID := fmt.Sprintf("%s-%d", project, *definition)
	value, ok := tu.timeRefByProject[projectID]
	if !ok {
		tu.timeRefByProject[projectID] = faker.GetRefTime()
	}

	return faker.ComputeStatus(value, statuses)
}

func (tu *azureDevOpsUsecase) computeDuration(project string, definition *int, duration time.Duration) time.Duration {
	projectID := fmt.Sprintf("%s-%d", project, *definition)
	value, ok := tu.timeRefByProject[projectID]
	if !ok {
		tu.timeRefByProject[projectID] = faker.GetRefTime()
	}

	return faker.ComputeDuration(value, duration)
}
