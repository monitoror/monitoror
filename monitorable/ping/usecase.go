package ping

import (
	. "github.com/jsdidierlaurent/monitowall/models/tiles"
	"github.com/jsdidierlaurent/monitowall/pkg/bind"
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
