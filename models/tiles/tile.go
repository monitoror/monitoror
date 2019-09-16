package tiles

type (
	// commonTile struct used by every response of monitorable route
	Tile struct {
		Category TileCategory `json:"category"`
		Type     TileType     `json:"type,omitempty"`
		Status   TileStatus   `json:"status,omitempty"`
		Label    string       `json:"label,omitempty"`
		Message  string       `json:"message,omitempty"`
	}

	TileCategory string // BUILD, HEALTH ...
	TileType     string // PING, JENKINS_BUILD ...
	TileStatus   string // SUCCESS, FAILURE ...
)

// List of all Status Code
const (
	SuccessStatus  TileStatus = "SUCCESS"
	FailedStatus   TileStatus = "FAILURE"
	RunningStatus  TileStatus = "RUNNING"
	QueuedStatus   TileStatus = "QUEUED"
	DisabledStatus TileStatus = "DISABLED"
	AbortedStatus  TileStatus = "ABORTED"
	WarningStatus  TileStatus = "WARNING"
	UnknownStatus  TileStatus = "UNKNOWN"
)
