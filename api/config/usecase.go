package config

import (
	configModels "github.com/monitoror/monitoror/api/config/models"
	"github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/pkg/monitoror/builder"
	"github.com/monitoror/monitoror/pkg/monitoror/utils"
)

type (
	Manager interface {
		RegisterTile(tileType models.TileType, variants []models.Variant, version string)
		EnableTile(tileType models.TileType, variant models.Variant, validator utils.Validator, path string, initialMaxDelay int)
		EnableDynamicTile(tileType models.TileType, variant models.Variant, Validator utils.Validator, builder builder.DynamicTileBuilder)
	}

	Usecase interface {
		Manager

		GetConfig(params *configModels.ConfigParams) *configModels.ConfigBag
		Verify(config *configModels.ConfigBag)
		Hydrate(config *configModels.ConfigBag)
	}
)
