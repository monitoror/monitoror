package service

import (
	"testing"

	. "github.com/monitoror/monitoror/config"

	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	conf := InitConfig()
	conf.Env = "test"
	server := Init(conf)
	assert.NotNil(t, server.Echo)
}

func TestRegisterTile(t *testing.T) {
	registerTile(func(s string) { assert.False(t, true) }, "TEST", false)
	registerTile(func(s string) { assert.True(t, true) }, "TEST", true)
}
