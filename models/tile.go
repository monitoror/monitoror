package models

type (
	Tile struct {
		Type   TileType   `json:"type"`
		Status TileStatus `json:"status"`

		Label   string `json:"label,omitempty"`
		Message string `json:"message,omitempty"`

		Value *TileValue `json:"value,omitempty"`
		Build *TileBuild `json:"build,omitempty"`
	}

	TileType   string //PING, PORT, ... (defined in usecase.go for each monitorable)
	TileStatus string
)

const (
	ActionRequiredStatus TileStatus = "ACTION_REQUIRED"
	CanceledStatus       TileStatus = "CANCELED"
	DisabledStatus       TileStatus = "DISABLED"
	FailedStatus         TileStatus = "FAILURE"
	QueuedStatus         TileStatus = "QUEUED"
	RunningStatus        TileStatus = "RUNNING"
	SuccessStatus        TileStatus = "SUCCESS"
	UnknownStatus        TileStatus = "UNKNOWN"
	WarningStatus        TileStatus = "WARNING"
)

func NewTile(t TileType) *Tile {
	return &Tile{Type: t}
}
