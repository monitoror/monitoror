package service

import (
	"testing"

	"github.com/GeertJohan/go.rice/embedded"
	"github.com/monitoror/monitoror/config"
	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	conf := config.InitConfig()
	conf.Env = "test"
	server := Init(conf)
	assert.NotNil(t, server.Echo)
}

func TestInitFront_Panic(t *testing.T) {
	conf := config.InitConfig()
	conf.Env = "production"

	server := &Server{
		config: conf,
	}
	server.initEcho()

	assert.Panics(t, func() {
		server.initUI()
	})
}

func TestInitFront_Success(t *testing.T) {
	conf := config.InitConfig()
	conf.Env = "production"

	server := &Server{
		config: conf,
	}
	server.initEcho()

	embedded.RegisterEmbeddedBox("../ui/dist", &embedded.EmbeddedBox{
		Name: "../ui/dist",
	})
	defer delete(embedded.EmbeddedBoxes, "../ui/dist")

	server.initUI()
}

func TestRegisterTile(t *testing.T) {
	registerTile(func(s string) { assert.False(t, true) }, "TEST", false)
	registerTile(func(s string) { assert.True(t, true) }, "TEST", true)
}
