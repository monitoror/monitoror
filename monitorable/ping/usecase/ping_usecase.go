//+build !faker

package usecase

import (
	. "github.com/jsdidierlaurent/monitoror/models/tiles"
	"github.com/jsdidierlaurent/monitoror/monitorable/ping"
	"github.com/jsdidierlaurent/monitoror/monitorable/ping/model"
)

type (
	pingUsecase struct {
		repository ping.Repository
	}
)

func NewPingUsecase(pr ping.Repository) ping.Usecase {
	return &pingUsecase{pr}
}

func (pu *pingUsecase) Ping(params *model.PingParams) (tile *HealthTile, err error) {
	tile = NewHealthTile(ping.PingTileSubType)
	tile.Label = params.Hostname

	ping, err := pu.repository.CheckPing(params.Hostname)
	if err == nil {
		tile.Status = SuccessStatus
		tile.Message = ping.Average.String()
	} else {
		tile.Status = FailStatus
		err = nil
	}

	return
}
