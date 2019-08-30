package usecase

import (
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"
	"strings"

	"github.com/monitoror/monitoror/config"

	"github.com/monitoror/monitoror/monitorable/config/models"
)

func (cu *configUsecase) Hydrate(conf *models.Config, host string) {
	cu.hydrateTiles(conf, &conf.Tiles, host)
}

func (cu *configUsecase) hydrateTiles(conf *models.Config, tiles *[]models.Tile, host string) {
	for i := 0; i < len(*tiles); i++ {
		tile := &((*tiles)[i])
		if tile.Type != EmptyTileType && tile.Type != GroupTileType {
			// Set ConfigVariant to DefaultVariant if empty
			if tile.ConfigVariant == "" {
				tile.ConfigVariant = config.DefaultVariant
			}
		}

		if _, exists := cu.dynamicTileConfigs[tile.Type]; !exists {
			cu.hydrateTile(conf, tile, host)
		} else {
			dynamicTiles := cu.hydrateDynamicTile(conf, tile)

			// Remove DynamicTile config and add real dynamic tiles in array
			temp := append((*tiles)[:i], dynamicTiles...)
			*tiles = append(temp, (*tiles)[i+1:]...)

			i--
		}
	}
}

func (cu *configUsecase) hydrateTile(conf *models.Config, tile *models.Tile, host string) {
	// Empty tile, skip
	if tile.Type == EmptyTileType {
		return
	}

	if tile.Type == GroupTileType {
		cu.hydrateTiles(conf, &tile.Tiles, host)
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

func (cu *configUsecase) hydrateDynamicTile(conf *models.Config, tile *models.Tile) (tiles []models.Tile) {
	dynamicTileConfig := cu.dynamicTileConfigs[tile.Type][tile.ConfigVariant]

	// Create new validator by reflexion
	rType := reflect.TypeOf(dynamicTileConfig.Validator)
	rInstance := reflect.New(rType.Elem()).Interface()

	// Marshal / Unmarshal the map[string]interface{} struct in new instance of Validator
	bParams, _ := json.Marshal(tile.Params)
	_ = json.Unmarshal(bParams, &rInstance)

	// Call builder and add inherited value from Dynamic tile
	results, err := dynamicTileConfig.Builder.ListDynamicTile(rInstance)
	if err != nil {
		// TODO : Replace that by errors.Is/As when go 1.13 will be released
		params, _ := json.Marshal(tile.Params)
		if strings.Contains(err.Error(), "unable to found job") {
			conf.AddErrors(fmt.Sprintf(`Error while listing %s dynamic tiles (params: %s). %v`, tile.Type, string(params), err))
		} else {
			conf.AddWarnings(fmt.Sprintf(`Warning while listing %s dynamic tiles (params: %s). %v`, tile.Type, string(params), err))
		}
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
