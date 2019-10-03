package builder

import "github.com/monitoror/monitoror/models"

type (
	DynamicTileBuilder interface {
		ListDynamicTile(params interface{}) ([]Result, error)
	}

	Result struct {
		TileType models.TileType
		Label    string
		Params   map[string]interface{}
	}
)
