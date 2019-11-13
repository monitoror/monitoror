//+build faker

package usecase

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorable/pingdom"
	pingdomModels "github.com/monitoror/monitoror/monitorable/pingdom/models"
	"github.com/monitoror/monitoror/pkg/monitoror/builder"
	"github.com/monitoror/monitoror/pkg/monitoror/utils/nonempty"
)

type (
	pingdomUsecase struct {
	}
)

var AvailableStatus = []models.TileStatus{models.SuccessStatus, models.FailedStatus, models.DisabledStatus}

func NewPingdomUsecase() pingdom.Usecase {
	return &pingdomUsecase{}
}

func (pu *pingdomUsecase) Check(params *pingdomModels.CheckParams) (tile *models.Tile, error error) {
	tile = models.NewTile(pingdom.PingdomCheckTileType)
	tile.Label = fmt.Sprintf("Check 1")

	// Init random generator
	rand.Seed(time.Now().UnixNano())

	// Code
	tile.Status = nonempty.Struct(params.Status, AvailableStatus[rand.Intn(len(AvailableStatus))]).(models.TileStatus)

	return
}

func (pu *pingdomUsecase) ListDynamicTile(params interface{}) ([]builder.Result, error) {
	panic("unimplemented")
}
