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

	config := InitConfig()

	assert.Equal(t, "production", config.Env)
	assert.Equal(t, 3000, config.Port)
	assert.Equal(t, 2, config.Monitorable.Ping.Count)
	assert.Equal(t, 2000, config.Monitorable.Port.Timeout)
	assert.Equal(t, "https://api.travis-ci.org/", config.Monitorable.TravisCI[DefaultVariant].Url)
	assert.Equal(t, true, config.Monitorable.Jenkins[DefaultVariant].SSLVerify)
	assert.Equal(t, 2000, config.Monitorable.Jenkins[DefaultVariant].Timeout)
	assert.Equal(t, "test", config.Monitorable.Jenkins[DefaultVariant].Url)
}

func TestPingdom_IsValid(t *testing.T) {
	config := InitConfig()

	config.Monitorable.Pingdom[DefaultVariant].Url = ""
	config.Monitorable.Pingdom[DefaultVariant].Token = ""
	assert.False(t, config.Monitorable.Pingdom[DefaultVariant].IsValid())

	config.Monitorable.Pingdom[DefaultVariant].Url = ""
	config.Monitorable.Pingdom[DefaultVariant].Token = "abcde"
	assert.True(t, config.Monitorable.Pingdom[DefaultVariant].IsValid())

	config.Monitorable.Pingdom[DefaultVariant].Url = "url%url"
	config.Monitorable.Pingdom[DefaultVariant].Token = "abcde"
	assert.False(t, config.Monitorable.Pingdom[DefaultVariant].IsValid())

	config.Monitorable.Pingdom[DefaultVariant].Url = "http://pingdom.com/api"
	config.Monitorable.Pingdom[DefaultVariant].Token = "abcde"
	assert.True(t, config.Monitorable.Pingdom[DefaultVariant].IsValid())
}

func TestTravisCI_IsValid(t *testing.T) {
	config := InitConfig()
	assert.True(t, config.Monitorable.TravisCI[DefaultVariant].IsValid())

	config.Monitorable.TravisCI[DefaultVariant].Url = "url%url"
	assert.False(t, config.Monitorable.TravisCI[DefaultVariant].IsValid())

	config.Monitorable.TravisCI[DefaultVariant].Url = ""
	assert.False(t, config.Monitorable.TravisCI[DefaultVariant].IsValid())
}

func TestJenkins_IsValid(t *testing.T) {
	config := InitConfig()

	config.Monitorable.Jenkins[DefaultVariant].Url = ""
	assert.False(t, config.Monitorable.Jenkins[DefaultVariant].IsValid())

	config.Monitorable.Jenkins[DefaultVariant].Url = "url%url"
	assert.False(t, config.Monitorable.Jenkins[DefaultVariant].IsValid())

	config.Monitorable.Jenkins[DefaultVariant].Url = "http://jenkins.test.com"
	assert.True(t, config.Monitorable.Jenkins[DefaultVariant].IsValid())
}
