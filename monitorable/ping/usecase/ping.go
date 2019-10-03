//+build !faker

package usecase

import (
	. "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorable/ping"
	"github.com/monitoror/monitoror/monitorable/ping/models"
)

type (
	pingUsecase struct {
		repository ping.Repository
	}
)

func NewPingUsecase(repository ping.Repository) ping.Usecase {
	return &pingUsecase{repository}
}

func (pu *pingUsecase) Ping(params *models.PingParams) (tile *Tile, err error) {
	tile = NewTile(ping.PingTileType)
	tile.Label = params.Hostname

	ping, err := pu.repository.ExecutePing(params.Hostname)
	if err == nil {
		tile.Status = SuccessStatus
		tile.Unit = MillisecondUnit
		tile.Values = []float64{float64(ping.Average.Milliseconds())}
	} else {
		tile.Status = FailedStatus
		err = nil
	}

	return
}
