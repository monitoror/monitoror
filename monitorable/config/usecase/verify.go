package usecase

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	. "github.com/monitoror/monitoror/pkg/monitoror/validator"

	"github.com/monitoror/monitoror/models/tiles"
	"github.com/monitoror/monitoror/monitorable/config/models"
)

func (cu *configUsecase) Verify(config *models.Config) error {
	err := models.NewConfigError()

	if exists := SupportedVersions[config.Version]; !exists {
		err.Add(fmt.Sprintf(`Unsupported "version" field. Must be %s.`, keys(SupportedVersions)))
	}

	if config.Columns == 0 {
		err.Add(`Missing or invalid "columns" field. Must be a positive integer.`)
	}

	if config.Tiles == nil || len(config.Tiles) == 0 {
		err.Add(`Missing or invalid "tiles" field. Must be an array not empty.`)
	} else {
		// Iterating through every tiles config
		for _, tile := range config.Tiles {
			cu.verifyTile(tile, false, err)
		}
	}

	if err.Count() > 0 {
		return err
	}
	return nil
}

func (cu *configUsecase) verifyTile(tile map[string]interface{}, group bool, err *models.ConfigError) {
	// Check if tile keys are authorized
	for key := range tile {
		if exists := AuthorizedTileKey[key]; !exists {
			err.Add(fmt.Sprintf(`Unknown key "%s" in tile definition. Must be %s.`, key, keys(AuthorizedTileKey)))
			return
		}
	}

	tileType := tiles.TileType(strings.ToUpper(tile[TypeKey].(string)))

	// Empty tile, skip
	if tileType == EmptyTileType {
		if group {
			err.Add(`Unauthorized "empty" type in group tile.`)
		}
		return
	}

	// Group tile, parse and call verifyTile for each grouped tile
	if tileType == GroupTileType {
		if group {
			err.Add(`Unauthorized "group" type in group tile.`)
			return
		}

		if _, exists := tile[ParamsKey]; exists {
			err.Add(fmt.Sprintf(`Unauthorized "%s" key in %s tile definition.`, ParamsKey, tile[TypeKey]))
			return
		}

		groupTiles, ok := tile[TilesKey].([]interface{})
		if !ok {
			err.Add(fmt.Sprintf(`Incorrect "%s" key in %s tile definition.`, TilesKey, tile[TypeKey]))
			return
		}

		for _, gt := range groupTiles {
			groupTile, ok := gt.(map[string]interface{})
			if !ok {
				err.Add(fmt.Sprintf(`Incorrect array element "%s" in group definition.`, gt))
				continue
			}

			cu.verifyTile(groupTile, true, err)
		}

		return
	}

	if _, exists := tile[ParamsKey]; !exists {
		err.Add(fmt.Sprintf(`Missing "%s" key in %s tile definition.`, ParamsKey, tile[TypeKey]))
		return
	}

	tileConfig, exists := cu.tileConfigs[tileType]
	if !exists {
		err.Add(fmt.Sprintf(`Unknown "%s" type in tile definition. Must be %s`, tile[TypeKey], keys(cu.tileConfigs)))
		return
	}

	// Create new validator by reflexion
	rType := reflect.TypeOf(tileConfig.Validator)
	rInstance := reflect.New(rType.Elem()).Interface()

	// Marshal / Unmarshal the map[string]interface{} struct in new instance of Validator
	bParams, _ := json.Marshal(tile[ParamsKey])
	unmarshalErr := json.Unmarshal(bParams, &rInstance)

	if unmarshalErr != nil || !rInstance.(Validator).IsValid() {
		err.Add(fmt.Sprintf(`Invalid params definition for "%s": "%s".`, tile[TypeKey], string(bParams)))
		return
	}

	return
}

// --- Utility functions ---
func keys(m interface{}) string {
	keys := reflect.ValueOf(m).MapKeys()
	strkeys := make([]string, len(keys))

	for i := 0; i < len(keys); i++ {
		strkeys[i] = strings.ToLower(fmt.Sprintf("%v", keys[i]))
	}

	return strings.Join(strkeys, ",")
}
