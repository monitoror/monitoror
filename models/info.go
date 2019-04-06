package models

import (
	"github.com/monitoror/monitoror/config"
)

type (
	//InfoResponse response for info route
	InfoResponse struct {
		BuildInfo BuildInfo     `json:"build-info"`
		Config    config.Config `json:"configuration"`
	}

	BuildInfo struct {
		Version   string `json:"version"`
		GitCommit string `json:"git-commit"`
		BuildTime string `json:"build-time"`
	}
)

func NewInfoResponse(version string, gitCommit string, buildTime string, config *config.Config) *InfoResponse {
	return &InfoResponse{
		BuildInfo: BuildInfo{
			Version:   version,
			GitCommit: gitCommit,
			BuildTime: buildTime,
		},
		Config: *config,
	}
}
