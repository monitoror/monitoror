//+build faker

package usecase

import (
	"math/rand"
	"strconv"
	"time"

	"github.com/monitoror/monitoror/internal/pkg/monitorable/faker"
	coreModels "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorables/http/api"
	"github.com/monitoror/monitoror/monitorables/http/api/models"
	"github.com/monitoror/monitoror/pkg/nonempty"
)

type (
	httpUsecase struct {
		timeRefByUrl map[string]time.Time
	}
)

var availableStatuses = faker.Statuses{
	{coreModels.SuccessStatus, time.Second * 30},
	{coreModels.FailedStatus, time.Second * 30},
}

func NewHTTPUsecase() api.Usecase {
	return &httpUsecase{make(map[string]time.Time)}
}

// HTTPStatus only check status code
func (hu *httpUsecase) HTTPStatus(params *models.HTTPStatusParams) (tile *coreModels.Tile, err error) {
	return hu.httpAll(api.HTTPStatusTileType, params.URL, params)
}

// HTTPRaw check status code and content
func (hu *httpUsecase) HTTPRaw(params *models.HTTPRawParams) (tile *coreModels.Tile, err error) {
	return hu.httpAll(api.HTTPRawTileType, params.URL, params)
}

func (hu *httpUsecase) HTTPFormatted(params *models.HTTPFormattedParams) (tile *coreModels.Tile, err error) {
	return hu.httpAll(api.HTTPFormattedTileType, params.URL, params)
}

// httpAll handle all http usecase by checking if params match interfaces listed in coreModels.params
func (hu *httpUsecase) httpAll(tileType coreModels.TileType, url string, params models.FakerParamsProvider) (tile *coreModels.Tile, err error) {
	tile = coreModels.NewTile(tileType)
	tile.Label = url

	tile.Status = nonempty.Struct(params.GetStatus(), hu.computeStatus(url)).(coreModels.TileStatus)
	if tile.Status == coreModels.SuccessStatus && tileType != api.HTTPStatusTileType {
		var values []string
		if len(params.GetValueValues()) != 0 {
			values = params.GetValueValues()
		}

		if len(values) == 0 {
			if rand.Intn(2) == 0 {
				values = append(values, "1000")
			} else {
				values = append(values, "random message")
			}
		}

		if params.GetValueUnit() == "" {
			if _, err := strconv.ParseFloat(values[0], 64); err == nil {
				tile.WithMetrics(coreModels.NumberUnit)
			} else {
				tile.WithMetrics(coreModels.RawUnit)
			}
		} else {
			tile.WithMetrics(params.GetValueUnit())
		}

		tile.Metrics.Values = values
	}

	if tile.Status == coreModels.FailedStatus {
		tile.Message = nonempty.String(params.GetMessage(), "Fake error message")
	}

	return
}

func (hu *httpUsecase) computeStatus(url string) coreModels.TileStatus {
	value, ok := hu.timeRefByUrl[url]
	if !ok {
		hu.timeRefByUrl[url] = faker.GetRefTime()
	}

	return faker.ComputeStatus(value, availableStatuses)
}
