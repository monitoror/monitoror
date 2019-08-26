package usecase

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/monitoror/monitoror/config"

	. "github.com/monitoror/monitoror/pkg/monitoror/validator"

	"github.com/monitoror/monitoror/monitorable/config/models"
)

func (cu *configUsecase) Verify(config *models.Config) error {
	err := models.NewConfigError()

	if exists := SupportedVersions[config.Version]; !exists {
		err.Add(fmt.Sprintf(`Unsupported "version" field. Must be %s.`, keys(SupportedVersions)))
	}

	if config.Columns <= 0 {
		err.Add(`Missing or invalid "columns" field. Must be a positive integer.`)
	}

	if config.Tiles == nil || len(config.Tiles) == 0 {
		err.Add(`Missing or invalid "tiles" field. Must be an array not empty.`)
	} else {
		// Iterating through every tiles config
		for _, tile := range config.Tiles {
			cu.verifyTile(&tile, false, err)
		}
	}

	if err.Count() > 0 {
		return err
	}
	return nil
}

func (cu *configUsecase) verifyTile(tile *models.Tile, group bool, err *models.ConfigError) {
	if tile.ColumnSpan != nil && *tile.ColumnSpan <= 0 {
		err.Add(`Invalid "columnSpan" field. Must be a positive integer.`)
		return
	}

	if tile.RowSpan != nil && *tile.RowSpan <= 0 {
		err.Add(`Invalid "rowSpan" field. Must be a positive integer.`)
		return
	}

	// Empty tile, skip
	if tile.Type == EmptyTileType {
		if group {
			err.Add(fmt.Sprintf(`Unauthorized "%s" type in %s tile.`, EmptyTileType, GroupTileType))
		}
		return
	}

	// Group tile, parse and call verifyTile for each grouped tile
	if tile.Type == GroupTileType {
		if group {
			err.Add(fmt.Sprintf(`Unauthorized "%s" type in %s tile.`, GroupTileType, GroupTileType))
			return
		}

		if tile.Params != nil {
			err.Add(fmt.Sprintf(`Unauthorized "params" key in %s tile definition.`, tile.Type))
			return
		}

		if tile.Tiles == nil || len(tile.Tiles) == 0 {
			err.Add(fmt.Sprintf(`Missing or empty "tiles" key in %s tile definition.`, tile.Type))
			return
		}

		for _, groupTile := range tile.Tiles {
			cu.verifyTile(&groupTile, true, err)
		}

		return
	}

	if _, exists := cu.tileConfigs[tile.Type]; !exists {
		err.Add(fmt.Sprintf(`Unknown "%s" type in tile definition. Must be %s`, tile.Type, keys(cu.tileConfigs)))
		return
	}

	if tile.Params == nil {
		err.Add(fmt.Sprintf(`Missing "params" key in %s tile definition.`, tile.Type))
		return
	}

	// Set ConfigVariant to DefaultVariant if empty
	if tile.ConfigVariant == "" {
		tile.ConfigVariant = config.DefaultVariant
	}

	// Get the validator for current tile
	// - for non dynamic tile, the validator is register in tileConfigs
	// - for dynamic tile, the validator is register in dynamicTileConfigs
	var validator Validator
	if _, exists := cu.dynamicTileConfigs[tile.Type]; !exists {
		if tileConfig, exists := cu.tileConfigs[tile.Type][tile.ConfigVariant]; !exists {
			err.Add(fmt.Sprintf(`Unknown "%s" variant for %s type in tile definition. Must be %s`,
				tile.ConfigVariant, tile.Type, keys(cu.tileConfigs[tile.Type])))
			return
		} else {
			validator = tileConfig.Validator
		}
	} else {
		if dynamicTileConfig, exists := cu.dynamicTileConfigs[tile.Type][tile.ConfigVariant]; !exists {
			err.Add(fmt.Sprintf(`Unknown "%s" variant for %s dynamic type in tile definition. Must be %s`,
				tile.ConfigVariant, tile.Type, keys(cu.dynamicTileConfigs[tile.Type])))
			return
		} else {
			validator = dynamicTileConfig.Validator
		}
	}

	// Create new validator by reflexion
	rType := reflect.TypeOf(validator)
	rInstance := reflect.New(rType.Elem()).Interface()

	// Marshal / Unmarshal the map[string]interface{} struct in new instance of Validator
	bParams, _ := json.Marshal(tile.Params)
	unmarshalErr := json.Unmarshal(bParams, &rInstance)

	if unmarshalErr != nil || !rInstance.(Validator).IsValid() {
		err.Add(fmt.Sprintf(`Invalid params definition for "%s": "%s".`, tile.Type, string(bParams)))
		return
	}

	return
}

// --- Utility functions ---
func keys(m interface{}) string {
	keys := reflect.ValueOf(m).MapKeys()
	strkeys := make([]string, len(keys))

	for i := 0; i < len(keys); i++ {
		strkeys[i] = fmt.Sprintf(`%v`, keys[i])
	}

	return strings.Join(strkeys, ",")
}
