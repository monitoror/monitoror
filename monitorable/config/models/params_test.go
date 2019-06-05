package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigParams_IsValid_Success(t *testing.T) {
	configParams := &ConfigParams{Path: "Path"}
	assert.True(t, configParams.IsValid())

	configParams = &ConfigParams{Url: "Url"}
	assert.True(t, configParams.IsValid())
}

func TestConfigParams_IsValid_Error(t *testing.T) {
	configParams := &ConfigParams{}
	assert.False(t, configParams.IsValid())

	configParams = &ConfigParams{Path: "Path", Url: "Url"}
	assert.False(t, configParams.IsValid())
}
