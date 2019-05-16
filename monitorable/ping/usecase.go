package ping

import (
	. "github.com/monitoror/monitoror/models/tiles"
	"github.com/monitoror/monitoror/monitorable/ping/model"
)

const (
	PingTileType TileType = "PING"
)

// Usecase represent the ping's usecases
type (
	Usecase interface {
		Ping(params *model.PingParams) (*HealthTile, error)
	}
)
