//+build !faker

package usecase

import (
	"fmt"

	"github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorable/port"
	portModels "github.com/monitoror/monitoror/monitorable/port/models"
)

type (
	portUsecase struct {
		repository port.Repository
	}
)

func NewPortUsecase(repository port.Repository) port.Usecase {
	return &portUsecase{repository}
}

func (pu *portUsecase) Port(params *portModels.PortParams) (tile *models.Tile, err error) {
	tile = models.NewTile(port.PortTileType)
	tile.Label = fmt.Sprintf("%s:%d", params.Hostname, params.Port)

	err = pu.repository.OpenSocket(params.Hostname, params.Port)
	if err == nil {
		tile.Status = models.SuccessStatus
	} else {
		tile.Status = models.FailedStatus
		err = nil
	}

	return
}
