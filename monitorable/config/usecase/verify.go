package usecase

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/monitoror/monitoror/config"
	"github.com/monitoror/monitoror/monitorable/config/models"
	"github.com/monitoror/monitoror/pkg/monitoror/utils"
)

func (cu *configUsecase) Verify(configBag *models.ConfigBag) {
	if configBag.Config.Version == nil {
		configBag.AddErrors(models.ConfigError{
			ID:      models.ConfigErrorMissingRequiredField,
			Message: fmt.Sprintf(`Required "version" field is missing. Current config version is: %s.`, CurrentVersion),
			Data: models.ConfigErrorData{
				FieldName: "version",
			},
		})
		return
	}

	if configBag.Config.Version.IsLessThan(MinimalVersion) || configBag.Config.Version.IsGreaterThan(CurrentVersion) {
		configBag.AddErrors(models.ConfigError{
			ID:      models.ConfigErrorUnsupportedVersion,
			Message: fmt.Sprintf(`Unsupported configuration version. Minimal supported version is %q. Current config version is: %q`, MinimalVersion, CurrentVersion),
			Data: models.ConfigErrorData{
				FieldName: "version",
				Value:     stringify(configBag.Config.Version),
				Expected:  fmt.Sprintf(`%q >= version >= %q`, MinimalVersion, CurrentVersion),
			},
		})
		return
	}

	if configBag.Config.Columns == nil {
		configBag.AddErrors(models.ConfigError{
			ID:      models.ConfigErrorMissingRequiredField,
			Message: fmt.Sprintf(`Required "columns" field is missing. Must be a positive integer.`),
			Data: models.ConfigErrorData{
				FieldName: "columns",
			},
		})
	} else if *configBag.Config.Columns <= 0 {
		configBag.AddErrors(models.ConfigError{
			ID:      models.ConfigErrorInvalidFieldValue,
			Message: fmt.Sprintf(`Invalid "columns" field. Must be a positive integer.`),
			Data: models.ConfigErrorData{
				FieldName: "columns",
				Value:     stringify(configBag.Config.Columns),
				Expected:  "columns > 0",
			},
		})
	}

	if configBag.Config.Zoom != nil && (*configBag.Config.Zoom <= 0 || *configBag.Config.Zoom > 10) {
		configBag.AddErrors(models.ConfigError{
			ID:      models.ConfigErrorInvalidFieldValue,
			Message: `Invalid "zoom" field. Must be a positive float between 0 and 10.`,
			Data: models.ConfigErrorData{
				FieldName: "zoom",
				Value:     stringify(configBag.Config.Zoom),
				Expected:  "0 < zoom <= 10",
			},
		})
	}

	if configBag.Config.Tiles == nil {
		configBag.AddErrors(models.ConfigError{
			ID:      models.ConfigErrorMissingRequiredField,
			Message: `Missing "tiles" field. Must be a non-empty array.`,
			Data: models.ConfigErrorData{
				FieldName: "tiles",
			},
		})
	} else if len(configBag.Config.Tiles) == 0 {
		configBag.AddErrors(models.ConfigError{
			ID:      models.ConfigErrorInvalidFieldValue,
			Message: `Invalid "tiles" field. Must be a non-empty array.`,
			Data: models.ConfigErrorData{
				FieldName:     "tiles",
				ConfigExtract: stringify(configBag.Config),
			},
		})
	} else {
		// Iterating through every config tiles
		for _, tile := range configBag.Config.Tiles {
			cu.verifyTile(configBag, &tile, nil)
		}
	}
}

func (cu *configUsecase) verifyTile(configBag *models.ConfigBag, tile *models.TileConfig, groupTile *models.TileConfig) {
	if tile.ColumnSpan != nil && *tile.ColumnSpan <= 0 {
		configBag.AddErrors(models.ConfigError{
			ID:      models.ConfigErrorInvalidFieldValue,
			Message: `Invalid "columnSpan" field. Must be a positive integer.`,
			Data: models.ConfigErrorData{
				FieldName: "columnSpan",
				Value:     stringify(*tile.ColumnSpan),
				Expected:  "columnSpan > 0",
			},
		})
		return
	}

	if tile.RowSpan != nil && *tile.RowSpan <= 0 {
		configBag.AddErrors(models.ConfigError{
			ID:      models.ConfigErrorInvalidFieldValue,
			Message: `Invalid "rowSpan" field. Must be a positive integer.`,
			Data: models.ConfigErrorData{
				FieldName: "rowSpan",
				Value:     stringify(*tile.RowSpan),
				Expected:  "rowSpan > 0"},
		})
		return
	}

	// Empty tile, skip
	if tile.Type == EmptyTileType {
		if groupTile != nil {
			configBag.AddErrors(models.ConfigError{
				ID:      models.ConfigErrorUnauthorizedSubtileType,
				Message: fmt.Sprintf(`Unauthorized %q type in %s tile.`, EmptyTileType, GroupTileType),
				Data: models.ConfigErrorData{
					ConfigExtract:          stringify(groupTile),
					ConfigExtractHighlight: stringify(tile),
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
					ConfigExtract:          stringify(groupTile),
					ConfigExtractHighlight: stringify(tile),
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
					ConfigExtract: stringify(tile),
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
					ConfigExtract: stringify(tile),
				},
			})
		} else if len(tile.Tiles) == 0 {
			configBag.AddErrors(models.ConfigError{
				ID:      models.ConfigErrorInvalidFieldValue,
				Message: fmt.Sprintf(`Invalid "tiles" field in %s tile definition. Must be a non-empty array.`, tile.Type),
				Data: models.ConfigErrorData{
					FieldName:     "tiles",
					ConfigExtract: stringify(tile),
				},
			})
			return
		}

		for _, groupTile := range tile.Tiles {
			cu.verifyTile(configBag, &groupTile, tile)
		}

		return
	}

	if _, exists := cu.tileConfigs[tile.Type]; !exists {
		configBag.AddErrors(models.ConfigError{
			ID:      models.ConfigErrorUnknownTileType,
			Message: fmt.Sprintf(`Unknown %q type in tile definition. Must be %s`, tile.Type, keys(cu.tileConfigs)),
			Data: models.ConfigErrorData{
				FieldName:     "type",
				ConfigExtract: stringify(tile),
				Expected:      keys(cu.tileConfigs),
			},
		})
		return
	}

	if tile.Params == nil {
		configBag.AddErrors(models.ConfigError{
			ID:      models.ConfigErrorMissingRequiredField,
			Message: fmt.Sprintf(`Missing "params" key in %s tile definition.`, tile.Type),
			Data: models.ConfigErrorData{
				FieldName:     "params",
				ConfigExtract: stringify(tile),
			},
		})
		return
	}

	// Set ConfigVariant to DefaultVariant if empty
	if tile.ConfigVariant == "" {
		tile.ConfigVariant = config.DefaultVariant
	}

	// Get the validator for current tile
	// - for non dynamic tile, the validator is register in tileConfigs
	// - for dynamic tile, the validator is register in dynamicTileConfigs
	var validator utils.Validator
	if _, exists := cu.dynamicTileConfigs[tile.Type]; !exists {
		tileConfig, exists := cu.tileConfigs[tile.Type][tile.ConfigVariant]
		if !exists {
			configBag.AddErrors(models.ConfigError{
				ID: models.ConfigErrorUnknownVariant,
				Message: fmt.Sprintf(`Unknown %q variant for %s type in tile definition. Must be %s`,
					tile.ConfigVariant, tile.Type, keys(cu.tileConfigs[tile.Type])),
				Data: models.ConfigErrorData{
					FieldName:     "configVariant",
					Value:         stringify(tile.ConfigVariant),
					ConfigExtract: stringify(tile),
				},
			})
			return
		}
		validator = tileConfig.Validator
	} else {
		dynamicTileConfig, exists := cu.dynamicTileConfigs[tile.Type][tile.ConfigVariant]
		if !exists {
			configBag.AddErrors(models.ConfigError{
				ID: models.ConfigErrorUnknownVariant,
				Message: fmt.Sprintf(`Unknown %q variant for %s dynamic type in tile definition. Must be %s`,
					tile.ConfigVariant, tile.Type, keys(cu.dynamicTileConfigs[tile.Type])),
				Data: models.ConfigErrorData{
					FieldName:     "configVariant",
					Value:         stringify(tile.ConfigVariant),
					ConfigExtract: stringify(tile),
				},
			})
			return
		}
		validator = dynamicTileConfig.Validator
	}

	// Create new validator by reflexion
	rType := reflect.TypeOf(validator)
	rInstance := reflect.New(rType.Elem()).Interface()

	// Marshal / Unmarshal the map[string]interface{} struct in new instance of Validator
	bytesParams, _ := json.Marshal(tile.Params)
	unmarshalErr := json.Unmarshal(bytesParams, &rInstance)

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
					ConfigExtract: stringify(tile),
					Expected:      keys(structParams),
				},
			})
			return
		}
	}

	if unmarshalErr != nil || !rInstance.(utils.Validator).IsValid() {
		configBag.AddErrors(models.ConfigError{
			ID:      models.ConfigErrorInvalidFieldValue,
			Message: fmt.Sprintf(`Invalid params definition for %q: %q.`, tile.Type, string(bytesParams)),
			Data: models.ConfigErrorData{
				FieldName:     "params",
				ConfigExtract: stringify(tile),
			},
		})
	}
}
