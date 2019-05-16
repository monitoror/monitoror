package port

import (
	. "github.com/monitoror/monitoror/models/tiles"
	"github.com/monitoror/monitoror/monitorable/port/model"
)

const (
	PortTileType TileType = "PORT"
)

// Usecase represent the port's usecases
type (
	Usecase interface {
		Port(params *model.PortParams) (*HealthTile, error)
	}
)
