//+build faker

package usecase

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/monitoror/monitoror/pkg/monitoror/utils/nonempty"

	. "github.com/monitoror/monitoror/models/tiles"
	"github.com/monitoror/monitoror/monitorable/pingdom"
	"github.com/monitoror/monitoror/monitorable/pingdom/models"
	"github.com/monitoror/monitoror/pkg/monitoror/builder"
)

type (
	pingdomUsecase struct {
	}
)

var AvailableStatus = []TileStatus{SuccessStatus, FailedStatus, DisabledStatus}

func NewPingdomUsecase() pingdom.Usecase {
	return &pingdomUsecase{}
}

func (pu *pingdomUsecase) Check(params *models.CheckParams) (tile *HealthTile, error error) {
	tile = NewHealthTile(pingdom.PingdomCheckTileType)
	tile.Label = fmt.Sprintf("Check 1")

	// Init random generator
	rand.Seed(time.Now().UnixNano())

	// Code
	tile.Status = nonempty.Struct(params.Status, AvailableStatus[rand.Intn(len(AvailableStatus))]).(TileStatus)

	return
}

func (pu *pingdomUsecase) ListDynamicTile(params interface{}) ([]builder.Result, error) {
	panic("unimplemented")
}
