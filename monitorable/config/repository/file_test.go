package repository

import (
	"io/ioutil"
	"os"
	"testing"

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
		_, err := repository.GetConfigFromPath(tmpFile.Name())
		assert.NoError(t, err)
	}
}

func TestConfigRepository_GetConfigFromPath_Error(t *testing.T) {
	repository := NewConfigRepository()
	_, err := repository.GetConfigFromPath("monitoror-missing-file")
	assert.Error(t, err)
}
