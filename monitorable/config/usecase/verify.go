package usecase

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/monitoror/monitoror/pkg/monitoror/utils"

	"github.com/monitoror/monitoror/models/tiles"
	"github.com/monitoror/monitoror/monitorable/config/models"
)

const (
	KeyType   = "type"
	KeyLabel  = "label"
	KeyParams = "params"

	EmptyTileType tiles.TileType = "EMPTY"
	GroupTileType tiles.TileType = "GROUP"
)

var (
	AuthorizedTileKey = map[string]bool{
		KeyType:   true,
		KeyLabel:  true,
		KeyParams: true,
	}
)

func (cu *configUsecase) Verify(config *models.Config) error {
	err := models.NewConfigError()
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

	tileType := tiles.TileType(strings.ToUpper(tile[KeyType].(string)))

	// Empty tile, skip
	if tileType == EmptyTileType {
		if group {
			err.Add(`Unauthorized "empty"" type in group tile.`)
		}
		return
	}

	if _, exists := tile[KeyParams]; !exists {
		err.Add(fmt.Sprintf(`Missing "%s" key in %s tile definition.`, KeyParams, tile[KeyType]))
		return
	}

	// Group tile, parse and call verifyTile for each grouped tile
	if tileType == GroupTileType {
		if group {
			err.Add(`Unauthorized "group"" type in group tile.`)
			return
		}

		groupTiles, ok := tile[KeyParams].([]interface{})
		if !ok {
			err.Add(fmt.Sprintf(`Incorrect "%s" key in group tile definition.`, KeyParams))
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

	validator, exists := cu.monitorableParams[tileType]
	if !exists {
		err.Add(fmt.Sprintf(`Unknown "%s" type in tile definition. Must be %s`, tile[KeyType], keys(cu.monitorableParams)))
		return
	}

	// Create new validator by reflexion
	rType := reflect.TypeOf(validator)
	rInstance := reflect.New(rType.Elem()).Interface()

	// Marshal / Unmarshal the map[string]interface{} struct in new instance of Validator
	bParams, _ := json.Marshal(tile[KeyParams])
	unmarshalErr := json.Unmarshal(bParams, &rInstance)

	if unmarshalErr != nil || !rInstance.(utils.Validator).IsValid() {
		err.Add(fmt.Sprintf(`Invalid params definition for "%s": "%s".`, tile[KeyType], string(bParams)))
		return
	}

	return
}

// --- Utility functions ---
func keys(m interface{}) string {
	keys := reflect.ValueOf(m).MapKeys()
	strkeys := make([]string, len(keys))

	for i := 0; i < len(keys); i++ {
		strkeys[i] = strings.ToLower(keys[i].String())
	}

	return strings.Join(strkeys, ",")
}
