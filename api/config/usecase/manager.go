package usecase

import (
	"github.com/monitoror/monitoror/api/config"
	"github.com/monitoror/monitoror/api/config/models"
	coreModels "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/pkg/validator"
)

type (
	ConfigData struct {
		tileConfigs        map[coreModels.TileType]map[coreModels.VariantName]*TileConfig
		dynamicTileConfigs map[coreModels.TileType]map[coreModels.VariantName]*DynamicTileConfig
	}

	// TileConfig struct is used by GetConfig endpoint to check / hydrate config
	TileConfig struct {
		Validator       validator.SimpleValidator
		Path            string
		InitialMaxDelay int
	}

	// DynamicTileConfig struct is used by GetConfig endpoint to check / hydrate config
	DynamicTileConfig struct {
		Validator validator.SimpleValidator
		Builder   config.DynamicTileBuilder
	}
)

func initConfigData() *ConfigData {
	// TODO

	return &ConfigData{
		tileConfigs:        make(map[coreModels.TileType]map[coreModels.VariantName]*TileConfig),
		dynamicTileConfigs: make(map[coreModels.TileType]map[coreModels.VariantName]*DynamicTileConfig),
	}
}

func (cu *configUsecase) RegisterTile(tileType coreModels.TileType, variant []coreModels.VariantName, version models.RawVersion) {
	// TODO
}

func (cu *configUsecase) EnableTile(
	tileType coreModels.TileType, variant coreModels.VariantName, clientConfigValidator validator.SimpleValidator, path string, initialMaxDelay int,
) {
	value, exists := cu.configData.tileConfigs[tileType]
	if !exists {
		value = make(map[coreModels.VariantName]*TileConfig)
		cu.configData.tileConfigs[tileType] = value
	}

	value[variant] = &TileConfig{
		Path:            path,
		Validator:       clientConfigValidator,
		InitialMaxDelay: initialMaxDelay,
	}
}

func (cu *configUsecase) EnableDynamicTile(
	tileType coreModels.TileType, variant coreModels.VariantName, clientConfigValidator validator.SimpleValidator, builder config.DynamicTileBuilder,
) {
	// Used for authorized type
	cu.configData.tileConfigs[tileType] = nil

	value, exists := cu.configData.dynamicTileConfigs[tileType]
	if !exists {
		value = make(map[coreModels.VariantName]*DynamicTileConfig)
	}

	value[variant] = &DynamicTileConfig{
		Validator: clientConfigValidator,
		Builder:   builder,
	}
	cu.configData.dynamicTileConfigs[tileType] = value
}
