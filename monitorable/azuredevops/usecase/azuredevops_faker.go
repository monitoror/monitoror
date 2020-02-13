//+build faker

package usecase

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/monitoror/monitoror/pkg/monitoror/faker"

	"github.com/AlekSi/pointer"

	"github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorable/azuredevops"
	azureModels "github.com/monitoror/monitoror/monitorable/azuredevops/models"
	"github.com/monitoror/monitoror/pkg/monitoror/utils/git"
	"github.com/monitoror/monitoror/pkg/monitoror/utils/nonempty"
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
	tile = models.NewTile(azuredevops.AzureDevOpsBuildTileType)

	branch := "master"
	if params.Branch != nil {
		branch = *params.Branch
	}
	tile.Label = fmt.Sprintf("%s (%d)\n%s - #12", params.Project, *params.Definition, git.HumanizeBranch(branch))

	tile.Status = nonempty.Struct(params.Status, tu.computeStatus(params.Project, params.Definition, availableBuildStatus)).(models.TileStatus)

	if tile.Status == models.WarningStatus {
		// Warning can be Unstable Build
		if rand.Intn(2) == 0 {
			tile.Message = "Fake error message"
			return
		}
	}

	tile.PreviousStatus = nonempty.Struct(params.PreviousStatus, models.SuccessStatus).(models.TileStatus)

	// Author
	if tile.Status != models.QueuedStatus {
		tile.Author = &models.Author{}
		tile.Author.Name = nonempty.String(params.AuthorName, "John Doe")
		tile.Author.AvatarURL = nonempty.String(params.AuthorAvatarURL, "https://monitoror.com/assets/images/avatar.png")
	}

	// StartedAt / FinishedAt
	if tile.Status == models.SuccessStatus || tile.Status == models.FailedStatus || tile.Status == models.WarningStatus || tile.Status == models.CanceledStatus {
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
		estimatedDuration := nonempty.Duration(time.Duration(params.EstimatedDuration), time.Second*300)
		tile.Duration = pointer.ToInt64(nonempty.Int64(params.Duration, int64(tu.computeDuration(params.Project, params.Definition, estimatedDuration).Seconds())))

		if tile.PreviousStatus != models.UnknownStatus {
			tile.EstimatedDuration = pointer.ToInt64(int64(estimatedDuration.Seconds()))
		} else {
			tile.EstimatedDuration = pointer.ToInt64(0)
		}
	}

	return
}

func (tu *azureDevOpsUsecase) Release(params *azureModels.ReleaseParams) (tile *models.Tile, err error) {
	tile = models.NewTile(azuredevops.AzureDevOpsReleaseTileType)
	tile.Label = fmt.Sprintf("%s (%d)\n#12", params.Project, *params.Definition)

	tile.Status = nonempty.Struct(params.Status, tu.computeStatus(params.Project, params.Definition, availableReleaseStatus)).(models.TileStatus)

	if tile.Status == models.WarningStatus {
		// Warning can be Unstable Build
		if rand.Intn(2) == 0 {
			tile.Message = "Fake error message"
			return
		}
	}

	tile.PreviousStatus = nonempty.Struct(params.PreviousStatus, models.SuccessStatus).(models.TileStatus)

	// Author
	if tile.Status == models.FailedStatus {
		tile.Author = &models.Author{}
		tile.Author.Name = nonempty.String(params.AuthorName, "John Doe")
		tile.Author.AvatarURL = nonempty.String(params.AuthorAvatarURL, "https://monitoror.com/assets/images/avatar.png")
	}

	// Duration / EstimatedDuration
	if tile.Status == models.RunningStatus {
		estimatedDuration := nonempty.Duration(time.Duration(params.EstimatedDuration), time.Second*300)
		tile.Duration = pointer.ToInt64(nonempty.Int64(params.Duration, int64(tu.computeDuration(params.Project, params.Definition, estimatedDuration).Seconds())))

		if tile.PreviousStatus != models.UnknownStatus {
			tile.EstimatedDuration = pointer.ToInt64(int64(estimatedDuration.Seconds()))
		} else {
			tile.EstimatedDuration = pointer.ToInt64(0)
		}
	}

	// StartedAt / FinishedAt
	if tile.Duration == nil {
		tile.StartedAt = pointer.ToTime(nonempty.Time(params.StartedAt, time.Now().Add(-time.Minute*10)))
	} else {
		tile.StartedAt = pointer.ToTime(nonempty.Time(params.StartedAt, time.Now().Add(-time.Second*time.Duration(*tile.Duration))))
	}

	if tile.Status != models.QueuedStatus && tile.Status != models.RunningStatus {
		tile.FinishedAt = pointer.ToTime(nonempty.Time(params.FinishedAt, tile.StartedAt.Add(time.Minute*5)))
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
