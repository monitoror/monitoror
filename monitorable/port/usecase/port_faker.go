//+build faker

package usecase

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorable/port"
	portModels "github.com/monitoror/monitoror/monitorable/port/models"
	"github.com/monitoror/monitoror/pkg/monitoror/utils/nonempty"
)

type (
	portUsecase struct{}
)

// NewArticleUsecase will create new an articleUsecase object representation of article.Usecase interface
func NewPortUsecase() port.Usecase {
	return &portUsecase{}
}

func (pu *portUsecase) Port(params *portModels.PortParams) (tile *models.Tile, err error) {
	tile = models.NewTile(port.PortTileType)
	tile.Label = fmt.Sprintf("%s:%d", params.Hostname, params.Port)

	// Init random generator
	rand.Seed(time.Now().UnixNano())

	// Code
	tile.Status = nonempty.Struct(params.Status, randomStatus()).(models.TileStatus)

	return
}

func randomStatus() models.TileStatus {
	if rand.Intn(2) == 0 {
		return models.SuccessStatus
	} else {
		return models.FailedStatus
	}
}
