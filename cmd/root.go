package main

import (
	"github.com/jsdidierlaurent/monitowall/config"
	"github.com/jsdidierlaurent/monitowall/service"
)

var (
	GitCommit, Version, BuildTime, OS, Arch string
)

func main() {
	// Setup BuildInfo struct
	buildInfo := &config.BuildInfo{
		GitCommit: GitCommit,
		Version:   Version,
		BuildTime: BuildTime,
		OS:        OS,
		Arch:      Arch,
	}

	// Load Config from File/Env
	config, err := config.Load()
	if err != nil {
		panic(err)
	}

	// Start Service
	service.Start(config, buildInfo)
}
