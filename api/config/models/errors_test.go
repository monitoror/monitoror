package models

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigFileNotFoundError(t *testing.T) {
	err := &ConfigFileNotFoundError{Err: errors.New("boom"), PathOrURL: "test"}
	assert.Equal(t, "Config not found at: test, boom", err.Error())
	assert.Equal(t, "boom", err.Unwrap().Error())

	err = &ConfigFileNotFoundError{PathOrURL: "test"}
	assert.Equal(t, "Config not found at: test", err.Error())
	assert.Equal(t, nil, err.Unwrap())
}

func TestConfigVersionFormatError(t *testing.T) {
	err := &ConfigVersionFormatError{WrongVersion: "10"}
	assert.Equal(t, "json: cannot unmarshal 10 into Go struct field Config.Version of type string and X.y format", err.Error())
}

func TestConfigUnmarshalError(t *testing.T) {
	err := &ConfigUnmarshalError{Err: errors.New("boom"), RawConfig: "test"}
	assert.Equal(t, "boom", err.Error())
	assert.Equal(t, "boom", err.Unwrap().Error())
}
