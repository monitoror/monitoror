//+build faker

package usecase

import (
	"math/rand"
	"time"

	. "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorable/ping"
	"github.com/monitoror/monitoror/monitorable/ping/models"
	"github.com/monitoror/monitoror/pkg/monitoror/utils/nonempty"
)

type (
	pingUsecase struct {
	}
)

func NewPingUsecase() ping.Usecase {
	return &pingUsecase{}
}

func (pu *pingUsecase) Ping(params *models.PingParams) (tile *Tile, err error) {
	tile = NewTile(ping.PingTileType)
	tile.Label = params.Hostname

	// Init random generator
	rand.Seed(time.Now().UnixNano())

	// Code
	tile.Status = nonempty.Struct(params.Status, randomStatus()).(TileStatus)

	// Message
	if tile.Status == SuccessStatus {
		tile.Unit = MillisecondUnit
		tile.Values = []float64{nonempty.Float64(params.Value, float64(time.Duration(rand.Intn(10000))))}
	}

	return
}

func randomStatus() TileStatus {
	if rand.Intn(2) == 0 {
		return SuccessStatus
	} else {
		return FailedStatus
	}
}
