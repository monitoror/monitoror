package tiles

const CounterTileCategory TileCategory = "BUILD"

type (
	// Used by Youtrack, Sentry ... as response structure
	CounterTile struct {
		*Tile

		//TODO
	}
)

func NewCounterTile(t TileType) *CounterTile {
	return &CounterTile{
		Tile: &Tile{
			Category: CounterTileCategory,
			Type:     t,
		},
	}
}
