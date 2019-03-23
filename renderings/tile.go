package renderings

type (
	// HealthCheckResponse response for PING / PORT / ...
	HealthCheckResponse struct {
		Type    TileType   `json:"type"`
		Status  TileStatus `json:"status"`
		Label   string     `json:"label"`
		Message string     `json:"message,omitempty"`
	}

	// BuildStatusResponse response for JENKINS_JOB / GITLAB_PIPELINE / ...
	BuildStatusResponse struct {
	}

	TileType   string
	TileStatus string
)

//List of all available types of tiles for monitowall
const (
	TypePing TileType = "PING"
)

// List of all Response Status
const (
	SuccessStatus TileStatus = "SUCCESS"
	FailStatus    TileStatus = "FAILURE"
)
