//+build faker

package usecase

import (
	"math/rand"
	"time"

	. "github.com/monitoror/monitoror/models/tiles"
	"github.com/monitoror/monitoror/monitorable/ping"
	"github.com/monitoror/monitoror/monitorable/ping/model"
)

type (
	pingUsecase struct {
	}
)

func NewPingUsecase() ping.Usecase {
	return &pingUsecase{}
}

func (pu *pingUsecase) Ping(params *model.PingParams) (tile *HealthTile, err error) {
	tile = NewHealthTile(ping.PingTileSubType)
	tile.Label = params.Hostname

	// Init random generator
	rand.Seed(time.Now().UnixNano())

	// Status
	if params.Status != "" {
		tile.Status = params.Status
	} else {
		if rand.Intn(2) == 0 {
			tile.Status = SuccessStatus
		} else {
			tile.Status = FailStatus
		}
	}

	// Message
	if tile.Status == SuccessStatus {
		if params.Message != "" {
			tile.Message = params.Message
		} else {
			tile.Message = (time.Duration(rand.Intn(10000)) * time.Millisecond).String()
		}
	}

	return
}
