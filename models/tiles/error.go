package tiles

const ErrorTileType TileType = "ERROR"

type (
	// commonTile struct used by every response of monitorable route
	ErrorTile struct {
		*Tile
	}
)

func NewErrorTile(label, message string) *ErrorTile {
	return &ErrorTile{
		Tile: &Tile{
			Type:    ErrorTileType,
			Label:   label,
			Message: message,
		},
	}
}
