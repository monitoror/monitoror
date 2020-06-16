package usecase

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/fatih/structs"

	"github.com/monitoror/monitoror/api/config/models"
	"github.com/monitoror/monitoror/api/config/versions"
	pkgConfig "github.com/monitoror/monitoror/internal/pkg/api/config"
	"github.com/monitoror/monitoror/internal/pkg/monitorable/params"
	"github.com/monitoror/monitoror/internal/pkg/validator"
	"github.com/monitoror/monitoror/internal/pkg/validator/available"
	"github.com/monitoror/monitoror/internal/pkg/validator/validate"
	coreModels "github.com/monitoror/monitoror/models"
	pkgStructs "github.com/monitoror/monitoror/pkg/structs"
	"github.com/monitoror/monitoror/registry"
)

func (cu *configUsecase) Verify(configBag *models.ConfigBag) {
	if configBag.Config.Version == nil {
		configBag.AddErrors(models.ConfigError{
			ID:      models.ConfigErrorMissingRequiredField,
			Message: fmt.Sprintf(`Required "version" field is missing. Current config version is: %s.`, versions.CurrentVersion),
			Data: models.ConfigErrorData{
				FieldName: "version",
			},
		})
		return
	}

	if configBag.Config.Version.IsLessThan(versions.MinimalVersion) || configBag.Config.Version.IsGreaterThan(versions.CurrentVersion) {
		configBag.AddErrors(models.ConfigError{
			ID:      models.ConfigErrorUnsupportedVersion,
			Message: fmt.Sprintf(`Unsupported configuration version. Minimal supported version is %q. Current config version is: %q`, versions.MinimalVersion, versions.CurrentVersion),
			Data: models.ConfigErrorData{
				FieldName: "version",
				Value:     pkgConfig.Stringify(configBag.Config.Version),
				Expected:  fmt.Sprintf(`%q <= version <= %q`, versions.MinimalVersion, versions.CurrentVersion),
			},
		})
		return
	}

	// Validate struct with "validate" and "available" tag
	errors := validateStruct(configBag.Config, configBag.Config.Version)
	if len(errors) > 0 {
		for _, err := range errors {
			// Convert validator.Error into ConfigError
			configError := convertValidatorError(err, configBag.Config, pkgConfig.Stringify(configBag.Config))
			configBag.AddErrors(*configError)
		}
		return
	}

	// Iterating through every config tiles
	for _, tile := range configBag.Config.Tiles {
		cu.verifyTile(configBag, &tile, nil)
	}
}

func (cu *configUsecase) verifyTile(configBag *models.ConfigBag, tile *models.TileConfig, groupTile *models.TileConfig) {
	// Validate struct with "validate" and "available" tag
	errors := validateStruct(tile, configBag.Config.Version)
	if len(errors) > 0 {
		for _, err := range errors {
			// Convert validator.Error into ConfigError
			configError := convertValidatorError(err, tile, pkgConfig.Stringify(tile))
			configBag.AddErrors(*configError)
		}
		return
	}

	// Empty tile, skip
	if tile.Type == EmptyTileType {
		if groupTile != nil {
			configBag.AddErrors(models.ConfigError{
				ID:      models.ConfigErrorUnauthorizedSubtileType,
				Message: fmt.Sprintf(`Unauthorized %q type in %s tile.`, EmptyTileType, GroupTileType),
				Data: models.ConfigErrorData{
					ConfigExtract:          pkgConfig.Stringify(groupTile),
					ConfigExtractHighlight: pkgConfig.Stringify(tile),
				},
			})
		}
		return
	}

	// Group tile, parse and call verifyTile for each grouped tile
	if tile.Type == GroupTileType {
		if groupTile != nil {
			configBag.AddErrors(models.ConfigError{
				ID:      models.ConfigErrorUnauthorizedSubtileType,
				Message: fmt.Sprintf(`Unauthorized %q type in %s tile.`, GroupTileType, GroupTileType),
				Data: models.ConfigErrorData{
					ConfigExtract:          pkgConfig.Stringify(groupTile),
					ConfigExtractHighlight: pkgConfig.Stringify(tile),
				},
			})
			return
		}

		if tile.Params != nil {
			configBag.AddErrors(models.ConfigError{
				ID:      models.ConfigErrorUnauthorizedField,
				Message: fmt.Sprintf(`Unauthorized "params" key in %s tile definition.`, tile.Type),
				Data: models.ConfigErrorData{
					FieldName:     "params",
					ConfigExtract: pkgConfig.Stringify(tile),
				},
			})
			return
		}

		if tile.Tiles == nil {
			configBag.AddErrors(models.ConfigError{
				ID:      models.ConfigErrorMissingRequiredField,
				Message: fmt.Sprintf(`Missing "tiles" field in %s tile definition. Must be a non-empty array.`, tile.Type),
				Data: models.ConfigErrorData{
					FieldName:     "tiles",
					ConfigExtract: pkgConfig.Stringify(tile),
				},
			})
		} else if len(tile.Tiles) == 0 {
			configBag.AddErrors(models.ConfigError{
				ID:      models.ConfigErrorInvalidFieldValue,
				Message: fmt.Sprintf(`Invalid "tiles" field in %s tile definition. Must be a non-empty array.`, tile.Type),
				Data: models.ConfigErrorData{
					FieldName:     "tiles",
					ConfigExtract: pkgConfig.Stringify(tile),
				},
			})
			return
		}

		for _, groupTile := range tile.Tiles {
			cu.verifyTile(configBag, &groupTile, tile)
		}

		return
	}

	// Set ConfigVariant to DefaultVariantName if empty
	if tile.ConfigVariant == "" {
		tile.ConfigVariant = coreModels.DefaultVariantName
	}

	// Get the metadataExplorer for current tile
	var metadataExplorer registry.TileMetadataExplorer
	var exists bool
	if tile.Type.IsGenerator() {
		// This tile type is a generator tile type
		metadataExplorer, exists = cu.registry.GeneratorMetadata[tile.Type]
		if !exists {
			configBag.AddErrors(models.ConfigError{
				ID:      models.ConfigErrorUnknownGeneratorTileType,
				Message: fmt.Sprintf(`Unknown %q generator type in tile definition. Must be %s`, tile.Type, pkgConfig.Keys(cu.registry.GeneratorMetadata)),
				Data: models.ConfigErrorData{
					FieldName:     "type",
					ConfigExtract: pkgConfig.Stringify(tile),
					Expected:      pkgConfig.Keys(cu.registry.GeneratorMetadata),
				},
			})
			return
		}
	} else {
		// This tile type is a normal tile type
		metadataExplorer, exists = cu.registry.TileMetadata[tile.Type]
		if !exists {
			configBag.AddErrors(models.ConfigError{
				ID:      models.ConfigErrorUnknownTileType,
				Message: fmt.Sprintf(`Unknown %q generator type in tile definition. Must be %s`, tile.Type, pkgConfig.Keys(cu.registry.TileMetadata)),
				Data: models.ConfigErrorData{
					FieldName:     "type",
					ConfigExtract: pkgConfig.Stringify(tile),
					Expected:      pkgConfig.Keys(cu.registry.TileMetadata),
				},
			})
			return
		}
	}

	// Check version on monitarable
	if configBag.Config.Version.IsLessThan(metadataExplorer.GetMinimalVersion()) {
		configBag.AddErrors(models.ConfigError{
			ID: models.ConfigErrorUnsupportedTileInThisVersion,
			Message: fmt.Sprintf(`%q tile type is not supported in version %q. Minimal supported version is %q`,
				tile.Type, configBag.Config.Version, string(metadataExplorer.GetMinimalVersion())),
			Data: models.ConfigErrorData{
				FieldName:     "type",
				ConfigExtract: pkgConfig.Stringify(tile),
				Expected:      fmt.Sprintf(`version >= %q`, string(metadataExplorer.GetMinimalVersion())),
			},
		})
		return
	}

	// Check if variant of monitoarable exist
	variantMetadataExplorer, exists := metadataExplorer.GetVariant(tile.ConfigVariant)
	if !exists {
		configBag.AddErrors(models.ConfigError{
			ID: models.ConfigErrorUnknownVariant,
			Message: fmt.Sprintf(`Unknown %q variant for %s type in tile definition. Must be %s`,
				tile.ConfigVariant, tile.Type, pkgConfig.Stringify(metadataExplorer.GetVariantsNames())),
			Data: models.ConfigErrorData{
				FieldName:     "configVariant",
				Value:         pkgConfig.Stringify(tile.ConfigVariant),
				Expected:      pkgConfig.Stringify(metadataExplorer.GetVariantsNames()),
				ConfigExtract: pkgConfig.Stringify(tile),
			},
		})
		return
	}

	// Check if this variant is enabled
	if !variantMetadataExplorer.IsEnabled() {
		configBag.AddErrors(models.ConfigError{
			ID:      models.ConfigErrorDisabledVariant,
			Message: fmt.Sprintf(`Variant %q is disabled for %s type. Check errors on the server side for the reason`, tile.ConfigVariant, tile.Type),
			Data: models.ConfigErrorData{
				FieldName:     "configVariant",
				Value:         pkgConfig.Stringify(tile.ConfigVariant),
				ConfigExtract: pkgConfig.Stringify(tile),
			},
		})
		return
	}

	// Check if param isn't empty
	if tile.Params == nil {
		configBag.AddErrors(models.ConfigError{
			ID:      models.ConfigErrorMissingRequiredField,
			Message: fmt.Sprintf(`Missing "params" key in %s tile definition.`, tile.Type),
			Data: models.ConfigErrorData{
				FieldName:     "params",
				ConfigExtract: pkgConfig.Stringify(tile),
			},
		})
		return
	}

	// Create new validator by reflexion
	rType := reflect.TypeOf(variantMetadataExplorer.GetValidator())
	rInstance := reflect.New(rType.Elem()).Interface()

	// Marshal / Unmarshal the map[string]interface{} struct in new instance of ParamsValidator
	bytesParams, _ := json.Marshal(tile.Params)
	if unmarshalErr := json.Unmarshal(bytesParams, &rInstance); unmarshalErr != nil {
		configBag.AddErrors(models.ConfigError{
			ID:      models.ConfigErrorUnexpectedError,
			Message: unmarshalErr.Error(),
			Data: models.ConfigErrorData{
				FieldName:     "params",
				ConfigExtract: pkgConfig.Stringify(tile),
			},
		})
		return
	}

	// Marshal / Unmarshal instance of validator into map[string]interface{} to identify unknown fields
	structParams := make(map[string]interface{})
	bytesParamInstance, _ := json.Marshal(rInstance)
	_ = json.Unmarshal(bytesParamInstance, &structParams)

	// Check if struct has unknown fields
	for field := range tile.Params {
		if _, ok := structParams[field]; !ok {
			configBag.AddErrors(models.ConfigError{
				ID:      models.ConfigErrorUnknownField,
				Message: fmt.Sprintf(`Unknown %q tile params field.`, field),
				Data: models.ConfigErrorData{
					FieldName:     field,
					ConfigExtract: pkgConfig.Stringify(tile),
					Expected:      pkgConfig.Keys(structParams),
				},
			})
			return
		}
	}

	// Validate struct with "validate" and "available" tag
	castedParams := rInstance.(params.Validator)
	errors = validateStruct(castedParams, configBag.Config.Version)
	errors = append(errors, castedParams.Validate()...)

	for _, vError := range errors {
		configError := convertValidatorError(vError, rInstance, pkgConfig.Stringify(tile))

		// UX HACK: if params is empty, inject "params:{}" to help users
		if len(tile.Params) == 0 {
			for _, field := range structs.Fields(tile) {
				if reflect.DeepEqual(field.Value(), tile.Params) {
					paramName := pkgStructs.GetJSONFieldName(field)
					configError.Data.ConfigExtract = strings.TrimSuffix(configError.Data.ConfigExtract, "}")
					configError.Data.ConfigExtract = fmt.Sprintf(`%s,"%s":{}}`, configError.Data.ConfigExtract, paramName)
					break
				}
			}
		}

		configBag.AddErrors(*configError)
	}
}

// validateStruct Validate struct with "validate" and "available" tag
func validateStruct(s interface{}, version *versions.ConfigVersion) []validator.Error {
	var errors []validator.Error

	// use "available" tag in struct definition to validate params
	errors = append(errors, available.Struct(s, version)...)
	// use "validate" tag in struct definition to validate params
	errors = append(errors, validate.Struct(s)...)

	return errors
}

// convertValidatorError into models.ConfigError
func convertValidatorError(vError validator.Error, instance interface{}, configExtract string) *models.ConfigError {
	configError := &models.ConfigError{Data: models.ConfigErrorData{}}

	// Convert error ID
	switch vError.GetErrorID() {
	case validator.ErrorRequired:
		configError.ID = models.ConfigErrorMissingRequiredField
	case validator.ErrorSince, validator.ErrorUntil:
		configError.ID = models.ConfigErrorUnsupportedTileParamInThisVersion
	default:
		configError.ID = models.ConfigErrorInvalidFieldValue
	}

	// Inject extract
	configError.Data.ConfigExtract = configExtract

	// Convert FieldName into json fieldName
	for _, field := range structs.Fields(instance) {
		if field.Name() == vError.GetFieldName() {
			// Replace FieldName By json FieldName
			if jsonTagValue := pkgStructs.GetJSONFieldName(field); jsonTagValue != "" && jsonTagValue != "-" {
				vError.SetFieldName(jsonTagValue)
			}
			break
		}
	}
	configError.Data.FieldName = vError.GetFieldName()

	// Convert Message
	configError.Message = vError.Error()

	// Convert Expected
	configError.Data.Expected = vError.Expected()

	return configError
}
