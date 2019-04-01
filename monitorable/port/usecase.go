package port

import (
	. "github.com/jsdidierlaurent/monitoror/models/tiles"
	"github.com/jsdidierlaurent/monitoror/monitorable/port/model"
)

const (
	PortTileSubType TileSubType = "PORT"
)

// Usecase represent the port's usecases
type (
	Usecase interface {
		Port(params *model.PortParams) (*HealthTile, error)
	}
)
