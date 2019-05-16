package tiles

type (
	// commonTile struct used by every response of monitorable route
	Tile struct {
		Category TileCategory `json:"category"`
		Type     TileType     `json:"type,omitempty"`
		Status   TileStatus   `json:"status,omitempty"`
		Label    string       `json:"label"`
		Message  string       `json:"message,omitempty"`
	}

	TileCategory string // BUILD, HEALTH ...
	TileType     string // PING, JENKINS_JOB ...
	TileStatus   string // SUCCESS, FAILURE ...
)

// List of all Response Code
const (
	SuccessStatus  TileStatus = "SUCCESS"
	FailedStatus   TileStatus = "FAILURE"
	RunningStatus  TileStatus = "RUNNING"
	QueuedStatus   TileStatus = "QUEUED"
	WarningStatus  TileStatus = "WARNING"
	CanceledStatus TileStatus = "CANCELED"
	UnknownStatus  TileStatus = "UNKNOWN"
)
