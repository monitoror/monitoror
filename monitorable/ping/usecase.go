package ping

import (
	. "github.com/monitoror/monitoror/models/tiles"
	"github.com/monitoror/monitoror/monitorable/ping/models"
)

const (
	PingTileType TileType = "PING"
)

// Usecase represent the ping's usecases
type (
	Usecase interface {
		Ping(params *models.PingParams) (*HealthTile, error)
	}
)
