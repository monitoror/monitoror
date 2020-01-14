package http

import (
	"github.com/monitoror/monitoror/models"
	httpModels "github.com/monitoror/monitoror/monitorable/http/models"
)

const (
	HTTPAnyTileType       models.TileType = "HTTP-ANY"
	HTTPRawTileType       models.TileType = "HTTP-RAW"
	HTTPFormattedTileType models.TileType = "HTTP-FORMATTED"
)

type (
	Usecase interface {
		HTTPAny(params *httpModels.HTTPAnyParams) (*models.Tile, error)
		HTTPRaw(params *httpModels.HTTPRawParams) (*models.Tile, error)
		HTTPFormatted(params *httpModels.HTTPFormattedParams) (*models.Tile, error)
	}
)
