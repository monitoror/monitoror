//go:generate mockery -name Usecase

package api

import (
	uiConfigModels "github.com/monitoror/monitoror/api/config/models"
	coreModels "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorables/pingdom/api/models"
)

const (
	PingdomCheckTileType coreModels.TileType = "PINGDOM-CHECK"
)

type (
	Usecase interface {
		Check(params *models.CheckParams) (*coreModels.Tile, error)
		CheckGenerator(params interface{}) ([]uiConfigModels.GeneratedTile, error)
	}
)
