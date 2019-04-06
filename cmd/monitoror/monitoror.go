package main

import (
	"github.com/monitoror/monitoror/cli"
	"github.com/monitoror/monitoror/config"
	"github.com/monitoror/monitoror/service"
)

func main() {
	// Load Config from File/Env
	config, err := config.InitConfig()
	if err != nil {
		panic(err)
	}

	// Banner
	cli.PrintBanner()

	// Start Service
	service.Start(config)
}
