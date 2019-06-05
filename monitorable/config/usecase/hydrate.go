package usecase

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/monitoror/monitoror/models/tiles"
	"github.com/monitoror/monitoror/monitorable/config/models"
)

func (cu *configUsecase) Hydrate(config *models.Config, host string) error {
	for _, tile := range config.Tiles {
		cu.hydrateTile(tile, host)
	}

	return nil
}

func (cu *configUsecase) hydrateTile(tile map[string]interface{}, host string) {
	tileType := tiles.TileType(strings.ToUpper(tile[TypeKey].(string)))

	// Empty tile, skip
	if tileType == EmptyTileType {
		return
	}

	if tileType == GroupTileType {
		groupTiles, _ := tile[TilesKey].([]interface{})
		for _, gt := range groupTiles {
			groupTile, _ := gt.(map[string]interface{})
			cu.hydrateTile(groupTile, host)
		}

		return
	}

	// Change Params by a valid Url
	path := cu.tileConfigs[tileType].Path
	params := url.Values{}
	for key, value := range tile[ParamsKey].(map[string]interface{}) {
		params.Add(key, fmt.Sprintf("%v", value))
	}

	tile[UrlKey] = fmt.Sprintf("%s%s?%s", host, path, params.Encode())

	// Remove Params
	delete(tile, ParamsKey)
}
