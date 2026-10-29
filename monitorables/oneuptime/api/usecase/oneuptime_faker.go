//go:build faker

package usecase

import (
	"time"

	uiConfigModels "github.com/monitoror/monitoror/api/config/models"
	"github.com/monitoror/monitoror/internal/pkg/monitorable/faker"
	coreModels "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorables/oneuptime/api"
	"github.com/monitoror/monitoror/monitorables/oneuptime/api/models"
	"github.com/monitoror/monitoror/pkg/nonempty"
)

type (
	oneuptimeUsecase struct {
		timeRefByMonitor map[string]time.Time
	}
)

var availableStatuses = faker.Statuses{
	{coreModels.SuccessStatus, time.Second * 30},
	{coreModels.FailedStatus, time.Second * 30},
	{coreModels.DisabledStatus, time.Second * 10},
	{coreModels.WarningStatus, time.Second * 10},
}

func NewOneuptimeUsecase() api.Usecase {
	return &oneuptimeUsecase{make(map[string]time.Time)}
}

func (ou *oneuptimeUsecase) Monitor(params *models.MonitorParams) (*coreModels.Tile, error) {
	tile := coreModels.NewTile(api.OneuptimeMonitorTileType)
	tile.Label = "Monitor " + params.ID
	tile.Status = nonempty.Struct(params.Status, ou.computeStatus(params.ID)).(coreModels.TileStatus)

	return tile, nil
}

func (ou *oneuptimeUsecase) MonitorGenerator(params interface{}) ([]uiConfigModels.GeneratedTile, error) {
	panic("unimplemented")
}

func (ou *oneuptimeUsecase) computeStatus(id string) coreModels.TileStatus {
	value, ok := ou.timeRefByMonitor[id]
	if !ok {
		value = faker.GetRefTime()
		ou.timeRefByMonitor[id] = value
	}

	return faker.ComputeStatus(value, availableStatuses)
}
