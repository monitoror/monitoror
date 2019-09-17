package http

import (
	. "github.com/monitoror/monitoror/models/tiles"
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
		HttpAny(params *models.HttpAnyParams) (*HealthTile, error)
		HttpRaw(params *models.HttpRawParams) (*HealthTile, error)
		HttpJson(params *models.HttpJsonParams) (*HealthTile, error)
		HttpYaml(params *models.HttpYamlParams) (*HealthTile, error)
	}
)
