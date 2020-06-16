package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitConfig_Default(t *testing.T) {
	config := InitConfig()
	assert.Equal(t, 8080, config.Port)
}

func TestInitConfig_WithEnv(t *testing.T) {
	assert.NoError(t, os.Setenv(EnvPrefix+"_PORT", "3000"))
	assert.NoError(t, os.Setenv(EnvPrefix+"_CONFIG", "default"))
	assert.NoError(t, os.Setenv(EnvPrefix+"_CONFIG_SCREEN1", "1"))

	config := InitConfig()

	assert.Equal(t, false, config.DisableUI)
	assert.Equal(t, 3000, config.Port)
	assert.Equal(t, "default", config.NamedConfigs["default"])
	assert.Equal(t, "1", config.NamedConfigs["screen1"])
}
