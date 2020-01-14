package usecase

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/monitoror/monitoror/config"
	"github.com/monitoror/monitoror/monitorable/config/models"
	"github.com/monitoror/monitoror/pkg/monitoror/utils"
)

func (cu *configUsecase) Verify(conf *models.Config) {
	if conf.Version == nil {
		conf.AddErrors(fmt.Sprintf(`Missing "version" field. Must be %s.`, keys(SupportedVersions)))
		return
	}

	if exists := SupportedVersions[*conf.Version]; !exists {
		conf.AddErrors(fmt.Sprintf(`Unsupported "version" field. Must be %s.`, keys(SupportedVersions)))
		return
	}

	if conf.Columns == nil || *conf.Columns <= 0 {
		conf.AddErrors(`Missing or invalid "columns" field. Must be a positive integer.`)
	}

	if conf.Tiles == nil || len(conf.Tiles) == 0 {
		conf.AddErrors(`Missing or invalid "tiles" field. Must be an array not empty.`)
	} else {
		// Iterating through every tiles conf
		for _, tile := range conf.Tiles {
			cu.verifyTile(conf, &tile, false)
		}
	}
}

func (cu *configUsecase) verifyTile(conf *models.Config, tile *models.Tile, group bool) {
	if tile.ColumnSpan != nil && *tile.ColumnSpan <= 0 {
		conf.AddErrors(`Invalid "columnSpan" field. Must be a positive integer.`)
		return
	}

	if tile.RowSpan != nil && *tile.RowSpan <= 0 {
		conf.AddErrors(`Invalid "rowSpan" field. Must be a positive integer.`)
		return
	}

	// Empty tile, skip
	if tile.Type == EmptyTileType {
		if group {
			conf.AddErrors(fmt.Sprintf(`Unauthorized "%s" type in %s tile.`, EmptyTileType, GroupTileType))
		}
		return
	}

	// Group tile, parse and call verifyTile for each grouped tile
	if tile.Type == GroupTileType {
		if group {
			conf.AddErrors(fmt.Sprintf(`Unauthorized "%s" type in %s tile.`, GroupTileType, GroupTileType))
			return
		}

		if tile.Params != nil {
			conf.AddErrors(fmt.Sprintf(`Unauthorized "params" key in %s tile definition.`, tile.Type))
			return
		}

		if tile.Tiles == nil || len(tile.Tiles) == 0 {
			conf.AddErrors(fmt.Sprintf(`Missing or empty "tiles" key in %s tile definition.`, tile.Type))
			return
		}

		for _, groupTile := range tile.Tiles {
			cu.verifyTile(conf, &groupTile, true)
		}

		return
	}

	if _, exists := cu.tileConfigs[tile.Type]; !exists {
		conf.AddErrors(fmt.Sprintf(`Unknown "%s" type in tile definition. Must be %s`, tile.Type, keys(cu.tileConfigs)))
		return
	}

	if tile.Params == nil {
		conf.AddErrors(fmt.Sprintf(`Missing "params" key in %s tile definition.`, tile.Type))
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
			conf.AddErrors(fmt.Sprintf(`Unknown "%s" variant for %s type in tile definition. Must be %s`,
				tile.ConfigVariant, tile.Type, keys(cu.tileConfigs[tile.Type])))
			return
		}
		validator = tileConfig.Validator
	} else {
		dynamicTileConfig, exists := cu.dynamicTileConfigs[tile.Type][tile.ConfigVariant]
		if !exists {
			conf.AddErrors(fmt.Sprintf(`Unknown "%s" variant for %s dynamic type in tile definition. Must be %s`,
				tile.ConfigVariant, tile.Type, keys(cu.dynamicTileConfigs[tile.Type])))
			return
		}
		validator = dynamicTileConfig.Validator
	}

	// Create new validator by reflexion
	rType := reflect.TypeOf(validator)
	rInstance := reflect.New(rType.Elem()).Interface()

	// Marshal / Unmarshal the map[string]interface{} struct in new instance of Validator
	bParams, _ := json.Marshal(tile.Params)
	unmarshalErr := json.Unmarshal(bParams, &rInstance)

	if unmarshalErr != nil || !rInstance.(utils.Validator).IsValid() {
		conf.AddErrors(fmt.Sprintf(`Invalid params definition for "%s": "%s".`, tile.Type, string(bParams)))
	}
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
