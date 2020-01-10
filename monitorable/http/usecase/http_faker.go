//+build faker

package usecase

import (
	"math/rand"
	"time"

	"github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorable/http"
	httpModels "github.com/monitoror/monitoror/monitorable/http/models"
	"github.com/monitoror/monitoror/pkg/monitoror/utils/nonempty"
)

type (
	httpUsecase struct {
	}
)

func NewHTTPUsecase() http.Usecase {
	return &httpUsecase{}
}

// HTTPAny only check status code
func (hu *httpUsecase) HTTPAny(params *httpModels.HTTPAnyParams) (tile *models.Tile, err error) {
	return hu.httpAll(http.HTTPAnyTileType, params.URL, params)
}

// HTTPRaw check status code and content
func (hu *httpUsecase) HTTPRaw(params *httpModels.HTTPRawParams) (tile *models.Tile, err error) {
	return hu.httpAll(http.HTTPRawTileType, params.URL, params)
}

func (hu *httpUsecase) HTTPFormatted(params *httpModels.HTTPFormattedParams) (tile *models.Tile, err error) {
	return hu.httpAll(http.HTTPFormattedTileType, params.URL, params)
}

// httpAll handle all http usecase by checking if params match interfaces listed in models.params
func (hu *httpUsecase) httpAll(tileType models.TileType, url string, params httpModels.FakerParamsProvider) (tile *models.Tile, err error) {
	tile = models.NewTile(tileType)
	tile.Label = url

	// Init random generator
	rand.Seed(time.Now().UnixNano())

	tile.Status = nonempty.Struct(params.GetStatus(), randomStatus()).(models.TileStatus)
	if tile.Status == models.SuccessStatus && tileType != http.HTTPAnyTileType {
		if len(params.GetValues()) != 0 {
			tile.Values = params.GetValues()
		} else if params.GetMessage() != "" {
			tile.Message = params.GetMessage()
		} else {
			if rand.Intn(2) == 0 {
				tile.Values = []float64{float64(rand.Intn(10000))}
			} else {
				tile.Message = "random message"
			}
		}

	}

	if tile.Status == models.FailedStatus {
		tile.Message = nonempty.String(params.GetMessage(), "random error message")
	}

	return
}

func randomStatus() models.TileStatus {
	if rand.Intn(2) == 0 {
		return models.SuccessStatus
	} else {
		return models.FailedStatus
	}
}
