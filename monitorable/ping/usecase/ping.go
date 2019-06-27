//+build !faker

package usecase

import (
	"context"

	. "github.com/monitoror/monitoror/models/tiles"
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

func (pu *pingUsecase) Ping(params *models.PingParams) (tile *HealthTile, err error) {
	tile = NewHealthTile(ping.PingTileType)
	tile.Label = params.Hostname

	ping, err := pu.repository.ExecutePing(context.Background(), params.Hostname)
	if err == nil {
		tile.Status = SuccessStatus
		tile.Message = ping.Average.String()
	} else {
		tile.Status = FailedStatus
		err = nil
	}

	return
}
