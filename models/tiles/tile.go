package tiles

type (
	// commonTile struct used by every response of monitorable route
	Tile struct {
		Type    TileType    `json:"type"`
		SubType TileSubType `json:"subtype,omitempty"`
		Status  TileStatus  `json:"status,omitempty"`
		Label   string      `json:"label"`
		Message string      `json:"message,omitempty"`
	}

	TileType    string // BUILD, HEALTH ...
	TileSubType string // PING, JENKINS_JOB ...
	TileStatus  string // SUCCESS, FAILURE ...
)

// List of all Response Status
const (
	SuccessStatus TileStatus = "SUCCESS"
	FailStatus    TileStatus = "FAILURE"
	TimeoutStatus TileStatus = "TIMEOUT"
)
