package http

import (
	. "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorable/http/models"
)

const (
	HttpAnyTileType TileType = "HTTP-ANY"
	HttpRawTileType TileType = "HTTP-RAW"

	HttpJsonTileType TileType = "HTTP-JSON"
	HttpYamlTileType TileType = "HTTP-YAML"
)

// Usecase represent the ping's usecases
type (
	Usecase interface {
		HttpAny(params *models.HttpAnyParams) (*Tile, error)
		HttpRaw(params *models.HttpRawParams) (*Tile, error)
		HttpJson(params *models.HttpJsonParams) (*Tile, error)
		HttpYaml(params *models.HttpYamlParams) (*Tile, error)
	}
)
