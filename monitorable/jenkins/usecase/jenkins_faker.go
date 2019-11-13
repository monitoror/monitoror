//+build faker

package usecase

import (
	"math/rand"
	"time"

	"github.com/AlekSi/pointer"

	"github.com/monitoror/monitoror/pkg/monitoror/utils/git"

	"github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorable/jenkins"
	jenkinsModels "github.com/monitoror/monitoror/monitorable/jenkins/models"
	"github.com/monitoror/monitoror/pkg/monitoror/builder"
	"github.com/monitoror/monitoror/pkg/monitoror/utils/nonempty"
)

var AvailableStatus = []models.TileStatus{
	models.SuccessStatus,
	models.FailedStatus,
	models.AbortedStatus,
	models.RunningStatus,
	models.QueuedStatus,
	models.WarningStatus,
	models.DisabledStatus,
}
var AvailablePreviousStatus = []models.TileStatus{
	models.SuccessStatus,
	models.FailedStatus,
	models.WarningStatus,
	models.UnknownStatus,
}

type (
	jenkinsUsecase struct {
		cachedRunningValue map[string]*durations
	}

	durations struct {
		duration          int64
		estimatedDuration int64
	}
)

func NewJenkinsUsecase() jenkins.Usecase {
	return &jenkinsUsecase{make(map[string]*durations)}
}

func (tu *jenkinsUsecase) Build(params *jenkinsModels.BuildParams) (tile *models.Tile, err error) {
	tile = models.NewTile(jenkins.JenkinsBuildTileType)
	tile.Label = params.Job
	if params.Branch != "" {
		tile.Message = git.HumanizeBranch(params.Branch)
	}

	// Init random generator
	rand.Seed(time.Now().UnixNano())

	tile.Status = nonempty.Struct(params.Status, randomStatus(AvailableStatus)).(models.TileStatus)

	if tile.Status == models.DisabledStatus {
		return
	}

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
		dur, ok := tu.cachedRunningValue[params.String()]
		if !ok {
			dur = &durations{}
			tu.cachedRunningValue[params.String()] = dur
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

func (tu *jenkinsUsecase) ListDynamicTile(params interface{}) ([]builder.Result, error) {
	panic("unimplemented")
}

func randomStatus(status []models.TileStatus) models.TileStatus {
	return status[rand.Intn(len(status))]
}
