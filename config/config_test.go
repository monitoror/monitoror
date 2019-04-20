package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitConfig_Default(t *testing.T) {
	config, err := InitConfig()

	assert.NoError(t, err)
	assert.Equal(t, 8080, config.Port)
}

func TestInitConfig_WithEnv(t *testing.T) {
	err := os.Setenv(EnvPrefix+"_PORT", "3000")
	assert.NoError(t, err, "unable to setEnv")

	config, err := InitConfig()

	assert.NoError(t, err)
	assert.Equal(t, 3000, config.Port)
	assert.Equal(t, 2, config.Monitorable.Ping.Count)
	assert.Equal(t, 2000, config.Monitorable.Port.Timeout)
}
