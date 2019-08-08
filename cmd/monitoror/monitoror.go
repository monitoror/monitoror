package main

import (
	"github.com/joho/godotenv"
	"github.com/labstack/gommon/log"

	"github.com/monitoror/monitoror/cli"
	"github.com/monitoror/monitoror/config"
	"github.com/monitoror/monitoror/service"
)

func main() {
	//  Default Logger
	log.SetPrefix("")
	log.SetHeader("[${level}]")
	log.SetLevel(log.INFO)

	// Load .env file
	_ = godotenv.Load(".env.local")
	_ = godotenv.Load(".env")

	// Load GetConfig from File/Env
	config := config.InitConfig()

	// Banner
	cli.PrintBanner()

	// Start Service
	server := service.Init(config)
	server.Start()
}
