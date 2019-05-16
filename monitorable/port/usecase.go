package port

import (
	. "github.com/monitoror/monitoror/models/tiles"
	"github.com/monitoror/monitoror/monitorable/port/models"
)

const (
	PortTileType TileType = "PORT"
)

// Usecase represent the port's usecases
type (
	Usecase interface {
		Port(params *models.PortParams) (*HealthTile, error)
	}
)
