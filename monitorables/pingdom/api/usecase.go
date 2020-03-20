package api

import (
	coreModels "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorables/pingdom/api/models"
	"github.com/monitoror/monitoror/pkg/monitoror/builder"
)

const (
	PingdomCheckTileType  coreModels.TileType = "PINGDOM-CHECK"
	PingdomChecksTileType coreModels.TileType = "PINGDOM-CHECKS"
)

type (
	Usecase interface {
		Check(params *models.CheckParams) (*coreModels.Tile, error)
		Checks(params interface{}) ([]builder.Result, error)
	}
)
