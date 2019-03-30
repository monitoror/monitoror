package config

import (
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitConfig_Default(t *testing.T) {
	config, err := InitConfig()

	assert.NoError(t, err)
	assert.Equal(t, 8080, config.Port)
}

func TestInitConfig_WithWrongFile(t *testing.T) {
	dest, err := copyConfig("./testdata/server-config-broken.yml")
	assert.NoError(t, err, "unable to copy config file")
	defer os.Remove(dest)

	config, err := InitConfig()

	assert.Error(t, err)
	assert.Nil(t, config)
	assert.EqualError(t, err, "unable to init config, While parsing config: yaml: line 11: could not find expected ':'")
}

func TestInitConfig_WithFile(t *testing.T) {
	dest, err := copyConfig("./testdata/server-config.yml")
	assert.NoError(t, err, "unable to copy config file")
	defer os.Remove(dest)

	config, err := InitConfig()

	assert.NoError(t, err)
	assert.Equal(t, 3000, config.Port)
}

func TestInitConfig_WithEnv(t *testing.T) {
	err := os.Setenv(EnvPrefix+"_PORT", "3000")
	assert.NoError(t, err, "unable to setEnv")

	config, err := InitConfig()

	assert.NoError(t, err)
	assert.Equal(t, 3000, config.Port)
}

func copyConfig(src string) (string, error) {
	in, err := os.Open(src)
	if err != nil {
		return "", err
	}
	defer in.Close()

	dest := "./" + ServerConfigFileName + filepath.Ext(src)
	out, err := os.Create(dest)
	if err != nil {
		return "", err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return "", err
	}
	return dest, out.Close()
}
