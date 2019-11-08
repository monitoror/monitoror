package http

import (
	"github.com/monitoror/monitoror/models"
	httpModels "github.com/monitoror/monitoror/monitorable/http/models"
)

const (
	HTTPAnyTileType models.TileType = "HTTP-ANY"
	HTTPRawTileType models.TileType = "HTTP-RAW"

	HTTPJsonTileType models.TileType = "HTTP-JSON"
	HTTPYamlTileType models.TileType = "HTTP-YAML"
)

type (
	Usecase interface {
		HTTPAny(params *httpModels.HTTPAnyParams) (*models.Tile, error)
		HTTPRaw(params *httpModels.HTTPRawParams) (*models.Tile, error)
		HTTPJson(params *httpModels.HTTPJsonParams) (*models.Tile, error)
		HTTPYaml(params *httpModels.HTTPYamlParams) (*models.Tile, error)
	}
)
