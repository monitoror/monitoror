//go:generate mockery --name Usecase

package api

import (
	uiConfigModels "github.com/monitoror/monitoror/api/config/models"
	coreModels "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorables/oneuptime/api/models"
)

const (
	OneuptimeMonitorTileType coreModels.TileType = "ONEUPTIME-MONITOR"
)

type (
	Usecase interface {
		Monitor(params *models.MonitorParams) (*coreModels.Tile, error)
		MonitorGenerator(params interface{}) ([]uiConfigModels.GeneratedTile, error)
	}
)
