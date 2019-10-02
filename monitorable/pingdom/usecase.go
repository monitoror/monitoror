package pingdom

import (
	. "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorable/pingdom/models"
	"github.com/monitoror/monitoror/pkg/monitoror/builder"
)

const (
	PingdomCheckTileType  TileType = "PINGDOM-CHECK"
	PingdomChecksTileType TileType = "PINGDOM-CHECKS"
)

// Usecase represent the pingdom's usecases
type (
	Usecase interface {
		Check(params *models.CheckParams) (*Tile, error)
		ListDynamicTile(params interface{}) ([]builder.Result, error)
	}
)
