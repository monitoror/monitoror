package configs

type (
	BuildInfo struct {
		GitCommit string `json:"git-commit"`
		Version   string `json:"version"`
		BuildTime string `json:"build-time"`
		OS        string `json:"os"`
		Arch      string `json:"arch"`
	}
)

func InitBuildInfo(gitCommit, version, buildTime, os, arch string) *BuildInfo {
	return &BuildInfo{
		GitCommit: gitCommit,
		Version:   version,
		BuildTime: buildTime,
		OS:        os,
		Arch:      arch,
	}
}
