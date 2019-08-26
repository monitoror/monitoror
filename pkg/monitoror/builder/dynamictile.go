package builder

import (
	"github.com/monitoror/monitoror/models/tiles"
)

type (
	DynamicTileBuilder interface {
		ListDynamicTile(params interface{}) ([]Result, error)
	}

	Result struct {
		TileType tiles.TileType
		Label    string
		Params   map[string]interface{}
	}
)
