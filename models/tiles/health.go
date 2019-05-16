package tiles

const HealthTileCategory TileCategory = "HEALTH"

type (
	// Used by Ping, Port ... as response structure
	HealthTile struct {
		*Tile
	}
)

func NewHealthTile(t TileType) *HealthTile {
	return &HealthTile{
		Tile: &Tile{
			Category: HealthTileCategory,
			Type:     t,
		},
	}
}
