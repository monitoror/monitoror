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
