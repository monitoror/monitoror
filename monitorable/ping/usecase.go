package ping

import (
	. "github.com/monitoror/monitoror/models/tiles"
	"github.com/monitoror/monitoror/monitorable/ping/model"
)

const (
	PingTileSubType TileSubType = "PING"
)

// Usecase represent the ping's usecases
type (
	Usecase interface {
		Ping(params *model.PingParams) (*HealthTile, error)
	}
)
