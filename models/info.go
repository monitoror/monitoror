package models

type (
	// InfoResponse response for info route
	InfoResponse struct {
		Version   string `json:"version"`
		GitCommit string `json:"git-commit"`
		BuildTime string `json:"build-time"`
	}
)

func NewInfoResponse(version string, gitCommit string, buildTime string) *InfoResponse {
	return &InfoResponse{
		Version:   version,
		GitCommit: gitCommit,
		BuildTime: buildTime,
	}
}
