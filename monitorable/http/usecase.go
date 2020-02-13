package http

import (
	"github.com/monitoror/monitoror/models"
	httpModels "github.com/monitoror/monitoror/monitorable/http/models"
)

const (
	HTTPStatusTileType    models.TileType = "HTTP-STATUS"
	HTTPRawTileType       models.TileType = "HTTP-RAW"
	HTTPFormattedTileType models.TileType = "HTTP-FORMATTED"
)

type (
	Usecase interface {
		HTTPStatus(params *httpModels.HTTPStatusParams) (*models.Tile, error)
		HTTPRaw(params *httpModels.HTTPRawParams) (*models.Tile, error)
		HTTPFormatted(params *httpModels.HTTPFormattedParams) (*models.Tile, error)
	}
)
