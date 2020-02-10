//+build faker

package usecase

import (
	"fmt"
	"time"

	"github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorable/travisci"
	travisModels "github.com/monitoror/monitoror/monitorable/travisci/models"
	"github.com/monitoror/monitoror/pkg/monitoror/faker"
	"github.com/monitoror/monitoror/pkg/monitoror/utils/git"
	"github.com/monitoror/monitoror/pkg/monitoror/utils/nonempty"

	"github.com/AlekSi/pointer"
)

type (
	travisCIUsecase struct {
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

func NewTravisCIUsecase() travisci.Usecase {
	return &travisCIUsecase{make(map[string]time.Time)}
}

func (tu *travisCIUsecase) Build(params *travisModels.BuildParams) (tile *models.Tile, err error) {
	tile = models.NewTile(travisci.TravisCIBuildTileType).WithBuild()
	tile.Label = params.Repository
	tile.Build.Branch = pointer.ToString(git.HumanizeBranch(params.Branch))

	tile.Status = nonempty.Struct(params.Status, tu.computeStatus(params)).(models.TileStatus)

	if tile.Status == models.WarningStatus {
		tile.Message = "Fake error message"
		return
	}

	tile.Build.ID = pointer.ToString("12")
	tile.Build.PreviousStatus = nonempty.Struct(params.PreviousStatus, models.SuccessStatus).(models.TileStatus)

	// Author
	if tile.Status == models.FailedStatus {
		tile.Build.Author = &models.Author{}
		tile.Build.Author.Name = nonempty.String(params.AuthorName, "John Doe")
		tile.Build.Author.AvatarURL = nonempty.String(params.AuthorAvatarURL, "https://monitoror.com/assets/images/avatar.png")
	}

	if tile.Status == models.RunningStatus {
		estimatedDuration := nonempty.Duration(time.Duration(params.EstimatedDuration), time.Second*300)
		tile.Build.Duration = pointer.ToInt64(nonempty.Int64(params.Duration, int64(tu.computeDuration(params, estimatedDuration).Seconds())))

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

func (tu *travisCIUsecase) computeStatus(params *travisModels.BuildParams) models.TileStatus {
	projectID := fmt.Sprintf("%s-%s-%s", params.Owner, params.Repository, params.Branch)
	value, ok := tu.timeRefByProject[projectID]
	if !ok {
		tu.timeRefByProject[projectID] = faker.GetRefTime()
	}

	return faker.ComputeStatus(value, availableBuildStatus)
}

func (tu *travisCIUsecase) computeDuration(params *travisModels.BuildParams, duration time.Duration) time.Duration {
	projectID := fmt.Sprintf("%s-%s-%s", params.Owner, params.Repository, params.Branch)
	value, ok := tu.timeRefByProject[projectID]
	if !ok {
		tu.timeRefByProject[projectID] = faker.GetRefTime()
	}

	return faker.ComputeDuration(value, duration)
}
