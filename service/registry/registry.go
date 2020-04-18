//go:generate mockery -name Registry|TileEnabler|GeneratorEnabler|TileMetadataExplorer|VariantMetadataExplorer -output ../mocks

package registry

import (
	"fmt"

	"github.com/monitoror/monitoror/api/config/models"
	"github.com/monitoror/monitoror/api/config/versions"
	coreModels "github.com/monitoror/monitoror/models"
)

type (
	// Registry is used to register Tile and Generator in config for verify / hydrate
	Registry interface {
		RegisterTile(tileType coreModels.TileType, minimalVersion versions.RawVersion, variantNames []coreModels.VariantName) TileEnabler
		RegisterGenerator(generatedTileType coreModels.TileType, minimalVersion versions.RawVersion, variantNames []coreModels.VariantName) GeneratorEnabler
	}
	// TileEnabler is returned to monitorable after register to enable monitorable tile with this variant if she is "valid"
	TileEnabler interface {
		Enable(variantName coreModels.VariantName, paramsValidator models.ParamsValidator, routePath string)
	}
	// GeneratorEnabler is returned to monitorable after register to enable monitorable generator with this variant if she is "valid"
	GeneratorEnabler interface {
		Enable(variantName coreModels.VariantName, generatorParamsValidator models.ParamsValidator, tileGeneratorFunction models.TileGeneratorFunction)
	}

	// TileMetadataExplorer is used in verify. Matching tileMetadata and generatorMetadata.
	TileMetadataExplorer interface {
		GetMinimalVersion() versions.RawVersion
		GetVariant(variantName coreModels.VariantName) (VariantMetadataExplorer, bool)
		GetVariantsNames() []coreModels.VariantName
	}
	// VariantMetadataExplorer is used in verify. Matching tileVariantMetadata and generatorVariantMetadata.
	VariantMetadataExplorer interface {
		IsEnabled() bool
		GetValidator() models.ParamsValidator
	}
)

type (
	MetadataRegistry struct {
		TileMetadata      map[coreModels.TileType]*tileMetadata
		GeneratorMetadata map[coreModels.TileType]*generatorMetadata
	}

	tileMetadata struct {
		// TileType
		TileType coreModels.TileType
		// MinimalVersion is the version that makes the tile available
		MinimalVersion versions.RawVersion
		// VariantsMetadata list all registered variants (can be available or not)
		VariantsMetadata map[coreModels.VariantName]*tileVariantMetadata
	}

	tileVariantMetadata struct {
		// Enabled define if variant is enabled
		Enabled bool

		// VariantName
		VariantName coreModels.VariantName

		// RoutePath path of the api endpoint for this tile. Used by hydrate
		RoutePath *string
		// ParamsValidator is used to validate given params
		ParamsValidator models.ParamsValidator
	}

	generatorMetadata struct {
		// TileType
		TileType coreModels.TileType
		// GeneratedTileType
		GeneratedTileType coreModels.TileType
		// MinimalVersion is the version that makes the tile available
		MinimalVersion versions.RawVersion
		// VariantsMetadata list all registered variants (can be available or not)
		VariantsMetadata map[coreModels.VariantName]*generatorVariantMetadata
	}

	generatorVariantMetadata struct {
		// Enabled define if variant is enabled
		Enabled bool

		// VariantName
		VariantName coreModels.VariantName

		// GeneratorFunction function used to generate tile config
		GeneratorFunction models.TileGeneratorFunction
		// GeneratorParamsValidator is used to validate given params for generator
		GeneratorParamsValidator models.ParamsValidator
	}
)

func NewRegistry() *MetadataRegistry {
	return &MetadataRegistry{
		TileMetadata:      make(map[coreModels.TileType]*tileMetadata),
		GeneratorMetadata: make(map[coreModels.TileType]*generatorMetadata),
	}
}

// REGISTRY
// ----------------------------------------
func (r *MetadataRegistry) RegisterTile(tileType coreModels.TileType, minimalVersion versions.RawVersion, variantNames []coreModels.VariantName) TileEnabler {
	tileSetting := &tileMetadata{
		TileType:         tileType,
		MinimalVersion:   minimalVersion,
		VariantsMetadata: make(map[coreModels.VariantName]*tileVariantMetadata),
	}

	// Register Variant with Enabled False
	for _, variantName := range variantNames {
		tileSetting.VariantsMetadata[variantName] = &tileVariantMetadata{
			VariantName: variantName,
		}
	}

	r.TileMetadata[tileType] = tileSetting

	return tileSetting
}

func (r *MetadataRegistry) RegisterGenerator(generatedTileType coreModels.TileType, minimalVersion versions.RawVersion, variantNames []coreModels.VariantName) GeneratorEnabler {
	// Boxing tile type into generator
	tileType := coreModels.NewGeneratorTileType(generatedTileType)

	generatorSetting := &generatorMetadata{
		TileType:          tileType,
		GeneratedTileType: generatedTileType,
		MinimalVersion:    minimalVersion,
		VariantsMetadata:  make(map[coreModels.VariantName]*generatorVariantMetadata),
	}

	// Register Variant with Enabled False
	for _, variantName := range variantNames {
		generatorSetting.VariantsMetadata[variantName] = &generatorVariantMetadata{
			VariantName: variantName,
		}
	}

	r.GeneratorMetadata[tileType] = generatorSetting

	return generatorSetting
}

// ----------------------------------------

// TILE METADATA
// ----------------------------------------
func (tm *tileMetadata) Enable(variantName coreModels.VariantName, paramsValidator models.ParamsValidator, routePath string) {
	variantMetadata, exists := tm.VariantsMetadata[variantName]
	if !exists {
		panic(fmt.Sprintf("unable to enable unknown variantName: %s for tile: %s. register it before.", variantName, tm.TileType))
	}

	variantMetadata.Enabled = true
	variantMetadata.ParamsValidator = paramsValidator
	variantMetadata.RoutePath = &routePath
}

func (tm *tileMetadata) GetMinimalVersion() versions.RawVersion {
	return tm.MinimalVersion
}

func (tm *tileMetadata) GetVariant(variantName coreModels.VariantName) (VariantMetadataExplorer, bool) {
	v, exists := tm.VariantsMetadata[variantName]
	return v, exists
}

func (tm *tileMetadata) GetVariantsNames() []coreModels.VariantName {
	var result []coreModels.VariantName
	for variantName := range tm.VariantsMetadata {
		result = append(result, variantName)
	}
	return result
}

// ----------------------------------------

// GENERATOR METADATA
// ----------------------------------------
func (gm *generatorMetadata) Enable(variantName coreModels.VariantName, generatorParamsValidator models.ParamsValidator, tileGeneratorFunction models.TileGeneratorFunction) {
	variantMetadata, exists := gm.VariantsMetadata[variantName]
	if !exists {
		panic(fmt.Sprintf("unable to enable unknown variantName: %s for tile: %s. register it before.", variantName, gm.TileType))
	}

	variantMetadata.Enabled = true
	variantMetadata.GeneratorParamsValidator = generatorParamsValidator
	variantMetadata.GeneratorFunction = tileGeneratorFunction
}

func (gm *generatorMetadata) GetMinimalVersion() versions.RawVersion {
	return gm.MinimalVersion
}

func (gm *generatorMetadata) GetVariant(variantName coreModels.VariantName) (VariantMetadataExplorer, bool) {
	v, exists := gm.VariantsMetadata[variantName]
	return v, exists
}

func (gm *generatorMetadata) GetVariantsNames() []coreModels.VariantName {
	var result []coreModels.VariantName
	for variantName := range gm.VariantsMetadata {
		result = append(result, variantName)
	}
	return result
}

// ----------------------------------------

// TILE VARIANT METADATA
// ----------------------------------------
func (tvm *tileVariantMetadata) IsEnabled() bool {
	return tvm.Enabled
}
func (tvm *tileVariantMetadata) GetValidator() models.ParamsValidator {
	return tvm.ParamsValidator
}

// ----------------------------------------

// GENERATOR VARIANT METADATA
// ----------------------------------------
func (gvm *generatorVariantMetadata) IsEnabled() bool {
	return gvm.Enabled
}
func (gvm *generatorVariantMetadata) GetValidator() models.ParamsValidator {
	return gvm.GeneratorParamsValidator
}

// ----------------------------------------
