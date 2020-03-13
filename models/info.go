package models

type (
	// InfoResponse response for info route
	InfoResponse struct {
		Version   string `json:"version"`
		GitCommit string `json:"git-commit"`
		BuildTime string `json:"build-time"`
		BuildTags string `json:"build-tags"`
	}
)

func NewInfoResponse(version, gitCommit, buildTime, buildTags string) *InfoResponse {
	return &InfoResponse{
		Version:   version,
		GitCommit: gitCommit,
		BuildTime: buildTime,
		BuildTags: buildTags,
	}
}
