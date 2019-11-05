package models

import (
	"fmt"
	"testing"

	"github.com/AlekSi/pointer"
	"github.com/stretchr/testify/assert"
)

func TestReleaseParams_IsValid(t *testing.T) {
	param := &ReleaseParams{}
	assert.False(t, param.IsValid())

	param.Project = "test"
	assert.False(t, param.IsValid())

	param.Definition = pointer.ToInt(1)
	assert.True(t, param.IsValid())
}

func TestReleaseParams_String(t *testing.T) {
	param := &ReleaseParams{
		Project:    "test",
		Definition: pointer.ToInt(1),
	}
	assert.Equal(t, "RELEASE-test-1", fmt.Sprint(param))
}
