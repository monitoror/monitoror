package config

import (
	"github.com/monitoror/monitoror/models"
	configModels "github.com/monitoror/monitoror/monitorable/config/models"
	"github.com/monitoror/monitoror/pkg/monitoror/builder"
	"github.com/monitoror/monitoror/pkg/monitoror/utils"
)

type (
	Helper interface {
		RegisterTile(tileType models.TileType, clientConfigValidator utils.Validator, path string, initialMaxDelay int)
		RegisterTileWithConfigVariant(tileType models.TileType, variant string, clientConfigValidator utils.Validator, path string, initialMaxDelay int)

		RegisterDynamicTile(tileType models.TileType, clientConfigValidator utils.Validator, builder builder.DynamicTileBuilder)
		RegisterDynamicTileWithConfigVariant(tileType models.TileType, configVariant string, clientConfigValidator utils.Validator, builder builder.DynamicTileBuilder)
	}

	Usecase interface {
		Helper

		GetConfig(params *configModels.ConfigParams) (*configModels.Config, error)
		Verify(config *configModels.Config)
		Hydrate(config *configModels.Config)
	}
)
