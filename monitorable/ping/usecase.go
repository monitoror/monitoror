package ping

import (
	. "github.com/jsdidierlaurent/monitoror/models/tiles"
	"github.com/jsdidierlaurent/monitoror/pkg/bind"
)

const (
	PingTileSubType TileSubType = "PING"
)

// Usecase represent the ping's usecases
type (
	Usecase interface {
		Ping(binder bind.Binder) (*HealthTile, error)
	}
)
