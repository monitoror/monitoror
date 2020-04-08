package models

import (
	"fmt"
	"testing"

	"github.com/AlekSi/pointer"
	"github.com/monitoror/monitoror/internal/pkg/monitorable/validator"
	"github.com/stretchr/testify/assert"
)

func TestBuildParams_Validate(t *testing.T) {
	param := &BuildParams{}
	assert.Error(t, validator.Validate(param))

	param.Project = "test"
	assert.Error(t, validator.Validate(param))

	param.Definition = pointer.ToInt(1)
	assert.NoError(t, validator.Validate(param))

	param.Branch = pointer.ToString("test")
	assert.NoError(t, validator.Validate(param))
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
