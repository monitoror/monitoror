package port

import (
	"github.com/monitoror/monitoror/models"
	portModels "github.com/monitoror/monitoror/monitorable/port/models"
)

const (
	PortTileType models.TileType = "PORT"
)

type (
	Usecase interface {
		Port(params *portModels.PortParams) (*models.Tile, error)
	}
)
