//+build faker

package usecase

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/monitoror/monitoror/pkg/monitoror/utils/git"

	. "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorable/jenkins"
	"github.com/monitoror/monitoror/monitorable/jenkins/models"
	"github.com/monitoror/monitoror/pkg/monitoror/builder"
	"github.com/monitoror/monitoror/pkg/monitoror/utils/nonempty"

	. "github.com/AlekSi/pointer"
)

var AvailableStatus = []TileStatus{SuccessStatus, FailedStatus, AbortedStatus, RunningStatus, QueuedStatus, WarningStatus, DisabledStatus}
var AvailablePreviousStatus = []TileStatus{SuccessStatus, FailedStatus, WarningStatus, UnknownStatus}

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

func (tu *jenkinsUsecase) Build(params *models.BuildParams) (tile *Tile, err error) {
	tile = NewTile(jenkins.JenkinsBuildTileType)
	if params.Branch == "" {
		tile.Label = params.Job
	} else {
		tile.Label = fmt.Sprintf("%s", params.Job)
		tile.Message = fmt.Sprintf("%s", git.HumanizeBranch(params.Branch))
	}

	// Init random generator
	rand.Seed(time.Now().UnixNano())

	tile.Status = nonempty.Struct(params.Status, randomStatus(AvailableStatus)).(TileStatus)

	if tile.Status == DisabledStatus {
		return
	}

	if tile.Status == WarningStatus {
		// Warning can be Unstable Build
		if rand.Intn(2) == 0 {
			tile.Message = "random error message"
			return
		}
	}

	tile.PreviousStatus = nonempty.Struct(params.PreviousStatus, randomStatus(AvailablePreviousStatus)).(TileStatus)

	// Author
	if tile.Status != QueuedStatus {
		tile.Author = &Author{}
		tile.Author.Name = nonempty.String(params.AuthorName, "Faker")
		tile.Author.AvatarUrl = nonempty.String(params.AuthorAvatarUrl, "https://www.gravatar.com/avatar/00000000000000000000000000000000")
	}

	// StartedAt / FinishedAt
	if tile.Status == SuccessStatus || tile.Status == FailedStatus || tile.Status == WarningStatus || tile.Status == AbortedStatus {
		min := time.Now().Unix() - int64(time.Hour.Seconds()*24*30) - 3600
		max := time.Now().Unix() - 3600
		delta := max - min

		tile.StartedAt = ToTime(nonempty.Time(params.StartedAt, time.Unix(rand.Int63n(delta)+min, 0)))
		tile.FinishedAt = ToTime(nonempty.Time(params.FinishedAt, tile.StartedAt.Add(time.Second*time.Duration(rand.Int63n(3600)))))
	}
	if tile.Status == QueuedStatus || tile.Status == RunningStatus {
		tile.StartedAt = ToTime(nonempty.Time(params.StartedAt, time.Now().Add(-time.Second*time.Duration(rand.Int63n(3600)))))
	}

	// Duration / EstimatedDuration
	if tile.Status == RunningStatus {
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

		tile.Duration = ToInt64(nonempty.Int64(params.Duration, dur.duration))

		if tile.PreviousStatus != UnknownStatus {
			tile.EstimatedDuration = ToInt64(dur.estimatedDuration)
		} else {
			tile.EstimatedDuration = ToInt64(0)
		}
	}

	return
}

func (tu *jenkinsUsecase) ListDynamicTile(params interface{}) ([]builder.Result, error) {
	panic("unimplemented")
}

func randomStatus(status []TileStatus) TileStatus {
	return status[rand.Intn(len(status))]
}
