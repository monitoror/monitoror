package tiles

const HealthTileType TileType = "HEALTH"

type (
	// Used by Ping, Port ... as response structure
	HealthTile struct {
		*Tile
	}
)

func NewHealthTile(subType TileSubType) *HealthTile {
	return &HealthTile{
		Tile: &Tile{
			Type:    HealthTileType,
			SubType: subType,
		},
	}
}
