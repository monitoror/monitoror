package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig_AddErrors(t *testing.T) {
	config := &ConfigBag{}
	config.AddErrors(ConfigError{})

	assert.Len(t, config.Errors, 1)
}
