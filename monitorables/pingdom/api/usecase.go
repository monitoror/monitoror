//go:generate mockery -name Usecase

package api

import (
	models2 "github.com/monitoror/monitoror/api/config/models"
	coreModels "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorables/pingdom/api/models"
)

const (
	PingdomCheckTileType  coreModels.TileType = "PINGDOM-CHECK"
	PingdomChecksTileType coreModels.TileType = "PINGDOM-CHECKS"
)

type (
	Usecase interface {
		Check(params *models.CheckParams) (*coreModels.Tile, error)
		Checks(params interface{}) ([]models2.DynamicTileResult, error)
	}
)
