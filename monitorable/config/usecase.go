package config

import (
	"github.com/monitoror/monitoror/models"
	configModels "github.com/monitoror/monitoror/monitorable/config/models"
	"github.com/monitoror/monitoror/pkg/monitoror/builder"
	"github.com/monitoror/monitoror/pkg/monitoror/utils"
)

type (
	Helper interface {
		RegisterTile(tileType models.TileType, validator utils.Validator, path string)
		RegisterTileWithConfigVariant(tileType models.TileType, variant string, validator utils.Validator, path string)

		RegisterDynamicTile(tileType models.TileType, validator utils.Validator, builder builder.DynamicTileBuilder)
		RegisterDynamicTileWithConfigVariant(tileType models.TileType, configVariant string, validator utils.Validator, builder builder.DynamicTileBuilder)
	}

	Usecase interface {
		Helper

		GetConfig(params *configModels.ConfigParams) (*configModels.Config, error)
		Verify(config *configModels.Config)
		Hydrate(config *configModels.Config)
	}
)
