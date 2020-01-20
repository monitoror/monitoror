//+build faker

package usecase

import (
	"fmt"
	"time"

	"github.com/monitoror/monitoror/pkg/monitoror/faker"

	"github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorable/port"
	portModels "github.com/monitoror/monitoror/monitorable/port/models"
	"github.com/monitoror/monitoror/pkg/monitoror/utils/nonempty"
)

type (
	portUsecase struct {
		timeRefByHostnamePort map[string]time.Time
	}
)

var availableStatuses = faker.Statuses{
	{models.SuccessStatus, time.Second * 30},
	{models.FailedStatus, time.Second * 30},
}

// NewArticleUsecase will create new an articleUsecase object representation of article.Usecase interface
func NewPortUsecase() port.Usecase {
	return &portUsecase{make(map[string]time.Time)}
}

func (pu *portUsecase) Port(params *portModels.PortParams) (tile *models.Tile, err error) {
	tile = models.NewTile(port.PortTileType)
	tile.Label = fmt.Sprintf("%s:%d", params.Hostname, params.Port)

	// Code
	tile.Status = nonempty.Struct(params.Status, pu.computeStatus(params)).(models.TileStatus)

	return
}

func (pu *portUsecase) computeStatus(params *portModels.PortParams) models.TileStatus {
	key := fmt.Sprintf("%s:%d", params.Hostname, params.Port)
	value, ok := pu.timeRefByHostnamePort[key]
	if !ok {
		pu.timeRefByHostnamePort[key] = faker.GetRefTime()
	}

	return faker.ComputeStatus(value, availableStatuses)
}
