package renderings

type (
	Response struct {
		Type    TileType   `json:"type"`
		Status  TileStatus `json:"status"`
		Message string     `json:"message,omitempty"`
	}

	// HealthCheckResponse response for PING / PORT / ...
	HealthCheckResponse struct {
		*Response
		Label string `json:"label"`
	}

	// BuildStatusResponse response for JENKINS_JOB / GITLAB_PIPELINE / ...
	BuildStatusResponse struct {
		*Response
	}

	TileType   string
	TileStatus string
)

func NewHealthCheckResponse() *HealthCheckResponse {
	return &HealthCheckResponse{
		Response: &Response{},
	}
}

func NewBuildStatusResponse() *BuildStatusResponse {
	return &BuildStatusResponse{
		Response: &Response{},
	}
}

//List of all available types of tiles for monitowall
const (
	TypePing TileType = "PING"
)

// List of all Response Status
const (
	SuccessStatus TileStatus = "SUCCESS"
	FailStatus    TileStatus = "FAILURE"
	TimeoutStatus TileStatus = "TIMEOUT"
)
