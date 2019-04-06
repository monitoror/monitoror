//+build faker

package usecase

import (
	"fmt"
	"math/rand"
	"time"

	. "github.com/monitoror/monitoror/models/tiles"
	"github.com/monitoror/monitoror/monitorable/port"
	"github.com/monitoror/monitoror/monitorable/port/model"
)

type (
	portUsecase struct{}
)

// NewArticleUsecase will create new an articleUsecase object representation of article.Usecase interface
func NewPortUsecase() port.Usecase {
	return &portUsecase{}
}

func (pu *portUsecase) Port(params *model.PortParams) (tile *HealthTile, err error) {
	tile = NewHealthTile(port.PortTileSubType)
	tile.Label = fmt.Sprintf("%s:%d", params.Hostname, params.Port)

	// Init random generator
	rand.Seed(time.Now().UnixNano())

	// Status
	if params.Status != "" {
		tile.Status = params.Status
	} else {
		if rand.Intn(2) == 0 {
			tile.Status = SuccessStatus
		} else {
			tile.Status = FailStatus
		}
	}

	return
}
