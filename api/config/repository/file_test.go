package repository

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/monitoror/monitoror/api/config/models"

	"github.com/stretchr/testify/assert"
)

// /!\ this is an integration test /!\
// Note : It may be necessary to separate them from unit tests

// TestConfigRepository_GetConfigFromPath test if os.Open get works
func TestConfigRepository_GetConfigFromPath(t *testing.T) {
	tmpFile, err := ioutil.TempFile(os.TempDir(), "test-config-GetConfigFromPath-")
	if assert.NoError(t, err) {
		defer os.Remove(tmpFile.Name())
		_, _ = tmpFile.WriteString("{}")

		repository := NewConfigRepository()
		_, err := repository.GetConfigFromPath("", tmpFile.Name())
		assert.NoError(t, err)
	}
}

func TestConfigRepository_UnableToParse(t *testing.T) {
	tmpFile, err := ioutil.TempFile(os.TempDir(), "monitoror-wrong-file")
	if assert.NoError(t, err) {
		defer os.Remove(tmpFile.Name())
		_, _ = tmpFile.WriteString("xxxxxx")

		repository := NewConfigRepository()
		_, err := repository.GetConfigFromPath("", tmpFile.Name())
		assert.Error(t, err)
		assert.Equal(t, "xxxxxx", err.(*models.ConfigUnmarshalError).RawConfig)
	}
}

func TestConfigRepository_GetConfigFromPath_MissingFile(t *testing.T) {
	repository := NewConfigRepository()
	_, err := repository.GetConfigFromPath("/test", "monitoror-missing-file")
	assert.Error(t, err)
	assert.Equal(t, "Config not found at: /test/monitoror-missing-file, open /test/monitoror-missing-file: no such file or directory", err.Error())
}
