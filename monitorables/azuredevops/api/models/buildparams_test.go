package models

import (
	"fmt"
	"testing"

	"github.com/AlekSi/pointer"
	"github.com/stretchr/testify/assert"
)

func TestBuildParams_IsValid(t *testing.T) {
	param := &BuildParams{}
	assert.False(t, param.IsValid())

	param.Project = "test"
	assert.False(t, param.IsValid())

	param.Definition = pointer.ToInt(1)
	assert.True(t, param.IsValid())

	param.Branch = pointer.ToString("test")
	assert.True(t, param.IsValid())
}

func TestBuildParams_String(t *testing.T) {
	param := &BuildParams{
		Project:    "test",
		Definition: pointer.ToInt(1),
	}
	assert.Equal(t, "BUILD-test-1", fmt.Sprint(param))

	param.Branch = pointer.ToString("test")
	assert.Equal(t, "BUILD-test-1-test", fmt.Sprint(param))
}
