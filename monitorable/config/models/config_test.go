package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig_AddErrors(t *testing.T) {
	conf := Config{
		Errors: []string{},
	}

	conf.AddErrors("error1", "error2")
	assert.Len(t, conf.Errors, 2)
}

func TestConfig_AddWarnings(t *testing.T) {
	conf := Config{
		Errors: []string{},
	}

	conf.AddWarnings("warn1", "warn2")
	assert.Len(t, conf.Warnings, 2)
}
