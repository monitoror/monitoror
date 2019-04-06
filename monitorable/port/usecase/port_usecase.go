//+build !faker

package usecase

import (
	"fmt"

	. "github.com/monitoror/monitoror/models/tiles"
	"github.com/monitoror/monitoror/monitorable/port"
	"github.com/monitoror/monitoror/monitorable/port/model"
)

type (
	portUsecase struct {
		repository port.Repository
	}
)

func NewPortUsecase(pr port.Repository) port.Usecase {
	return &portUsecase{pr}
}

func (pu *portUsecase) Port(params *model.PortParams) (tile *HealthTile, err error) {
	tile = NewHealthTile(port.PortTileSubType)
	tile.Label = fmt.Sprintf("%s:%d", params.Hostname, params.Port)

	err = pu.repository.CheckPort(params.Hostname, params.Port)
	if err == nil {
		tile.Status = SuccessStatus
	} else {
		tile.Status = FailStatus
		err = nil
	}

	return
}
