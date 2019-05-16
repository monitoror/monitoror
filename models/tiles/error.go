package tiles

const ErrorTileCategory TileCategory = "ERROR"

type (
	// commonTile struct used by every response of monitorable route
	ErrorTile struct {
		*Tile
	}
)

func NewErrorTile(label, message string) *ErrorTile {
	return &ErrorTile{
		Tile: &Tile{
			Category: ErrorTileCategory,
			Label:    label,
			Message:  message,
		},
	}
}
