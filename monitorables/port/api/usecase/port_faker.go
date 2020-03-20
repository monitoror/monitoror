//+build faker

package usecase

import (
	"fmt"
	"time"

	coreModels "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorables/port/api"
	"github.com/monitoror/monitoror/monitorables/port/api/models"
	"github.com/monitoror/monitoror/pkg/monitoror/faker"
	"github.com/monitoror/monitoror/pkg/monitoror/utils/nonempty"
)

type (
	portUsecase struct {
		timeRefByHostnamePort map[string]time.Time
	}
)

var availableStatuses = faker.Statuses{
	{coreModels.SuccessStatus, time.Second * 30},
	{coreModels.FailedStatus, time.Second * 30},
}

// NewArticleUsecase will create new an articleUsecase object representation of article.Usecase interface
func NewPortUsecase() api.Usecase {
	return &portUsecase{make(map[string]time.Time)}
}

func (pu *portUsecase) Port(params *models.PortParams) (tile *coreModels.Tile, err error) {
	tile = coreModels.NewTile(api.PortTileType)
	tile.Label = fmt.Sprintf("%s:%d", params.Hostname, params.Port)

	// Code
	tile.Status = nonempty.Struct(params.Status, pu.computeStatus(params)).(coreModels.TileStatus)

	return
}

func (pu *portUsecase) computeStatus(params *models.PortParams) coreModels.TileStatus {
	key := fmt.Sprintf("%s:%d", params.Hostname, params.Port)
	value, ok := pu.timeRefByHostnamePort[key]
	if !ok {
		pu.timeRefByHostnamePort[key] = faker.GetRefTime()
	}

	return faker.ComputeStatus(value, availableStatuses)
}
