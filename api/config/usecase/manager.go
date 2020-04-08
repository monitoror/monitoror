package usecase

import (
	"fmt"

	"github.com/monitoror/monitoror/api/config"
	"github.com/monitoror/monitoror/api/config/models"
	coreModels "github.com/monitoror/monitoror/models"
)

type (
	ConfigData struct {
		tileSettings          map[coreModels.TileType]*TileSetting
		tileGeneratorSettings map[coreModels.TileType]*TileGeneratorSetting
	}

	Variant interface {
		IsEnabled() bool
		GetValidator() models.ParamsValidator
	}

	TileSetting struct {
		//TileType
		TileType coreModels.TileType
		//MinimalVersion is the version that makes the tile available
		MinimalVersion models.RawVersion
		//Variants list all registered variants (can be available or not)
		Variants map[coreModels.VariantName]*VariantSettings
	}

	VariantSettings struct {
		//Enabled define if variant is enabled
		Enabled bool

		//VariantName
		VariantName coreModels.VariantName

		//Route path of the api endpoint for this tile. Used by hydrate
		RoutePath *string
		//InitialMaxDelay is forward to ui to delay initial requests
		InitialMaxDelay *int
		//ParamsValidator is used to validate given params
		ParamsValidator models.ParamsValidator
	}

	TileGeneratorSetting struct {
		//TileType
		TileType coreModels.TileType
		//TileType
		GeneratedTileType coreModels.TileType
		//MinimalVersion is the version that makes the tile available
		MinimalVersion models.RawVersion
		//Variants list all registered variants (can be available or not)
		Variants map[coreModels.VariantName]*VariantGeneratorSettings
	}

	VariantGeneratorSettings struct {
		//Enabled define if variant is enabled
		Enabled bool

		//VariantName
		VariantName coreModels.VariantName

		//GeneratorFunction function used to generate tile config
		GeneratorFunction models.TileGeneratorFunction
		//GeneratorParamsValidator is used to validate given params for generator
		GeneratorParamsValidator models.ParamsValidator
	}
)

func initConfigData() *ConfigData {
	return &ConfigData{
		tileSettings:          make(map[coreModels.TileType]*TileSetting),
		tileGeneratorSettings: make(map[coreModels.TileType]*TileGeneratorSetting),
	}
}

func (cu *configUsecase) Register(tileType coreModels.TileType, minimalVersion models.RawVersion, variants []coreModels.VariantName) config.TileEnabler {
	tileSetting := &TileSetting{
		TileType:       tileType,
		MinimalVersion: minimalVersion,
		Variants:       make(map[coreModels.VariantName]*VariantSettings),
	}

	// Register Variant with Enabled False
	for _, variant := range variants {
		tileSetting.Variants[variant] = &VariantSettings{
			VariantName: variant,
		}
	}

	cu.configData.tileSettings[tileType] = tileSetting

	return tileSetting
}

func (cu *configUsecase) RegisterGenerator(generatedTileType coreModels.TileType, minimalVersion models.RawVersion, variants []coreModels.VariantName) config.TileGeneratorEnabler {
	tileType := coreModels.NewGeneratorTileType(generatedTileType)
	tileSetting := &TileGeneratorSetting{
		TileType:          tileType,
		GeneratedTileType: generatedTileType,
		MinimalVersion:    minimalVersion,
		Variants:          make(map[coreModels.VariantName]*VariantGeneratorSettings),
	}

	// Register Variant with Enabled False
	for _, variant := range variants {
		tileSetting.Variants[variant] = &VariantGeneratorSettings{
			VariantName: variant,
		}
	}

	cu.configData.tileGeneratorSettings[tileType] = tileSetting

	return tileSetting
}

func (ts *TileSetting) Enable(variant coreModels.VariantName, paramsValidator models.ParamsValidator, routePath string, initialMaxDelay int) {
	variantSetting, exists := ts.Variants[variant]
	if !exists {
		panic(fmt.Sprintf("unable to enable unknown variant: %s for tile: %s. register it before.", variant, ts.TileType))
	}

	variantSetting.Enabled = true
	variantSetting.ParamsValidator = paramsValidator
	variantSetting.RoutePath = &routePath
	variantSetting.InitialMaxDelay = &initialMaxDelay
}

func (ts *TileGeneratorSetting) Enable(variant coreModels.VariantName, generatorParamsValidator models.ParamsValidator, tileGeneratorFunction models.TileGeneratorFunction) {
	variantSetting, exists := ts.Variants[variant]
	if !exists {
		panic(fmt.Sprintf("unable to enable unknown variant: %s for tile: %s. register it before.", variant, ts.TileType))
	}

	variantSetting.Enabled = true
	variantSetting.GeneratorParamsValidator = generatorParamsValidator
	variantSetting.GeneratorFunction = tileGeneratorFunction
}

func (v *VariantSettings) IsEnabled() bool {
	return v.Enabled
}

func (v *VariantSettings) GetValidator() models.ParamsValidator {
	return v.ParamsValidator
}

func (v *VariantGeneratorSettings) IsEnabled() bool {
	return v.Enabled
}

func (v *VariantGeneratorSettings) GetValidator() models.ParamsValidator {
	return v.GeneratorParamsValidator
}
