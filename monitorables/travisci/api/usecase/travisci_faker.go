//+build faker

package usecase

import (
	"fmt"
	"time"

	"github.com/monitoror/monitoror/internal/pkg/monitorable/faker"
	coreModels "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorables/travisci/api"
	"github.com/monitoror/monitoror/monitorables/travisci/api/models"
	"github.com/monitoror/monitoror/pkg/git"
	"github.com/monitoror/monitoror/pkg/nonempty"

	"github.com/AlekSi/pointer"
	cmap "github.com/orcaman/concurrent-map"
)

type (
	travisCIUsecase struct {
		timeRefByProject cmap.ConcurrentMap
	}
)

var availableBuildStatus = faker.Statuses{
	{coreModels.SuccessStatus, time.Second * 30},
	{coreModels.FailedStatus, time.Second * 30},
	{coreModels.CanceledStatus, time.Second * 20},
	{coreModels.RunningStatus, time.Second * 60},
	{coreModels.QueuedStatus, time.Second * 30},
	{coreModels.WarningStatus, time.Second * 20},
}

func NewTravisCIUsecase() api.Usecase {
	return &travisCIUsecase{cmap.New()}
}

func (tu *travisCIUsecase) Build(params *models.BuildParams) (tile *coreModels.Tile, err error) {
	tile = coreModels.NewTile(api.TravisCIBuildTileType).WithBuild()
	tile.Label = params.Repository
	tile.Build.Branch = pointer.ToString(git.HumanizeBranch(params.Branch))

	tile.Status = nonempty.Struct(params.Status, tu.computeStatus(params)).(coreModels.TileStatus)

	if tile.Status == coreModels.WarningStatus {
		tile.Message = "Fake error message"
		return
	}

	tile.Build.ID = pointer.ToString("12")
	tile.Build.PreviousStatus = nonempty.Struct(params.PreviousStatus, coreModels.SuccessStatus).(coreModels.TileStatus)

	// Author
	if tile.Status == coreModels.FailedStatus {
		tile.Build.Author = &coreModels.Author{}
		tile.Build.Author.Name = nonempty.String(params.AuthorName, "John Doe")
		tile.Build.Author.AvatarURL = nonempty.String(params.AuthorAvatarURL, "https://monitoror.com/assets/images/avatar.png")
	}

	if tile.Status == coreModels.RunningStatus {
		estimatedDuration := nonempty.Duration(time.Duration(params.EstimatedDuration), time.Second*300)
		tile.Build.Duration = pointer.ToInt64(nonempty.Int64(params.Duration, int64(tu.computeDuration(params, estimatedDuration).Seconds())))

		if tile.Build.PreviousStatus != coreModels.UnknownStatus {
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

	if tile.Status != coreModels.QueuedStatus && tile.Status != coreModels.RunningStatus {
		tile.Build.FinishedAt = pointer.ToTime(nonempty.Time(params.FinishedAt, tile.Build.StartedAt.Add(time.Minute*5)))
	}

	return
}

func (tu *travisCIUsecase) computeStatus(params *models.BuildParams) coreModels.TileStatus {
	projectID := fmt.Sprintf("%s-%s-%s", params.Owner, params.Repository, params.Branch)
	value, ok := tu.timeRefByProject.Get(projectID)
	if !ok || value == nil {
		value = faker.GetRefTime()
		tu.timeRefByProject.Set(projectID, value)
	}

	return faker.ComputeStatus(value.(time.Time), availableBuildStatus)
}

func (tu *travisCIUsecase) computeDuration(params *models.BuildParams, duration time.Duration) time.Duration {
	projectID := fmt.Sprintf("%s-%s-%s", params.Owner, params.Repository, params.Branch)
	value, ok := tu.timeRefByProject.Get(projectID)
	if !ok || value == nil {
		value = faker.GetRefTime()
		tu.timeRefByProject.Set(projectID, value)
	}

	return faker.ComputeDuration(value.(time.Time), duration)
}
