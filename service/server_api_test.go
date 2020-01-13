package service

import (
	"testing"

	"github.com/monitoror/monitoror/config"
)

func initServer() *Server {
	conf := config.InitConfig()
	server := &Server{config: conf}
	server.initEcho()
	server.initMiddleware()
	server.api = server.Echo.Group("/test/api/v1")

	server.registerConfig()

	return server
}

func TestRegisterPing(t *testing.T) {
	server := initServer()
	server.registerPing(config.DefaultVariant)
}

func TestRegisterPort(t *testing.T) {
	server := initServer()
	server.registerPort(config.DefaultVariant)
}

func TestRegisterHTTP(t *testing.T) {
	server := initServer()
	server.registerHTTP(config.DefaultVariant)
}

func TestRegisterPingdom(t *testing.T) {
	server := initServer()
	server.registerPingdom(config.DefaultVariant)
}

func TestRegisterJenkins(t *testing.T) {
	server := initServer()
	server.registerJenkins(config.DefaultVariant)
}

func TestRegisterTravisCI(t *testing.T) {
	server := initServer()
	server.registerTravisCI(config.DefaultVariant)
}

func TestRegisterAzureDevOps(t *testing.T) {
	server := initServer()
	server.registerAzureDevOps(config.DefaultVariant)
}
