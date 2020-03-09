//+build faker

package usecase

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/AlekSi/pointer"
	"github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorable/github"
	githubModels "github.com/monitoror/monitoror/monitorable/github/models"
	"github.com/monitoror/monitoror/pkg/monitoror/builder"
	"github.com/monitoror/monitoror/pkg/monitoror/faker"
	"github.com/monitoror/monitoror/pkg/monitoror/utils/git"
	"github.com/monitoror/monitoror/pkg/monitoror/utils/nonempty"

	cmap "github.com/orcaman/concurrent-map"
)

type (
	githubUsecase struct {
		timeRefByProject cmap.ConcurrentMap
	}
)

var availableBuildStatus = faker.Statuses{
	{models.SuccessStatus, time.Second * 30},
	{models.FailedStatus, time.Second * 30},
	{models.CanceledStatus, time.Second * 20},
	{models.RunningStatus, time.Second * 60},
	{models.QueuedStatus, time.Second * 30},
	{models.WarningStatus, time.Second * 20},
	{models.DisabledStatus, time.Second * 20},
	{models.ActionRequiredStatus, time.Second * 20},
}

func NewGithubUsecase() github.Usecase {
	return &githubUsecase{cmap.New()}
}

func (gu *githubUsecase) Count(params *githubModels.CountParams) (*models.Tile, error) {
	tile := models.NewTile(github.GithubCountTileType).WithValue(models.NumberUnit)
	tile.Label = params.Query

	tile.Status = models.SuccessStatus

	if len(params.ValueValues) != 0 {
		tile.Value.Values = params.ValueValues
	} else {
		tile.Value.Values = []string{"42"}
	}

	return tile, nil
}

func (gu *githubUsecase) Checks(params *githubModels.ChecksParams) (tile *models.Tile, err error) {
	tile = models.NewTile(github.GithubChecksTileType).WithBuild()
	tile.Label = params.Repository

	tile.Status = nonempty.Struct(params.Status, gu.computeStatus(params)).(models.TileStatus)

	if tile.Status == models.DisabledStatus {
		return
	}

	if tile.Status == models.WarningStatus {
		// Warning can be Unstable Build
		if rand.Intn(2) == 0 {
			tile.Message = "Fake error message"
			return
		}
	}

	tile.Build.Branch = pointer.ToString(git.HumanizeBranch(params.Ref))
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
		tile.Build.Duration = pointer.ToInt64(nonempty.Int64(params.Duration, int64(gu.computeDuration(params, estimatedDuration).Seconds())))

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

	return tile, nil
}

func (gu *githubUsecase) ListDynamicTile(params interface{}) ([]builder.Result, error) {
	panic("unimplemented")
}

func (gu *githubUsecase) computeStatus(params *githubModels.ChecksParams) models.TileStatus {
	projectID := fmt.Sprintf("%s-%s-%s", params.Owner, params.Repository, params.Ref)
	value, ok := gu.timeRefByProject.Get(projectID)
	if !ok || value == nil {
		value = faker.GetRefTime()
		gu.timeRefByProject.Set(projectID, value)
	}

	return faker.ComputeStatus(value.(time.Time), availableBuildStatus)
}

func (gu *githubUsecase) computeDuration(params *githubModels.ChecksParams, duration time.Duration) time.Duration {
	projectID := fmt.Sprintf("%s-%s-%s", params.Owner, params.Repository, params.Ref)
	value, ok := gu.timeRefByProject.Get(projectID)
	if !ok || value == nil {
		value = faker.GetRefTime()
		gu.timeRefByProject.Set(projectID, value)
	}

	return faker.ComputeDuration(value.(time.Time), duration)
}
