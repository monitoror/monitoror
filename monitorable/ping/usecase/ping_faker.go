//+build faker

package usecase

import (
	"math/rand"
	"time"

	"github.com/monitoror/monitoror/pkg/monitoror/utils/nonempty"

	. "github.com/monitoror/monitoror/models/tiles"
	"github.com/monitoror/monitoror/monitorable/ping"
	"github.com/monitoror/monitoror/monitorable/ping/models"
)

type (
	pingUsecase struct {
	}
)

func NewPingUsecase() ping.Usecase {
	return &pingUsecase{}
}

func (pu *pingUsecase) Ping(params *models.PingParams) (tile *HealthTile, err error) {
	tile = NewHealthTile(ping.PingTileType)
	tile.Label = params.Hostname

	// Init random generator
	rand.Seed(time.Now().UnixNano())

	// Code
	tile.Status = nonempty.Struct(params.Status, randomStatus()).(TileStatus)

	// Message
	if tile.Status == SuccessStatus {
		tile.Message = nonempty.String(params.Message, (time.Duration(rand.Intn(10000)) * time.Millisecond).String())
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
