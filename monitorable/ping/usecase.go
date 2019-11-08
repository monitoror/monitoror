package ping

import (
	"github.com/monitoror/monitoror/models"
	pingModels "github.com/monitoror/monitoror/monitorable/ping/models"
)

const (
	PingTileType models.TileType = "PING"
)

type (
	Usecase interface {
		Ping(params *pingModels.PingParams) (*models.Tile, error)
	}
)
