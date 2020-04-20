package main

import (
	"os"
	"path/filepath"

	"github.com/monitoror/monitoror/cli"
	"github.com/monitoror/monitoror/config"
	"github.com/monitoror/monitoror/service"

	"github.com/joho/godotenv"
	"github.com/labstack/gommon/log"
)

func main() {
	//  Default Logger
	log.SetPrefix("")
	log.SetHeader("[${level}]")
	log.SetLevel(log.INFO)

	// GetConfig .env file
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))

	_ = godotenv.Load(".env")
	_ = godotenv.Load(filepath.Join(dir, ".env"))

	// GetConfig GetConfig from File/Env
	conf := config.InitConfig()

	// CLI
	cli := cli.New()
	cli.PrintBanner()

	// Start Service
	server := service.Init(conf, cli)
	server.Start()
}
