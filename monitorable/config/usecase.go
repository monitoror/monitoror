package config

import (
	"github.com/monitoror/monitoror/models"
	configModels "github.com/monitoror/monitoror/monitorable/config/models"
	"github.com/monitoror/monitoror/pkg/monitoror/builder"
	"github.com/monitoror/monitoror/pkg/monitoror/utils"
)

type (
	Manager interface {
		RegisterTile(tileType models.TileType, variant string, clientConfigValidator utils.Validator, path string, initialMaxDelay int)
		RegisterDynamicTile(tileType models.TileType, variant string, clientConfigValidator utils.Validator, builder builder.DynamicTileBuilder)
		DisableTile(tileType models.TileType, variant string)
	}

	Usecase interface {
		Manager

		GetConfig(params *configModels.ConfigParams) *configModels.ConfigBag
		Verify(config *configModels.ConfigBag)
		Hydrate(config *configModels.ConfigBag)
	}
)
