//+build !faker

package usecase

import (
	"fmt"

	. "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorable/port"
	"github.com/monitoror/monitoror/monitorable/port/models"
)

type (
	portUsecase struct {
		repository port.Repository
	}
)

func NewPortUsecase(repository port.Repository) port.Usecase {
	return &portUsecase{repository}
}

func (pu *portUsecase) Port(params *models.PortParams) (tile *Tile, err error) {
	tile = NewTile(port.PortTileType)
	tile.Label = fmt.Sprintf("%s:%d", params.Hostname, params.Port)

	err = pu.repository.OpenSocket(params.Hostname, params.Port)
	if err == nil {
		tile.Status = SuccessStatus
	} else {
		tile.Status = FailedStatus
		err = nil
	}

	return
}
