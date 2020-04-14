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
	err := os.Setenv(EnvPrefix+"_PORT", "3000")
	assert.NoError(t, err, "unable to setEnv")
	err = os.Setenv(EnvPrefix+"_MONITORABLE_JENKINS_URL", "test")
	assert.NoError(t, err, "unable to setEnv")
	err = os.Setenv(EnvPrefix+"_MONITORABLE_AZUREDEVOPS_URL", "test")
	assert.NoError(t, err, "unable to setEnv")

	config := InitConfig()

	assert.Equal(t, "production", config.Env)
	assert.Equal(t, 3000, config.Port)
}
