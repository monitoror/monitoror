package config

type (
	BuildInfo struct {
		GitCommit string `json:"git-commit"`
		Version   string `json:"version"`
		BuildTime string `json:"build-time"`
		OS        string `json:"os"`
		Arch      string `json:"arch"`
	}
)
