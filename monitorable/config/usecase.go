package config

import (
	. "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorable/config/models"
	. "github.com/monitoror/monitoror/pkg/monitoror/builder"
	. "github.com/monitoror/monitoror/pkg/monitoror/utils"
)

// Usecase represent the config's usecases
type (
	Helper interface {
		RegisterTile(tileType TileType, validator Validator, path string)
		RegisterTileWithConfigVariant(tileType TileType, variant string, validator Validator, path string)

		RegisterDynamicTile(tileType TileType, validator Validator, builder DynamicTileBuilder)
		RegisterDynamicTileWithConfigVariant(tileType TileType, configVariant string, validator Validator, builder DynamicTileBuilder)
	}

	Usecase interface {
		Helper

		GetConfig(params *models.ConfigParams) (*models.Config, error)
		Verify(config *models.Config)
		Hydrate(config *models.Config)
	}
)
