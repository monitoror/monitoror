package renderings

type (
	Type   string
	Status string

	// HealthCheckResponse response for PING / PORT / ...
	HealthCheckResponse struct {
		Type    Type   `json:"type"`
		Status  Status `json:"status"`
		Label   string `json:"label"`
		Message string `json:"message,omitempty"`
	}

	// BuildStatusResponse response for JENKINS_JOB / GITLAB_PIPELINE / ...
	BuildStatusResponse struct {
	}
)

//List of all available types of tiles for monitowall
const (
	TypePing Type = "PING"
)

// List of all Response Status
const (
	SuccessStatus Status = "SUCCESS"
	FailStatus    Status = "FAILURE"
)
