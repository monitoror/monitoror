package usecase

import (
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"

	"github.com/monitoror/monitoror/config"

	"github.com/monitoror/monitoror/monitorable/config/models"
)

func (cu *configUsecase) Hydrate(config *models.Config, host string) error {
	err := models.NewConfigError()

	cu.hydrateTiles(&config.Tiles, host, err)

	if err.Count() > 0 {
		return err
	}
	return nil
}

func (cu *configUsecase) hydrateTiles(tiles *[]models.Tile, host string, err *models.ConfigError) {
	for i := 0; i < len(*tiles); i++ {
		tile := &(*tiles)[i]
		if tile.Type != EmptyTileType && tile.Type != GroupTileType {
			// Set ConfigVariant to DefaultVariant if empty
			if tile.ConfigVariant == "" {
				tile.ConfigVariant = config.DefaultVariant
			}
		}

		if _, exists := cu.dynamicTileConfigs[tile.Type]; !exists {
			cu.hydrateTile(tile, host, err)
		} else {
			dynamicTiles := cu.hydrateDynamicTile(tile, err)

			temp := append((*tiles)[:i], dynamicTiles...)
			*tiles = append(temp, (*tiles)[i+1:]...)

			i--
		}
	}
}

func (cu *configUsecase) hydrateTile(tile *models.Tile, host string, err *models.ConfigError) {
	// Empty tile, skip
	if tile.Type == EmptyTileType {
		return
	}

	if tile.Type == GroupTileType {
		cu.hydrateTiles(&tile.Tiles, host, err)
		return
	}

	// Change Params by a valid Url
	path := cu.tileConfigs[tile.Type][tile.ConfigVariant].Path
	urlParams := url.Values{}
	for key, value := range tile.Params {
		urlParams.Add(key, fmt.Sprintf("%v", value))
	}

	tile.Url = fmt.Sprintf("%s%s?%s", host, path, urlParams.Encode())

	// Remove Params / Variant
	tile.Params = nil
	tile.ConfigVariant = ""
}

func (cu *configUsecase) hydrateDynamicTile(tile *models.Tile, err *models.ConfigError) (tiles []models.Tile) {
	config := cu.dynamicTileConfigs[tile.Type][tile.ConfigVariant]

	// Create new validator by reflexion
	rType := reflect.TypeOf(config.Validator)
	rInstance := reflect.New(rType.Elem()).Interface()

	// Marshal / Unmarshal the map[string]interface{} struct in new instance of Validator
	bParams, _ := json.Marshal(tile.Params)
	_ = json.Unmarshal(bParams, &rInstance)

	// Call builder and add inherited value from Dynamic tile
	results, e := config.Builder.ListDynamicTile(rInstance)
	if e != nil {
		err.Add(fmt.Sprintf(`Error while listing %s dynamic tiles. %v`, tile.Type, e))
		return
	}

	tiles = []models.Tile{}
	for _, result := range results {
		newTile := models.Tile{
			Type:          result.TileType,
			Label:         result.Label,
			Params:        result.Params,
			ConfigVariant: tile.ConfigVariant,
			ColumnSpan:    tile.ColumnSpan,
			RowSpan:       tile.RowSpan,
		}

		tiles = append(tiles, newTile)
	}

	return
}
