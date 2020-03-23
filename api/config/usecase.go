package config

import (
	"github.com/monitoror/monitoror/api/config/models"
	coreModels "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/pkg/validator"
)

type (
	Manager interface {
		RegisterTile(tileType coreModels.TileType, variants []coreModels.Variant, version string)
		EnableTile(tileType coreModels.TileType, variant coreModels.Variant, validator validator.SimpleValidator, path string, initialMaxDelay int)
		EnableDynamicTile(tileType coreModels.TileType, variant coreModels.Variant, Validator validator.SimpleValidator, builder DynamicTileBuilder)
	}
	DynamicTileBuilder func(params interface{}) ([]models.DynamicTileResult, error)

	Usecase interface {
		Manager

		GetConfig(params *models.ConfigParams) *models.ConfigBag
		Verify(config *models.ConfigBag)
		Hydrate(config *models.ConfigBag)
	}
)
