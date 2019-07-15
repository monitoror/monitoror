package service

import (
	"testing"

	"github.com/monitoror/monitoror/config"

	"github.com/stretchr/testify/assert"
)

// This tests are realy basic and juste check if Server member are not nil
func TestInit_WithAllTile(t *testing.T) {
	conf := config.InitConfig()
	conf.Env = "Test"
	conf.Monitorable.Jenkins.Url = "http://jenkins.test.com"

	server := &Server{config: conf}
	server.initEcho()
	server.initMiddleware()
	server.initFront()
	server.initApis()

	assert.NotNil(t, server.Echo)
}

// This tests are realy basic and juste check if Server member are not nil
func TestInit_WithoutAllTile(t *testing.T) {
	conf := config.InitConfig()
	conf.Env = "Test"
	conf.Monitorable.TravisCI.Url = ""

	server := &Server{config: conf}
	server.initEcho()
	server.initMiddleware()
	server.initFront()
	server.initApis()

	assert.NotNil(t, server.Echo)
}
