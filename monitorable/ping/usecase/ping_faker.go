//+build faker

package usecase

import (
	"math/rand"
	"time"

	"github.com/monitoror/monitoror/pkg/monitoror/faker"

	"github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorable/ping"
	pingModels "github.com/monitoror/monitoror/monitorable/ping/models"
	"github.com/monitoror/monitoror/pkg/monitoror/utils/nonempty"
)

type (
	pingUsecase struct {
		timeRefByHostname map[string]time.Time
	}
)

var availableStatuses = faker.Statuses{
	{models.SuccessStatus, time.Second * 30},
	{models.FailedStatus, time.Second * 30},
}

func NewPingUsecase() ping.Usecase {
	return &pingUsecase{make(map[string]time.Time)}
}

func (pu *pingUsecase) Ping(params *pingModels.PingParams) (tile *models.Tile, err error) {
	tile = models.NewTile(ping.PingTileType)
	tile.Label = params.Hostname

	// Code
	tile.Status = nonempty.Struct(params.Status, pu.computeStatus(params)).(models.TileStatus)

	// Message
	if tile.Status == models.SuccessStatus {
		tile.Unit = models.MillisecondUnit

		if len(params.Values) != 0 {
			tile.Values = params.Values
		} else {
			tile.Values = []float64{300 * rand.Float64()}
		}
	}

	return
}

func (pu *pingUsecase) computeStatus(params *pingModels.PingParams) models.TileStatus {
	value, ok := pu.timeRefByHostname[params.Hostname]
	if !ok {
		pu.timeRefByHostname[params.Hostname] = faker.GetRefTime()
	}

	return faker.ComputeStatus(value, availableStatuses)
}
