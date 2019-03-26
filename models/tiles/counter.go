package tiles

const CounterTileType TileType = "BUILD"

type (
	// Used by Youtrack, Sentry ... as response structure
	CounterTile struct {
		*Tile

		//TODO
	}
)

func NewCounterTile(subType TileSubType) *CounterTile {
	return &CounterTile{
		Tile: &Tile{
			Type:    CounterTileType,
			SubType: subType,
		},
	}
}
