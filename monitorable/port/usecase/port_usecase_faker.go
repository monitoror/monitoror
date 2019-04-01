//+build faker

package usecase

import (
	"fmt"
	"math/rand"
	"time"

	. "github.com/jsdidierlaurent/monitoror/models/tiles"
	"github.com/jsdidierlaurent/monitoror/monitorable/ping/model"
	"github.com/jsdidierlaurent/monitoror/monitorable/port"
)

type (
	portUsecase struct {
		repository port.Repository
	}
)

// NewArticleUsecase will create new an articleUsecase object representation of article.Usecase interface
func NewPortUsecase(pr port.Repository) port.Usecase {
	return &portUsecase{pr}
}

func (pu *portUsecase) CheckPort(params *model.PortParams) (tile *HealthTile, err error) {
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
