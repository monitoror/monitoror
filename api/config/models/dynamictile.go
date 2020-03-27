package models

import "github.com/monitoror/monitoror/models"

type (
	DynamicTileResult struct {
		TileType models.TileType
		Label    string
		Params   map[string]interface{}
	}
)
