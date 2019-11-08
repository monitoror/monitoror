package service

import (
	"testing"

	. "github.com/monitoror/monitoror/config"

	"github.com/stretchr/testify/assert"
)

// This tests are really basic and just check if Server member are not nil
func TestInit_WithAllTile(t *testing.T) {
	conf := InitConfig()
	conf.Env = "Test"

	conf.Monitorable.Jenkins["jenkins"] = &Jenkins{URL: "http://jenkins.test.com"}
	conf.Monitorable.Pingdom["pingdom"] = &Pingdom{Token: "abcdef"}
	conf.Monitorable.AzureDevOps["azure-devops"] = &AzureDevOps{URL: "https://dev.azure.com/test", Token: "abcdef"}

	server := &Server{config: conf}
	server.initEcho()
	server.initMiddleware()
	server.initFront()
	server.initApis()

	assert.NotNil(t, server.Echo)
}

// This tests are really basic and just check if Server member are not nil
func TestInit_WithoutAllTile(t *testing.T) {
	conf := InitConfig()
	conf.Env = "Test"
	conf.Monitorable.TravisCI[DefaultVariant].URL = ""

	server := &Server{config: conf}
	server.initEcho()
	server.initMiddleware()
	server.initFront()
	server.initApis()

	assert.NotNil(t, server.Echo)
}
