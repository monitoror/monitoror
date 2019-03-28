//+build faker

package usecase

import (
	"math/rand"
	"time"

	"github.com/jsdidierlaurent/monitowall/models/errors"
	. "github.com/jsdidierlaurent/monitowall/models/tiles"
	"github.com/jsdidierlaurent/monitowall/monitorable/ping"
	"github.com/jsdidierlaurent/monitowall/monitorable/ping/model"
	"github.com/jsdidierlaurent/monitowall/pkg/bind"
)

type (
	pingUsecase struct {
	}
)

// NewArticleUsecase will create new an articleUsecase object representation of article.Usecase interface
func NewPingUsecase() ping.Usecase {
	return &pingUsecase{}
}

func (pu *pingUsecase) Ping(binder bind.Binder) (*HealthTile, error) {
	tile := NewHealthTile(ping.PingTileSubType)

	// Bind / Validate Params
	params := &model.PingParams{}
	err := binder.Bind(params)
	if err != nil || !params.Validate() {
		return nil, errors.NewQueryParamsError(tile.Tile, err)
	}

	tile.Label = params.Hostname

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

	// Message
	if tile.Status == SuccessStatus {
		if params.Message != "" {
			tile.Message = params.Message
		} else {
			tile.Message = (time.Duration(rand.Intn(10000)) * time.Millisecond).String()
		}
	}

	return tile, nil
}
