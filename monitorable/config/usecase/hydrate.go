package usecase

import (
	"fmt"
	"net/url"

	. "github.com/monitoror/monitoror/config"

	"github.com/monitoror/monitoror/monitorable/config/models"
)

func (cu *configUsecase) Hydrate(config *models.Config, host string) error {
	for i := range config.Tiles {
		cu.hydrateTile(&config.Tiles[i], host)
	}

	return nil
}

func (cu *configUsecase) hydrateTile(tile *models.Tile, host string) {
	// Empty tile, skip
	if tile.Type == EmptyTileType {
		return
	}

	if tile.Type == GroupTileType {
		for i := range tile.Tiles {
			cu.hydrateTile(&tile.Tiles[i], host)
		}
		return
	}

	// Define config Variant
	if tile.ConfigVariant == "" {
		tile.ConfigVariant = DefaultVariant
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
