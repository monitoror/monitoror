package main

import (
	"github.com/jsdidierlaurent/monitowall/config"
	"github.com/jsdidierlaurent/monitowall/service"
)

func main() {
	// Load Config from File/Env
	config, err := config.InitConfig()
	if err != nil {
		panic(err)
	}

	// Start Service
	service.Start(config)
}
