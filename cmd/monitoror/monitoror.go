package main

import (
	"github.com/jsdidierlaurent/monitoror/cli"
	"github.com/jsdidierlaurent/monitoror/config"
	"github.com/jsdidierlaurent/monitoror/service"
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
