//+build debug

package main

import (
	"github.com/jsdidierlaurent/monitowall/configs"
	"github.com/jsdidierlaurent/monitowall/service"
)

var (
	GitCommit, Version, BuildTime, OS, Arch string
)

func main() {
	// Setup BuildInfo struct
	buildInfo := configs.InitBuildInfo(GitCommit, Version, BuildTime, OS, Arch)

	// Load Config from File/Env
	config, err := configs.InitConfig()
	if err != nil {
		panic(err)
	}

	// Start Service
	service.StartDebug(config, buildInfo)
}
