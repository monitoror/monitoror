//+build faker

package usecase

import (
	"math/rand"
	"strconv"
	"time"

	"github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorable/http"
	httpModels "github.com/monitoror/monitoror/monitorable/http/models"
	"github.com/monitoror/monitoror/pkg/monitoror/faker"
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

// HTTPStatus only check status code
func (hu *httpUsecase) HTTPStatus(params *httpModels.HTTPStatusParams) (tile *models.Tile, err error) {
	return hu.httpAll(http.HTTPStatusTileType, params.URL, params)
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
	if tile.Status == models.SuccessStatus && tileType != http.HTTPStatusTileType {
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
				tile.WithValue(models.NumberUnit)
			} else {
				tile.WithValue(models.RawUnit)
			}
		} else {
			tile.WithValue(params.GetValueUnit())
		}

		tile.Value.Values = values
	}

	if tile.Status == models.FailedStatus {
		tile.Message = nonempty.String(params.GetMessage(), "Fake error message")
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
