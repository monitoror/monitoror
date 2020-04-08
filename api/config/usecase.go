//go:generate mockery -name TileSettingManager|TileEnabler|TileGeneratorEnabler|Usecase

package config

import (
	"github.com/monitoror/monitoror/api/config/models"
	coreModels "github.com/monitoror/monitoror/models"
)

type (
	TileSettingManager interface {
		Register(tileType coreModels.TileType, minimalVersion models.RawVersion, variants []coreModels.VariantName) TileEnabler
		RegisterGenerator(tileType coreModels.TileType, minimalVersion models.RawVersion, variants []coreModels.VariantName) TileGeneratorEnabler
	}

	TileEnabler interface {
		Enable(variant coreModels.VariantName, paramsValidator models.ParamsValidator, routePath string, initialMaxDelay int)
	}

	TileGeneratorEnabler interface {
		Enable(variant coreModels.VariantName, generatorParamsValidator models.ParamsValidator, tileGeneratorFunction models.TileGeneratorFunction)
	}

	Usecase interface {
		TileSettingManager

		GetConfig(params *models.ConfigParams) *models.ConfigBag
		Verify(config *models.ConfigBag)
		Hydrate(config *models.ConfigBag)
	}
)
