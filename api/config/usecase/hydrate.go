package usecase

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"reflect"

	"github.com/monitoror/monitoror/api/config/models"
	coreModels "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/pkg/humanize"
)

func (cu *configUsecase) Hydrate(configBag *models.ConfigBag) {
	cu.hydrateTiles(configBag, &configBag.Config.Tiles)
}

func (cu *configUsecase) hydrateTiles(configBag *models.ConfigBag, tiles *[]models.TileConfig) {
	for i := 0; i < len(*tiles); i++ {
		tile := &((*tiles)[i])
		if tile.Type != EmptyTileType && tile.Type != GroupTileType {
			// Set ConfigVariant to DefaultVariant if empty
			if tile.ConfigVariant == "" {
				tile.ConfigVariant = coreModels.DefaultVariant
			}
		}

		if _, exists := cu.configData.dynamicTileConfigs[tile.Type]; !exists {
			cu.hydrateTile(configBag, tile)

			if tile.Type == GroupTileType && len(tile.Tiles) == 0 {
				*tiles = append((*tiles)[:i], (*tiles)[i+1:]...)
				i--
			}
		} else {
			dynamicTiles := cu.hydrateDynamicTile(configBag, tile)

			// Remove DynamicTile config and add real dynamic tiles in array
			temp := append((*tiles)[:i], dynamicTiles...)
			*tiles = append(temp, (*tiles)[i+1:]...)

			i--
		}
	}
}

func (cu *configUsecase) hydrateTile(configBag *models.ConfigBag, tile *models.TileConfig) {
	// Empty tile, skip
	if tile.Type == EmptyTileType {
		return
	}

	if tile.Type == GroupTileType {
		cu.hydrateTiles(configBag, &tile.Tiles)
		return
	}

	tileConfig := cu.configData.tileConfigs[tile.Type][tile.ConfigVariant]

	// Change Params by a valid URL
	urlParams := url.Values{}
	for key, value := range tile.Params {
		// Array of value
		if reflect.TypeOf(value).Kind() == reflect.Slice {
			for _, v := range value.([]interface{}) {
				urlParams.Add(key, humanize.Interface(v))
			}
		} else {
			urlParams.Add(key, humanize.Interface(value))
		}
	}
	tile.URL = fmt.Sprintf("%s?%s", tileConfig.Path, urlParams.Encode())

	// Add initial max delay from config
	tile.InitialMaxDelay = &tileConfig.InitialMaxDelay

	// Remove Params / Variant
	tile.Params = nil
	tile.ConfigVariant = ""
}

func (cu *configUsecase) hydrateDynamicTile(configBag *models.ConfigBag, tile *models.TileConfig) []models.TileConfig {
	dynamicTileConfig := cu.configData.dynamicTileConfigs[tile.Type][tile.ConfigVariant]

	// Create new validator by reflexion
	rType := reflect.TypeOf(dynamicTileConfig.Validator)
	rInstance := reflect.New(rType.Elem()).Interface()

	// Marshal / Unmarshal the map[string]interface{} struct in new instance of Validator
	bParams, _ := json.Marshal(tile.Params)
	_ = json.Unmarshal(bParams, &rInstance)

	// Call builder and add inherited value from Dynamic tile
	cacheKey := fmt.Sprintf("%s:%s_%s_%s", DynamicTileStoreKeyPrefix, tile.Type, tile.ConfigVariant, string(bParams))
	results, err := dynamicTileConfig.Builder(rInstance)
	if err != nil {
		if os.IsTimeout(err) {
			// Get previous value in cache
			if err := cu.dynamicTileStore.Get(cacheKey, &results); err != nil {
				configBag.AddErrors(models.ConfigError{
					ID:      models.ConfigErrorUnableToHydrate,
					Message: fmt.Sprintf(`Error while listing %s dynamic tiles (params: %s). Timeout or host unreachable`, tile.Type, string(bParams)),
					Data: models.ConfigErrorData{
						ConfigExtract: stringify(tile),
					},
				})
			}
		} else {
			configBag.AddErrors(models.ConfigError{
				ID:      models.ConfigErrorUnableToHydrate,
				Message: fmt.Sprintf(`Error while listing %s dynamic tiles (params: %s). %v`, tile.Type, string(bParams), err),
				Data: models.ConfigErrorData{
					ConfigExtract: stringify(tile),
				},
			})
		}
	} else {
		// Add result in cache
		_ = cu.dynamicTileStore.Set(cacheKey, results, cu.cacheExpiration)
	}

	var tiles []models.TileConfig
	for _, result := range results {
		newTile := models.TileConfig{
			Type:          result.TileType,
			Label:         result.Label,
			Params:        result.Params,
			ConfigVariant: tile.ConfigVariant,
			ColumnSpan:    tile.ColumnSpan,
			RowSpan:       tile.RowSpan,
		}

		// Override tile if dynamicTile has a label
		if tile.Label != "" {
			newTile.Label = tile.Label
		}

		tiles = append(tiles, newTile)
	}

	return tiles
}
