//+build faker

package usecase

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/monitoror/monitoror/internal/pkg/monitorable/faker"
	coreModels "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorables/ping/api"
	"github.com/monitoror/monitoror/monitorables/ping/api/models"
	"github.com/monitoror/monitoror/pkg/nonempty"
)

type (
	pingUsecase struct {
		timeRefByHostname map[string]time.Time
	}
)

var availableStatuses = faker.Statuses{
	{coreModels.SuccessStatus, time.Second * 30},
	{coreModels.FailedStatus, time.Second * 30},
}

func NewPingUsecase() api.Usecase {
	return &pingUsecase{make(map[string]time.Time)}
}

func (pu *pingUsecase) Ping(params *models.PingParams) (tile *coreModels.Tile, err error) {
	tile = coreModels.NewTile(api.PingTileType)
	tile.Label = params.Hostname

	// Code
	tile.Status = nonempty.Struct(params.Status, pu.computeStatus(params)).(coreModels.TileStatus)

	// Message
	if tile.Status == coreModels.SuccessStatus {
		tile.WithMetrics(coreModels.MillisecondUnit)
		if len(params.ValueValues) != 0 {
			tile.Metrics.Values = params.ValueValues
		} else {
			tile.Metrics.Values = append(tile.Metrics.Values, fmt.Sprintf("%d", rand.Int31n(300)))
		}
	}

	return
}

func (pu *pingUsecase) computeStatus(params *models.PingParams) coreModels.TileStatus {
	value, ok := pu.timeRefByHostname[params.Hostname]
	if !ok {
		pu.timeRefByHostname[params.Hostname] = faker.GetRefTime()
	}

	return faker.ComputeStatus(value, availableStatuses)
}
