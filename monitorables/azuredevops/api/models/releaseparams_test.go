package models

import (
	"fmt"
	"testing"

	"github.com/AlekSi/pointer"
	"github.com/monitoror/monitoror/internal/pkg/monitorable/validator"
	"github.com/stretchr/testify/assert"
)

func TestReleaseParams_Validate(t *testing.T) {
	param := &ReleaseParams{}
	assert.Error(t, validator.Validate(param))

	param.Project = "test"
	assert.Error(t, validator.Validate(param))

	param.Definition = pointer.ToInt(1)
	assert.NoError(t, validator.Validate(param))
}

func TestReleaseParams_String(t *testing.T) {
	param := &ReleaseParams{
		Project:    "test",
		Definition: pointer.ToInt(1),
	}
	assert.Equal(t, "RELEASE-test-1", fmt.Sprint(param))
}
