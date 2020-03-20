//+build faker

package usecase

import (
	"fmt"
	"time"

	"github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorables/pingdom/api"
	pingdomModels "github.com/monitoror/monitoror/monitorables/pingdom/api/models"
	"github.com/monitoror/monitoror/pkg/monitoror/builder"
	"github.com/monitoror/monitoror/pkg/monitoror/faker"
	"github.com/monitoror/monitoror/pkg/monitoror/utils/nonempty"
)

type (
	pingdomUsecase struct {
		timeRefByCheck map[int]time.Time
	}
)

var availableStatuses = faker.Statuses{
	{models.SuccessStatus, time.Second * 30},
	{models.FailedStatus, time.Second * 30},
	{models.DisabledStatus, time.Second * 10},
}

func NewPingdomUsecase() api.Usecase {
	return &pingdomUsecase{make(map[int]time.Time)}
}

func (pu *pingdomUsecase) Check(params *pingdomModels.CheckParams) (tile *models.Tile, error error) {
	tile = models.NewTile(api.PingdomCheckTileType)
	tile.Label = fmt.Sprintf(fmt.Sprintf("Check %d", *params.ID))

	// Code
	tile.Status = nonempty.Struct(params.Status, pu.computeStatus(params)).(models.TileStatus)

	return
}

func (pu *pingdomUsecase) Checks(params interface{}) ([]builder.Result, error) {
	panic("unimplemented")
}

func (pu *pingdomUsecase) computeStatus(params *pingdomModels.CheckParams) models.TileStatus {
	value, ok := pu.timeRefByCheck[*params.ID]
	if !ok {
		pu.timeRefByCheck[*params.ID] = faker.GetRefTime()
	}

	return faker.ComputeStatus(value, availableStatuses)
}
