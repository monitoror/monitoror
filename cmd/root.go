package main

import (
	"github.com/jsdidierlaurent/monitowall/config"
	"github.com/jsdidierlaurent/monitowall/service"
)

func main() {
	//TODO: Adding "debug-configuration" flag to cmd to print config and missing config file

	config, err := config.Load()
	if err != nil {
		panic(err)
	}

	service.Start(config)
}
