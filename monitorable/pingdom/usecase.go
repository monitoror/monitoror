package pingdom

import (
	"github.com/monitoror/monitoror/models"
	pingdomModels "github.com/monitoror/monitoror/monitorable/pingdom/models"
	"github.com/monitoror/monitoror/pkg/monitoror/builder"
)

const (
	PingdomCheckTileType  models.TileType = "PINGDOM-CHECK"
	PingdomChecksTileType models.TileType = "PINGDOM-CHECKS"
)

type (
	Usecase interface {
		Check(params *pingdomModels.CheckParams) (*models.Tile, error)
		Checks(params interface{}) ([]builder.Result, error)
	}
)
