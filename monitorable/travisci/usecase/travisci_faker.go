//+build faker

package usecase

import (
	"fmt"
	"math/rand"
	"time"

	. "github.com/AlekSi/pointer"

	"github.com/monitoror/monitoror/pkg/monitoror/utils/nonempty"

	. "github.com/monitoror/monitoror/models/tiles"
	"github.com/monitoror/monitoror/monitorable/travisci"
	"github.com/monitoror/monitoror/monitorable/travisci/models"
)

var AvailableStatus = []TileStatus{SuccessStatus, FailedStatus, RunningStatus, QueuedStatus, WarningStatus}
var AvailablePreviousStatus = []TileStatus{SuccessStatus, FailedStatus}

type (
	travisCIUsecase struct {
		cachedRunningValue map[string]*durations
	}

	durations struct {
		duration          int64
		estimatedDuration int64
	}
)

func NewTravisCIUsecase() travisci.Usecase {
	return &travisCIUsecase{make(map[string]*durations)}
}

func (tu *travisCIUsecase) Build(params *models.BuildParams) (tile *BuildTile, err error) {
	tile = NewBuildTile(travisci.TravisCIBuildTileType)
	tile.Label = fmt.Sprintf("%s : #%s", params.Repository, params.Branch)

	// Init random generator
	rand.Seed(time.Now().UnixNano())

	tile.Status = nonempty.Struct(params.Status, randomStatus(AvailableStatus)).(TileStatus)

	if tile.Status == WarningStatus {
		tile.Message = "random error message"
		return
	}

	tile.Author = &Author{}
	tile.Author.Name = nonempty.String(params.AuthorName, "Faker")
	tile.Author.AvatarUrl = nonempty.String(params.AuthorAvatarUrl, "https://www.gravatar.com/avatar/00000000000000000000000000000000")

	if tile.Status == SuccessStatus || tile.Status == FailedStatus {
		min := time.Now().Unix() - int64(time.Hour.Seconds()*24*30) - 3600
		max := time.Now().Unix() - 3600
		delta := max - min

		tile.StartedAt = ToInt64(nonempty.Int64(params.StartedAt, time.Unix(rand.Int63n(delta)+min, 0).Unix()))
		tile.FinishedAt = ToInt64(nonempty.Int64(params.FinishedAt, *tile.StartedAt+rand.Int63n(3600)))
	}

	if tile.Status == QueuedStatus || tile.Status == RunningStatus {
		tile.StartedAt = ToInt64(nonempty.Int64(params.StartedAt, time.Now().Unix()-rand.Int63n(3600)))
		tile.PreviousStatus = nonempty.Struct(params.PreviousStatus, randomStatus(AvailablePreviousStatus)).(TileStatus)
	}

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
		tile.EstimatedDuration = ToInt64(dur.estimatedDuration)
	}

	return
}

func randomStatus(status []TileStatus) TileStatus {
	return status[rand.Intn(len(status))]
}
