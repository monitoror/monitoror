package usecase

import (
	"github.com/jsdidierlaurent/monitowall/models/errors"
	. "github.com/jsdidierlaurent/monitowall/models/tiles"
	"github.com/jsdidierlaurent/monitowall/monitorable/ping"
	"github.com/jsdidierlaurent/monitowall/monitorable/ping/model"
	"github.com/jsdidierlaurent/monitowall/pkg/bind"
)

type (
	pingUsecase struct {
		repository ping.Repository
	}
)

const (
	PingTileSubType TileSubType = "PING"
)

// NewArticleUsecase will create new an articleUsecase object representation of article.Usecase interface
func NewPingUsecase(pr ping.Repository) ping.Usecase {
	return &pingUsecase{pr}
}

func (pu *pingUsecase) Ping(binder bind.Binder) (*HealthTile, error) {
	tile := NewHealthTile(PingTileSubType)

	// Bind / Validate Params
	params := &model.PingParams{}
	err := binder.Bind(params)
	if err != nil || !params.Validate() {
		return nil, errors.NewQueryParamsError(tile.Tile, err)
	}

	tile.Label = params.Hostname

	ping, err := pu.repository.Ping(params.Hostname)
	if err != nil {
		tile.Status = FailStatus
		return tile, nil
	}

	tile.Status = SuccessStatus
	tile.Message = ping.Average.String()

	return tile, nil
}
