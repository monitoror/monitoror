package config

import (
	"github.com/monitoror/monitoror/models/tiles"
	"github.com/monitoror/monitoror/monitorable/config/models"
	. "github.com/monitoror/monitoror/pkg/monitoror/builder"
	. "github.com/monitoror/monitoror/pkg/monitoror/validator"
)

// Usecase represent the config's usecases
type (
	Helper interface {
		RegisterTile(tileType tiles.TileType, validator Validator, path string)
		RegisterTileWithConfigVariant(tileType tiles.TileType, variant string, validator Validator, path string)

		RegisterDynamicTile(tileType tiles.TileType, validator Validator, builder DynamicTileBuilder)
		RegisterDynamicTileWithConfigVariant(tileType tiles.TileType, configVariant string, validator Validator, builder DynamicTileBuilder)
	}

	Usecase interface {
		Helper

		GetConfig(params *models.ConfigParams) (*models.Config, error)
		Verify(config *models.Config)
		Hydrate(config *models.Config, host string)
	}
)
