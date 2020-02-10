//+build faker

package usecase

import (
	"math/rand"
	"time"

	"github.com/monitoror/monitoror/pkg/monitoror/faker"

	"github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorable/http"
	httpModels "github.com/monitoror/monitoror/monitorable/http/models"
	"github.com/monitoror/monitoror/pkg/monitoror/utils/nonempty"
)

type (
	httpUsecase struct {
		timeRefByUrl map[string]time.Time
	}
)

var availableStatuses = faker.Statuses{
	{models.SuccessStatus, time.Second * 30},
	{models.FailedStatus, time.Second * 30},
}

func NewHTTPUsecase() http.Usecase {
	return &httpUsecase{make(map[string]time.Time)}
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

	tile.Status = nonempty.Struct(params.GetStatus(), hu.computeStatus(url)).(models.TileStatus)
	if tile.Status == models.SuccessStatus && tileType != http.HTTPAnyTileType {
		if len(params.GetValues()) != 0 {
			tile.Values = params.GetValues()
		} else if params.GetMessage() != "" {
			tile.Message = params.GetMessage()
		} else {
			if rand.Intn(2) == 0 {
				tile.Values = []float64{1000}
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

func (hu *httpUsecase) computeStatus(url string) models.TileStatus {
	value, ok := hu.timeRefByUrl[url]
	if !ok {
		hu.timeRefByUrl[url] = faker.GetRefTime()
	}

	return faker.ComputeStatus(value, availableStatuses)
}
